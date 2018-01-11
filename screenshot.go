package qanswer

import (
	"image"

	"github.com/silenceper/qanswer/config"
	"github.com/silenceper/qanswer/proto"
	"github.com/silenceper/qanswer/screenshot"
)

//Screenshot 获取屏幕截图
type Screenshot interface {
	GetImage() (image.Image, error)
}

//NewScreenshot new
func NewScreenshot(cfg *config.Config) Screenshot {
	if cfg.Device == proto.DeviceiOS {
		return screenshot.NewIOS(cfg)
	}
	return screenshot.NewAndroid(cfg)
}
