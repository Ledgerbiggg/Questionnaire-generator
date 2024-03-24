package req

import (
	"QuestionnaireDataGenerator/config"
	"QuestionnaireDataGenerator/log"
	"QuestionnaireDataGenerator/model/common"
	"QuestionnaireDataGenerator/model/i"
	utils "QuestionnaireDataGenerator/utils"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"time"
)

type XunFeiAIHeader struct {
	AppID string `json:"app_id"`
	UID   string `json:"uid"`
}

type XunFeiAIChatParameter struct {
	Chat *XunFeiAIChatDomain `json:"chat"`
}

type XunFeiAIChatDomain struct {
	Domain      string  `json:"domain"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

type XunFeiAIMessage struct {
	Text []*common.MessageContent `json:"text"`
}

func newXunFeiAIMessageContent(role string, content string) *common.MessageContent {
	return &common.MessageContent{Role: role, Content: content}
}

type XunFeiAIPayload struct {
	Message *XunFeiAIMessage `json:"message"`
}

type XunFeiAIRequest struct {
	Header    *XunFeiAIHeader        `json:"header"`
	Parameter *XunFeiAIChatParameter `json:"parameter"`
	Payload   *XunFeiAIPayload       `json:"payload"`
}

func (x *XunFeiAIRequest) SetMessages(messages []*common.MessageContent) {
	x.Payload.Message.Text = messages
}
func (x *XunFeiAIRequest) ShowModel() string {
	return x.Parameter.Chat.Domain
}

func (x *XunFeiAIRequest) GetMessages() []*common.MessageContent {
	return x.Payload.Message.Text
}

func (x *XunFeiAIRequest) Chat(request i.AiModel) (i.AiModel, error) {
	log.Println(x.ShowModel() + "调用chat")
	hostUrl := "wss://spark-api.xf-yun.com/v3.5/chat"
	apiSecret := config.Configs.XunFeiConfig.ApiSecret
	apiKey := config.Configs.XunFeiConfig.ApiKey
	d := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	//握手并建立websocket 连接
	conn, resp, err := d.Dial(utils.AssembleAuthUrl(hostUrl, apiKey, apiSecret), nil)
	if err != nil {
		logs.Println("connect error:", err)
		return request, err
	} else if resp.StatusCode != 101 {
		logs.Println("connect error:", readResp(resp))
		return request, err
	}

	go func() {
		err = conn.WriteJSON(request)
		if err != nil {
			logs.Println("write message error:", err)
		}
	}()

	var answer = ""
	//获取返回的数据
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			logs.Println("read message error:", err)
			break
		}

		var data map[string]interface{}

		defer func() {
			if r := recover(); r != nil {
				logs.Println("Error occurred:", r)
				// 在出现错误时再次进行json.Marshal()
				marshal, err := json.Marshal(data)
				if err != nil {
					logs.Println("Error marshaling JSON:", err)
					return
				}
				logs.Println("========================" + string(marshal) + "结束========================")
			}
		}()

		err = json.Unmarshal(msg, &data)
		if err != nil {
			logs.Println("Error parsing JSON:", err)
			return request, err
		}
		//解析数据
		payload := data["payload"].(map[string]interface{})
		choices := payload["choices"].(map[string]interface{})
		header := data["header"].(map[string]interface{})
		code := header["code"].(float64)

		if code != 0 {
			logs.Println(data["payload"])
			return request, err
		}
		status := choices["status"].(float64)
		text := choices["text"].([]interface{})
		content := text[0].(map[string]interface{})["content"].(string)
		if status != 2 {
			answer += content
		} else {
			answer += content
			conn.Close()
			break
		}

	}
	request.AddMessage(answer)
	return request, nil
}

func (x *XunFeiAIRequest) GetLastMessage() string {
	return x.Payload.Message.Text[len(x.Payload.Message.Text)-1].Content
}

func (x *XunFeiAIRequest) AddMessage(content string) {
	x2 := &common.MessageContent{}
	if len(x.Payload.Message.Text)%2 == 0 {
		x2 = newXunFeiAIMessageContent("assistant", content)
	} else {
		x2 = newXunFeiAIMessageContent("user", content)
	}
	x.Payload.Message.Text = append(x.Payload.Message.Text, x2)
}

func NewXunFeiAIRequest(system string) *XunFeiAIRequest {
	feiConfig := config.Configs.XunFeiConfig
	newUUID, _ := uuid.NewUUID()

	header := XunFeiAIHeader{
		AppID: feiConfig.AppID,
		UID:   newUUID.String(),
	}
	chatDomain := &XunFeiAIChatDomain{
		Domain:      "generalv3.5",
		Temperature: 0.5,
		MaxTokens:   8192,
	}
	parameter := XunFeiAIChatParameter{
		Chat: chatDomain,
	}

	payload := XunFeiAIPayload{
		Message: &XunFeiAIMessage{
			Text: []*common.MessageContent{
				newXunFeiAIMessageContent("system", system),
			},
		},
	}
	return &XunFeiAIRequest{
		Header:    &header,
		Parameter: &parameter,
		Payload:   &payload,
	}
}
func readResp(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		logs.Println("read resp error:", err)
	}
	return fmt.Sprintf("code=%d,body=%s", resp.StatusCode, string(b))
}
