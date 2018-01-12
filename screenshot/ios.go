package screenshot

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"

	"github.com/silenceper/qanswer/config"
	"github.com/silenceper/qanswer/util"
)

//IOS 获取iOS截图
type IOS struct {
	wdaAddress string
}

type screenshotRes struct {
	Value     string `json:"value"`
	SessionID string `json:"sessionId"`
	Status    int    `json:"status"`
}

//NewIOS new
func NewIOS(cfg *config.Config) *IOS {
	ios := new(IOS)
	ios.wdaAddress = cfg.WdaAddress
	if ios.wdaAddress == "" {
		panic("请指定 wda 连接地址")
	}
	return ios
}

//GetImage 返回图片生成的路径
func (ios *IOS) GetImage() (img image.Image, err error) {
	body, e := util.HTTPGet(fmt.Sprintf("http://%s/screenshot", ios.wdaAddress), 3)
	if e != nil {
		err = fmt.Errorf("WebDriverAgentRunner 连接失败, err=%v", e)
		return
	}

	res := new(screenshotRes)
	e = json.Unmarshal(body, res)
	if err != nil {
		err = fmt.Errorf("WebDriverAgentRunner 响应数据异常，请检查 WebDriverAgentRunner 运行状态, err=%v", e)
		return
	}
	pngValue, e := base64.StdEncoding.DecodeString(res.Value)
	if err != nil {
		err = fmt.Errorf("图片解码失败, err=%v", e)
		return
	}

	src, err := png.Decode(bytes.NewReader(pngValue))
	if err != nil {
		err = fmt.Errorf("图片解码失败, err=%v", e)
		return
	}
	img = src
	return
}
