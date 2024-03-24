package main

import (
	"QuestionnaireDataGenerator/config"
	logs "QuestionnaireDataGenerator/log"
	"QuestionnaireDataGenerator/service"
)

func init() {
	logs.InitLogStyle()
	err := config.LoadConfig()
	if err != nil {
		logs.Println("load config fail " + err.Error())
	}
	err = config.ReadTitle()
	config.Configs.DataNum = config.Configs.DataNum * config.Configs.ReliabilityBias
}

func main() {
	err := service.Start()
	if err != nil {
		logs.Println("程序出错啦！！！", err)
	}
}
