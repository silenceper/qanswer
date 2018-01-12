package qanswer

import (
	"fmt"
	"image"

	"github.com/ngaut/log"
	"github.com/silenceper/qanswer/config"
	"github.com/silenceper/qanswer/proto"
	"github.com/silenceper/qanswer/util"
)

func saveImage(png image.Image, cfg *config.Config) error {
	screenshotPath := fmt.Sprintf("%s/screenshot.png", proto.ImagePath)
	err := util.SavePNG(screenshotPath, png)
	if err != nil {
		return fmt.Errorf("保存截图失败，%v", err)
	}
	log.Debugf("保存完整截图成功，%s", screenshotPath)

	//裁剪图片
	questionImg, answerImg, err := cutImage(png, cfg)
	if err != nil {
		return fmt.Errorf("截图失败，%v", err)
	}
	err = util.SavePNG(proto.QuestionImage, questionImg)
	if err != nil {
		return fmt.Errorf("保存question截图失败，%v", err)
	}
	log.Debugf("保存question截图成功")
	err = util.SavePNG(proto.ImagePath+"/answer.png", answerImg)
	if err != nil {
		return fmt.Errorf("保存answer截图失败，%v", err)

	}
	log.Debugf("保存answer截图成功")
	return nil
}

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
