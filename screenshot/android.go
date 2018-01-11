package screenshot

import (
	"image"
	"os/exec"

	"github.com/silenceper/qanswer/config"
	"github.com/silenceper/qanswer/proto"
	"github.com/silenceper/qanswer/util"
)

//Android android
type Android struct{}

//NewAndroid new
func NewAndroid(cfg *config.Config) *Android {
	return new(Android)
}

//GetImage 通过adb获取截图
func (android *Android) GetImage() (img image.Image, err error) {
	err = exec.Command("adb", "shell", "screencap", "-p", "/sdcard/screenshot.png").Run()
	if err != nil {
		return
	}
	originImagePath := proto.ImagePath + "origin.png"
	err = exec.Command("adb", "pull", "/sdcard/screenshot.png", originImagePath).Run()
	if err != nil {
		return
	}
	img, err = util.OpenPNG(originImagePath)
	return
}
