package service

import (
	"QuestionnaireDataGenerator/basic/data/i"
	logs "QuestionnaireDataGenerator/log"
	"QuestionnaireDataGenerator/utils"
)

type QuestionnaireService struct {
	QuestionnaireName string
	TableHead         []string
	Generation        i.Generation
	ReliabilityBias   int
}

func NewQuestionnaireService(questionnaireName string, tableHead []string, generation i.Generation, reliabilityBias int) *QuestionnaireService {
	return &QuestionnaireService{QuestionnaireName: questionnaireName, TableHead: tableHead, Generation: generation, ReliabilityBias: reliabilityBias}
}

func (q *QuestionnaireService) Generate() error {
	getData, err := q.Generation.GetData()
	if err != nil {
		logs.Println("获取问卷数据失败" + err.Error())
		return err
	}
	refreshData, err := q.Generation.RefreshData(q.TableHead, getData, q.ReliabilityBias)
	if err != nil {
		logs.Println("刷新问卷数据失败" + err.Error())
		return err
	}
	helper := utils.NewExcelHelper(q.QuestionnaireName+".xlsx", "sheet1", q.TableHead, refreshData)
	helper.SetMapOrder(q.TableHead)
	err = helper.GenerateExcel()
	if err != nil {
		logs.Println("生成问卷失败" + err.Error())
		return err
	}

	return nil

}
