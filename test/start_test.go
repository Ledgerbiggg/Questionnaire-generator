package test

import (
	"QuestionnaireDataGenerator/service"
	"testing"
)

func TestStart(t *testing.T) {
	err := service.Start()
	if err != nil {
		t.Error(err)
	}
}
