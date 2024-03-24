package api

import (
	"QuestionnaireDataGenerator/model/common"
	"QuestionnaireDataGenerator/model/i"
	"sync"
)

var Apis *AIApis

type AIApis struct {
	lock           sync.Mutex
	CurrentAiModel int         `json:"currentAiModel"`
	Apis           []i.AiModel `json:"apis"`
	Messages       []*common.MessageContent
}

func (a *AIApis) Chat(message string) (string, error) {
	api := a.GetApi()
	if a.Messages != nil {
		api.SetMessages(a.Messages)
	}
	api.AddMessage(message)
	chat, err := api.Chat(api)
	if err != nil {
		return "", err
	}
	a.Messages = chat.GetMessages()
	return chat.GetLastMessage(), nil
}

func (a *AIApis) AddApi(apis ...i.AiModel) {
	a.Apis = append(a.Apis, apis...)
}

// GetApi 轮训所有的api
func (a *AIApis) GetApi() i.AiModel {
	defer func() {
		// 循环
		a.CurrentAiModel = (a.CurrentAiModel + 1) % len(a.Apis)
		defer a.lock.Unlock()
	}()
	a.lock.Lock()
	return a.Apis[a.CurrentAiModel]
}

func NewAIApis() *AIApis {
	a := &AIApis{}
	a.Apis = make([]i.AiModel, 0)
	return a
}
