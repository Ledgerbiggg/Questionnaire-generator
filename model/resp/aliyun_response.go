package resp

type AliyunOutput struct {
	FinishReason string `json:"finish_reason"`
	Text         string `json:"text"`
}

type AliyunUsage struct {
	TotalTokens  int `json:"total_tokens"`
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type AliyunResponse struct {
	Output    AliyunOutput `json:"output"`
	Usage     AliyunUsage  `json:"usage"`
	RequestID string       `json:"request_id"`
}

func NewAliyunResponse() *AliyunResponse {
	return &AliyunResponse{}
}
