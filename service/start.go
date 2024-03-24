package service

import (
	"QuestionnaireDataGenerator/basic/data"
	"QuestionnaireDataGenerator/config"
)

func Start() error {
	generation := data.NewGeneration1(config.QuestionnaireFrame.Questions, config.Configs.MinReliabilityAndValidity, config.Configs.DataNum, 6)
	var tableHead []string
	questions := config.QuestionnaireFrame.Questions
	for _, v := range questions {
		tableHead = append(tableHead, v.Title)
	}
	service := NewQuestionnaireService(config.Configs.PaperName+"问卷", tableHead, generation, config.Configs.ReliabilityBias)
	err := service.Generate()
	if err != nil {
		return err
	}
	return nil
}
