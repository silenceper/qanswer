package ocr

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/silenceper/qanswer/config"
	"github.com/silenceper/qanswer/util"
)

//Baidu baidu ocr api
type Baidu struct {
	apiKey    string
	secretKey string

	accessToken string
	sync.RWMutex
}

type accessTokenRes struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int32  `json:"expires_in"`
}

//wordsResults 匹配
type wordsResults struct {
	WordsNum    int32 `json:"words_result_num"`
	WordsResult []struct {
		Words string `json:"words"`
	} `json:"words_result"`
}

//NewBaidu new
func NewBaidu(cfg *config.Config) *Baidu {
	baidu := new(Baidu)
	baidu.apiKey = cfg.BaiduAPIKey
	baidu.secretKey = cfg.BaiduSecretKey
	return baidu
}

//GetText 识别图片中的文字
func (baidu *Baidu) GetText(imgPath string) (string, error) {
	accessToken, err := baidu.getAccessToken()
	if err != nil {
		return "", err
	}
	base64Data, err := util.OpenImageToBase64(imgPath)
	if err != nil {
		return "", err
	}
	uri := fmt.Sprintf("https://aip.baidubce.com/rest/2.0/ocr/v1/general_basic?access_token=%s", accessToken)

	postData := url.Values{}
	postData.Add("image", base64Data)
	body, err := util.PostForm(uri, postData, 5)
	if err != nil {
		return "", err
	}
	wordResults := new(wordsResults)
	err = json.Unmarshal(body, wordResults)
	if err != nil {
		return "", err
	}
	var text string
	replaceRe, _ := regexp.Compile(`^[1-9]{0,1}`)
	for k, words := range wordResults.WordsResult {
		//TIPS: 去除第一个数字 1-9
		//虽然有12个字符，但是 10-12 与最后的数字识别混在一起了
		if k == 0 {
			words.Words = replaceRe.ReplaceAllString(words.Words, "")
		}
		text = fmt.Sprintf("%s\n%s", text, strings.TrimSpace(words.Words))
	}
	text = strings.TrimLeft(text, "\n")
	return text, nil
}

func (baidu *Baidu) getAccessToken() (accessToken string, err error) {
	baidu.Lock()
	defer baidu.Unlock()
	if baidu.accessToken != "" {
		accessToken = baidu.accessToken
		return
	}
	uri := fmt.Sprintf("https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s", baidu.apiKey, baidu.secretKey)
	body, e := util.PostForm(uri, nil, 5)
	if e != nil {
		err = e
		return
	}
	res := new(accessTokenRes)
	err = json.Unmarshal(body, res)
	if err != nil {
		return
	}
	accessToken = res.AccessToken
	baidu.accessToken = accessToken
	return
}
