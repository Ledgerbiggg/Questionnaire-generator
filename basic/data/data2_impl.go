package data

import (
	"QuestionnaireDataGenerator/config"
	logs "QuestionnaireDataGenerator/log"
	"QuestionnaireDataGenerator/utils"
	"errors"
	"github.com/google/uuid"
	"strings"
)

type Generation1 struct {
	Questions                 []*config.Question `json:"questions"`
	MinReliabilityAndValidity float64            `json:"minReliabilityAndValidity"`
	DataNum                   int                `json:"dataNum"`
	ThreadChan                chan int           `json:"thread"`
}

func NewGeneration1(questions []*config.Question, minReliabilityAndValidity float64, dataNum int, thread int) *Generation1 {
	c := make(chan int, thread)
	for i := 0; i < thread; i++ {
		c <- 1
	}
	return &Generation1{Questions: questions, MinReliabilityAndValidity: minReliabilityAndValidity, DataNum: dataNum, ThreadChan: c}
}

func (g *Generation1) GetData() ([][]string, error) {
	datas := make([][]string, 0)
	for _, q := range g.Questions {
		var data []string

		if q.IsLiKeTe {
			data = utils.GenerateRandomArray(g.DataNum)
		} else if q.IsMulti {
			if q.Options == nil || len(q.Options) == 0 {
				logs.Println(q.Title, "这个的options配置项配置不正确!!!,多选题必须配置选择的项目option!!!!")
				return nil, errors.New("这个的options配置项配置不正确!!!,多选题必须配置选择的项目option")
			}
			for i := 0; i < g.DataNum; i++ {
				subset := utils.RandomSubset(q.Options)
				join := strings.Join(subset, ",")
				data = append(data, join)
			}
		} else if q.IsFill {
			for i := 0; i < g.DataNum; i++ {
				data = append(data, uuid.New().String())
			}
		} else {
			if q.Options == nil || len(q.Options) == 0 {
				logs.Println(q.Title, "这个的options配置项配置不正确!!!,单选题必须配置选择的项目option!!!!")
				return nil, errors.New("这个的options配置项配置不正确!!!,单选题必须配置选择的项目option")
			}
			for i := 0; i < g.DataNum; i++ {
				randomString := utils.RandomString(q.Options)
				data = append(data, randomString)
			}
		}
		datas = append(datas, data)
	}
	return datas, nil
}

func (g *Generation1) ClearData(data string) ([]string, error) {
	return nil, nil
}

func (g *Generation1) RefreshData(tableHead []string, datas [][]string, reliabilityBias int) ([]any, error) {
	if tableHead == nil || len(tableHead) == 0 || datas == nil || len(datas) == 0 {
		return nil, errors.New("表头和数据不能为空")
	}
	if len(tableHead) != len(datas) {
		return nil, errors.New("表头和数据长度不一致")
	}
	// 最终数据
	for i := 0; i < len(datas); i++ {
		if len(datas[i]) > 10 {
			datas[i] = datas[i][:len(datas[i])/reliabilityBias]
		}
	}
	var tableData []any

	var count = 0
	for count < len(datas[0]) {
		var mapData = make(map[string]string)
		for j, v := range tableHead {
			mapData[v] = datas[j][count]
		}
		tableData = append(tableData, mapData)
		count++
	}

	return tableData, nil
}
