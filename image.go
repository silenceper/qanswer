package qanswer

import (
	"fmt"
	"image"
	"image/draw"
	"sync"

	"github.com/ngaut/log"
	"github.com/silenceper/qanswer/config"
	"github.com/silenceper/qanswer/proto"
	"github.com/silenceper/qanswer/util"
)

func saveImage(png image.Image, cfg *config.Config) error {
	go func() {
		screenshotPath := fmt.Sprintf("%sscreenshot.png", proto.ImagePath)
		err := util.SavePNG(screenshotPath, png)
		if err != nil {
			log.Errorf("保存截图失败，%v", err)
		}
		log.Debugf("保存完整截图成功，%s", screenshotPath)
	}()

	//裁剪图片
	questionImg, answerImg, err := cutImage(png, cfg)
	if err != nil {
		return fmt.Errorf("截图失败，%v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		pic := thresholdingImage(questionImg)
		err = util.SavePNG(proto.QuestionImage, pic)
		if err != nil {
			log.Errorf("保存question截图失败，%v", err)
		}
		log.Debugf("保存question截图成功")
	}()

	go func() {
		defer wg.Done()
		pic := thresholdingImage(answerImg)
		err = util.SavePNG(proto.AnswerImage, pic)
		if err != nil {
			log.Errorf("保存answer截图失败，%v", err)
		}
		log.Debugf("保存answer截图成功")
	}()

	wg.Wait()
	return nil
}

//裁剪图片
func cutImage(src image.Image, cfg *config.Config) (questionImg image.Image, answerImg image.Image, err error) {
	questionImg, err = util.CutImage(src, cfg.QuestionX, cfg.QuestionY, cfg.QuestionW, cfg.QuestionH)
	// questionImg, err = util.CutImage(src, 30, 250, 660, 135)
	if err != nil {
		return
	}

	answerImg, err = util.CutImage(src, cfg.AnswerX, cfg.AnswerY, cfg.AnswerW, cfg.AnswerH)
	// answerImg, err = util.CutImage(src, 30, 420, 690, 430)
	if err != nil {
		return
	}
	return
}

//二值化图片
func thresholdingImage(img image.Image) image.Image {
	size := img.Bounds()
	pic := image.NewGray(size)
	draw.Draw(pic, size, img, size.Min, draw.Src)

	width := size.Dx()
	height := size.Dy()
	zft := make([]int, 256) //用于保存每个像素的数量，注意这里用了int类型，在某些图像上可能会溢出。
	var idx int
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			idx = i*height + j
			zft[pic.Pix[idx]]++ //image对像有一个Pix属性，它是一个slice，里面保存的是所有像素的数据。
		}
	}

	fz := getOSTUThreshold(zft)
	for i := 0; i < len(pic.Pix); i++ {
		if int(pic.Pix[i]) > fz {
			pic.Pix[i] = 255
		} else {
			pic.Pix[i] = 0
		}
	}
	return pic
}

//getOSTUThreshold OSTU大律法 计算阀值
func getOSTUThreshold(HistGram []int) int {
	var Y, Amount int
	var PixelBack, PixelFore, PixelIntegralBack, PixelIntegralFore, PixelIntegral int
	var OmegaBack, OmegaFore, MicroBack, MicroFore, SigmaB, Sigma float64 // 类间方差;
	var MinValue, MaxValue int
	var Threshold int
	for MinValue = 0; MinValue < 256 && HistGram[MinValue] == 0; MinValue++ {
	}
	for MaxValue = 255; MaxValue > MinValue && HistGram[MinValue] == 0; MaxValue-- {
	}
	if MaxValue == MinValue {
		return MaxValue // 图像中只有一个颜色
	}
	if MinValue+1 == MaxValue {
		return MinValue // 图像中只有二个颜色
	}
	for Y = MinValue; Y <= MaxValue; Y++ {
		Amount += HistGram[Y] //  像素总数
	}
	PixelIntegral = 0
	for Y = MinValue; Y <= MaxValue; Y++ {
		PixelIntegral += HistGram[Y] * Y
	}
	SigmaB = -1
	for Y = MinValue; Y < MaxValue; Y++ {
		PixelBack = PixelBack + HistGram[Y]
		PixelFore = Amount - PixelBack
		OmegaBack = float64(PixelBack) / float64(Amount)
		OmegaFore = float64(PixelFore) / float64(Amount)
		PixelIntegralBack += HistGram[Y] * Y
		PixelIntegralFore = PixelIntegral - PixelIntegralBack
		MicroBack = float64(PixelIntegralBack) / float64(PixelBack)
		MicroFore = float64(PixelIntegralFore) / float64(PixelFore)
		Sigma = OmegaBack * OmegaFore * (MicroBack - MicroFore) * (MicroBack - MicroFore)
		if Sigma > SigmaB {
			SigmaB = Sigma
			Threshold = Y
		}
	}
	return Threshold
}
