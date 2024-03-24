package api

import (
	"QuestionnaireDataGenerator/config"
	"QuestionnaireDataGenerator/model/req"
	"strings"
)

func InitAllApi(system string) {
	// 初始化所有的ai的api
	apis := NewAIApis()
	// 阿里云的所有模型
	aliyunModel := config.Configs.AliyunConfig.Model
	if aliyunModel != "" {
		split := strings.Split(aliyunModel, ",")
		for _, v := range split {
			apis.AddApi(req.NewAliyunRequest(v, system))
		}
	}
	// 科大讯飞模型
	if config.Configs.XunFeiConfig.ApiKey != "" {
		apis.AddApi(req.NewXunFeiAIRequest(system))
	}
	if len(apis.Apis) == 0 {
		panic("至少配置一种ai模型")
	}
	Apis = apis
}
