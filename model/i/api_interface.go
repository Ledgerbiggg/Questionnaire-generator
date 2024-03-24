package i

import (
	"QuestionnaireDataGenerator/model/common"
)

type AiModel interface {
	// GetLastMessage 获取最后一条消息
	GetLastMessage() string
	// AddMessage 添加一条消息
	AddMessage(string)
	// Chat 聊天
	Chat(AiModel) (AiModel, error)
	// GetMessages 获取所有的
	GetMessages() []*common.MessageContent
	// ShowModel 显示模型
	ShowModel() string
	// SetMessages 设置消息
	SetMessages([]*common.MessageContent)
}
