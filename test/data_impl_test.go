package test

import (
	"QuestionnaireDataGenerator/basic/data"
	"QuestionnaireDataGenerator/config"
	logs "QuestionnaireDataGenerator/log"
	"testing"
)

func init() {
	logs.InitLogStyle()
	err := config.LoadConfig()
	if err != nil {
		logs.Println("load config fail " + err.Error())
	}
	err = config.ReadTitle()
}

func TestGeneration_GetData(t *testing.T) {
	generation := data.NewGeneration1(config.QuestionnaireFrame.Questions, config.Configs.MinReliabilityAndValidity, config.Configs.DataNum, 2)
	getData, err := generation.GetData()
	if err != nil {
		t.Error(err)
	}

	t.Log(getData)
}
