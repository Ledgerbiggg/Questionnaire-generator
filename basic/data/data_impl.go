package data

import (
	"QuestionnaireDataGenerator/api"
	"QuestionnaireDataGenerator/config"
	logs "QuestionnaireDataGenerator/log"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Generation struct {
	Questions                 []*config.Question `json:"questions"`
	MinReliabilityAndValidity float64            `json:"minReliabilityAndValidity"`
	DataNum                   int                `json:"dataNum"`
	ThreadChan                chan int           `json:"thread"`
}

func NewGeneration(questions []*config.Question, minReliabilityAndValidity float64, dataNum int, thread int) *Generation {
	c := make(chan int, thread)
	for i := 0; i < thread; i++ {
		c <- 1
	}
	return &Generation{Questions: questions, MinReliabilityAndValidity: minReliabilityAndValidity, DataNum: dataNum, ThreadChan: c}
}

func (g *Generation) GetData() ([][]string, error) {
	datas := make([][]string, 0)
	apis := api.Apis
	for _, q := range g.Questions {
		<-g.ThreadChan
		go func(q *config.Question) {
			defer func() {
				g.ThreadChan <- 1
			}()
			for i := 0; i < 5; i++ {
				msg := "我要这个问卷的题目是:" + q.Title + "然后你给我这个问卷生成数据,你给我返回一个json的数组的数据" +
					"类似于:['问题的第一个回答','问题的第二个回答'],返回格式一定不能错" +
					"类似于:['问题的第一个回答','问题的第二个回答'],返回格式一定不能错" +
					"类似于:['问题的第一个回答','问题的第二个回答'],返回格式一定不能错"
				if q.IsLiKeTe {
					msg += "我这个是一个李克特量表,返回的数据必须是数字1-5中的其中一个,然后信效度不能低于" + fmt.Sprintf("%f", g.MinReliabilityAndValidity)
				} else if q.IsMulti {
					if q.Options == nil || len(q.Options) == 0 {
						logs.Println(q.Title, "这个的options配置项配置不正确!!!,多选题必须配置选择的项目option!!!!")
						return
					}
					marshal, err := json.Marshal(q.Options)
					if err != nil {
						logs.Println(q.Title, "这个的配置的options配置项不正确!!!,多选题必须配置选择的项目option!!!!")
						return
					}
					msg += "我这个是一个多选题,选项在" + string(marshal) + "中选一个或者多个,多个的时候使用逗号分割所有的选项,批量返回一个字符串数组"
				} else if q.IsFill {
					msg += "我这个是一个填空题,你给我批量返回一个回答这个题目的字符串数组"
				} else {
					if q.Options == nil || len(q.Options) == 0 {
						logs.Println(q.Title, "这个的options配置项配置不正确!!!,多选题必须配置选择的项目option!!!!")
						return
					}
					marshal, err := json.Marshal(q.Options)
					if err != nil {
						logs.Println(q.Title, "这个的配置的options配置项不正确!!!,单选题必须配置选择的项目option!!!!")
						return
					}
					msg += "我这个是一个单选题,选项在" + string(marshal) + "中选一个!!,批量返回这个问题问题答案的一个字符串数组"
				}
				msg += "我需要的回答的答案的数组长度为" + fmt.Sprintf("%d", g.DataNum)
				data, err := apis.Chat(msg)
				if err != nil {
					logs.Println("获取数据失败:" + err.Error())
					logs.Println("尝试重试中...")
					time.Sleep(300 * time.Millisecond)
					continue
				}
				clearData, err := g.ClearData(data)
				if err != nil {
					continue
				}
				datas = append(datas, clearData)
				break
			}
		}(q)
	}
	return datas, nil
}

func (g *Generation) ClearData(data string) ([]string, error) {
	// 定义切割函数
	splitFunc := func(r rune) bool {
		return strings.ContainsAny(string(r), "[") || strings.ContainsAny(string(r), "]")
	}
	fieldsFunc := strings.FieldsFunc(data, splitFunc)
	if len(fieldsFunc) != 3 {
		return nil, errors.New("数据格式不正确")
	}
	s := fieldsFunc[1]
	split := strings.Split(s, ",")
	for i, v := range split {
		all := strings.ReplaceAll(v, "\"", "")
		strings.ReplaceAll(all, "'", "")
		split[i] = all
	}
	return split, nil
}

func (g *Generation) RefreshData(tableHead []string, datas [][]string) ([]any, error) {
	if tableHead == nil || len(tableHead) == 0 || datas == nil || len(datas) == 0 {
		return nil, errors.New("表头和数据不能为空")
	}
	if len(tableHead) != len(datas) {
		return nil, errors.New("表头和数据长度不一致")
	}
	// 最终数据
	var tableData []any
	for i, v := range tableHead {
		var mapData = make(map[string]string)
		for _, d := range datas {
			mapData[v] = d[i]
		}
		tableData = append(tableData, mapData)
	}
	return tableData, nil
}
