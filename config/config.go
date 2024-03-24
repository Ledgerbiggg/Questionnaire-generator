package config

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

// Configs 使用全局的配置变量
var (
	Configs            = &QConfig{}
	QuestionnaireFrame = &Questions{}
)

// LoadConfig viper读取yaml
func LoadConfig() error {
	// yaml
	vconfig := viper.New()
	//表示 先预加载匹配的环境变量
	vconfig.AutomaticEnv()
	//设置环境变量分割符，将点号和横杠替换为下划线
	vconfig.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	// 设置读取的配置文件
	vconfig.SetConfigName("config")
	// 添加读取的配置文件路径
	vconfig.AddConfigPath(".")
	// 读取文件类型
	vconfig.SetConfigType("yaml")

	err := vconfig.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	var cng QConfig
	if err := vconfig.Unmarshal(&cng); err != nil {
		log.Panicln("unmarshal cng file fail " + err.Error())
	}
	// 赋值全局变量
	Configs = &cng
	return err
}

func init() {
	err := LoadConfig()
	if err != nil {
		log.Println("load config fail " + err.Error())
	}
}

// QConfig 系统整体配置
type QConfig struct {
	PaperName       string `yaml:"paperName"`
	ReliabilityBias int    `yaml:"reliabilityBias"`
	// 问卷描述
	Description string `yaml:"description"`
	// 问卷数据量
	DataNum int `yaml:"dataNum"`
	// 最低信效度
	MinReliabilityAndValidity float64 `yaml:"minReliabilityAndValidity"`
	// ai模型
	XunFeiConfig XunFeiConfig `mapstructure:"xunFei"`
	AliyunConfig AliyunConfig `mapstructure:"aliyun"`
}
type XunFeiConfig struct {
	ApiKey    string `yaml:"apiKey"`
	ApiSecret string `yaml:"apiSecret"`
	AppID     string `yaml:"appID"`
}

type AliyunConfig struct {
	Model         string `yaml:"model"`
	Authorization string `yaml:"authorization"`
}

func ReadTitle() error {
	file, err := os.OpenFile("title.json", os.O_RDONLY, 0)
	if err != nil {
		log.Println("open file fail " + err.Error())
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(QuestionnaireFrame)
	if err != nil {
		log.Println("decoder fail " + err.Error())
		return err
	}
	return nil
}
