package i

type Generation interface {
	GetData() ([][]string, error)
	ClearData(data string) ([]string, error)
	RefreshData(tableHead []string, datas [][]string, refreshData int) ([]any, error)
}
