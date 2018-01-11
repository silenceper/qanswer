package ocr

import (
	"github.com/otiai10/gosseract"
	"github.com/silenceper/qanswer/config"
)

//Tesseract tesseract 识别
type Tesseract struct{}

//NewTesseract new
func NewTesseract(cfg *config.Config) *Tesseract {
	return new(Tesseract)
}

//GetText 根据图片路径获取识别文字
func (tesseract *Tesseract) GetText(imgPath string) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	client.SetImage(imgPath)
	client.SetLanguage("chi_sim")

	return client.Text()
}
