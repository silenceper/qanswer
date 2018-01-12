package config

import (
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v1"
)

//Config 全局配置
type Config struct {
	Debug      bool   `yaml:"debug"`
	Device     string `yaml:"device"`
	OcrType    string `yaml:"ocr_type"`
	WdaAddress string `yaml:"wda_address"`

	//baidu ocr
	BaiduAPIKey    string `yaml:"baidu_api_key"`
	BaiduSecretKey string `yaml:"baidu_secret_key"`

	//截图题目位置
	QuestionX int `yaml:"question_x"`
	QuestionY int `yaml:"question_y"`
	QuestionW int `yaml:"question_w"`
	QuestionH int `yaml:"question_h"`

	//截取答案位置
	AnswerX int `yaml:"answer_x"`
	AnswerY int `yaml:"answer_y"`
	AnswerW int `yaml:"answer_w"`
	AnswerH int `yaml:"answer_h"`
}

var cfg *Config

var cfgFilename = "./config.yml"

//SetConfigFile 设置配置文件地址
func SetConfigFile(path string) {
	cfgFilename = path
}

//GetConfig 解析配置
func GetConfig() *Config {
	if cfg != nil {
		return cfg
	}
	filename, _ := filepath.Abs(cfgFilename)
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}
	var c *Config
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		panic(err)
	}
	cfg = c
	return cfg
}
