/*
@author: ledger
@since: 2024/2/23
*/

package utils

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"reflect"
	"strconv"
)

type ExcelHelper struct {
	FileName  string
	SheetName string
	Header    []string
	Data      []any
	MapOrder  []string
}

func (e *ExcelHelper) SetMapOrder(mapOrder []string) {
	e.MapOrder = mapOrder
}

func (e *ExcelHelper) GenerateExcel() error {
	file := excelize.NewFile()
	for i := 0; i < len(e.Header); i++ {
		column := e.intToExcelColumn(i)
		file.SetCellValue(e.SheetName, column+"1", e.Header[i])
	}
	for i := 0; i < len(e.Data); i++ {
		row := i + 2 // 行索引从2开始，因为第一行是表头
		a := e.Data[i]
		t := reflect.ValueOf(a)
		if t.Kind() == reflect.Ptr {
			t = t.Elem() // 获取指针指向的值的类型
		}
		if t.Kind() == reflect.Struct {
			// 遍历结构体的字段
			for j := 0; j < t.NumField(); j++ {
				field := t.Field(j)
				file.SetCellValue(e.SheetName, e.intToExcelColumn(j)+strconv.Itoa(row), field.Interface())
			}
		} else if t.Kind() == reflect.Map {
			for j, key := range e.MapOrder {
				value := t.MapIndex(reflect.ValueOf(key))
				// 检查值是否有效
				if value.IsValid() {
					// 处理键值对，例如输出
					column := e.intToExcelColumn(j)
					file.SetCellValue(e.SheetName, column+strconv.Itoa(row), value.Interface())
				} else {
					return errors.New("map key " + key + " not found")
				}
			}
		}

	}
	err := file.SaveAs(e.FileName)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func NewExcelHelper(fileName string, sheetName string, header []string, data []any) *ExcelHelper {
	helper := ExcelHelper{FileName: fileName, SheetName: sheetName, Header: header, Data: data}
	return &helper
}

// 将数字转换为 Excel 列的名称
func (e *ExcelHelper) intToExcelColumn(n int) string {
	// 处理特殊情况，0对应'A'
	if n == 0 {
		return "A"
	}

	result := ""
	n++
	for n > 0 {
		// 计算当前最低位对应的字母
		remainder := n % 26
		// 如果remainder为0，说明当前位应该是'Z'，而不是'A'
		if remainder == 0 {
			result = "Z" + result
			// 对应的数字需要减去26
			n = n/26 - 1
		} else {
			// 否则按照'A' + remainder - 1的逻辑计算当前位的字母
			result = fmt.Sprintf("%c", 'A'+remainder-1) + result
			n /= 26
		}
	}
	return result
}
