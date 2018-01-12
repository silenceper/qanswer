package ocr

import (
	"os/exec"

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
	body, err := exec.Command("tesseract", imgPath, "stdout", "-l", "chi_sim").Output()
	if err != nil {
		return "", err
	}
	return string(body), nil
}
