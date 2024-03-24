package utils

import "testing"

type person struct {
	Name string
	Age  int
}

func TestExcelHelper_GenerateExcel(t *testing.T) {

	i := make([]any, 3)
	i[0] = map[string]string{"a": "1", "b": "2", "c": "3"}
	i[1] = map[string]string{"a": "1", "b": "2", "c": "3"}
	i[2] = map[string]string{"a": "1", "b": "2", "c": "3"}

	helper := NewExcelHelper("test.xlsx", "Sheet1", []string{"a", "b", "c"}, i)

	helper.SetMapOrder([]string{"a", "b", "c"})
	err := helper.GenerateExcel()
	if err != nil {
		t.Error(err)
	}
}
