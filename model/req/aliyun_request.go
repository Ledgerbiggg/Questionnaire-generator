package req

import (
	"QuestionnaireDataGenerator/config"
	"QuestionnaireDataGenerator/model/common"
	"QuestionnaireDataGenerator/model/i"
	"QuestionnaireDataGenerator/model/resp"
	utils "QuestionnaireDataGenerator/utils"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"strings"
)

func NewAliyunInputMessage(role string, content string) *common.MessageContent {
	return &common.MessageContent{Role: role, Content: content}
}

type AliyunInput struct {
	Messages []*common.MessageContent `json:"messages"`
}

type AliyunRequest struct {
	HostUrl string       `json:"hostUrl"`
	Model   string       `json:"model"`
	Input   *AliyunInput `json:"input"`
}

func (a *AliyunRequest) ShowModel() string {
	return a.Model
}

func (a *AliyunRequest) SetMessages(messages []*common.MessageContent) {
	a.Input.Messages = messages
}

func (a *AliyunRequest) GetMessages() []*common.MessageContent {
	return a.Input.Messages
}

func (a *AliyunRequest) Chat(request i.AiModel) (i.AiModel, error) {
	log.Println(a.ShowModel() + "调用chat")
	headers := map[string]string{}
	headers["Content-Type"] = "application/json"
	headers["Accept"] = "text/event-stream"
	headers["Authorization"] = "Bearer " + config.Configs.AliyunConfig.Authorization
	headers["X-DashScope-SSE"] = "enable"
	marshal, err := json.Marshal(request)
	if err != nil {
		log.Println("Error reading request body:", err)
		return request, err
	}
	post, err := utils.NewHttpDos(
		a.HostUrl,
		nil,
		marshal,
		headers,
	).Post()

	if err != nil {
		log.Println("Error reading request body:", err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(post))
	var allLine []string
	for scanner.Scan() {
		line := scanner.Text()
		allLine = append(allLine, line)
	}
	if err = scanner.Err(); err != nil {
		log.Println("Error reading request body:", err)
		return request, err
	}

	if len(allLine) <= 2 {
		return request, errors.New("received less than four lines of data")
	} else {
		resline := len(allLine) - 2
		replace := strings.Replace(allLine[resline], "data:", "", 1)
		response := resp.NewAliyunResponse()
		err = json.Unmarshal([]byte(replace), response)
		if err != nil {
			log.Println("Error reading request body:", err)
			return request, err
		}
		// 消息发送成功
		request.AddMessage(response.Output.Text)
		return request, err
	}
}

func (a *AliyunRequest) GetLastMessage() string {
	return a.Input.Messages[len(a.Input.Messages)-1].Content
}

func (a *AliyunRequest) AddMessage(content string) {
	a2 := &common.MessageContent{}
	if len(a.Input.Messages)%2 == 0 {
		a2 = NewAliyunInputMessage("assistant", content)
	} else {
		a2 = NewAliyunInputMessage("user", content)
	}
	a.Input.Messages = append(a.Input.Messages, a2)
}

func NewAliyunRequest(
	model string,
	system string,
) *AliyunRequest {
	input := &AliyunInput{}
	input.Messages = append(input.Messages, NewAliyunInputMessage("system", system))
	return &AliyunRequest{HostUrl: "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation", Model: model, Input: input}
}
