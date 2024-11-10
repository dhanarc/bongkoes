package console

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"reflect"
)

func GetHeaders[T any](model T) []string {
	var headers []string
	modelStruct := reflect.TypeOf(model)
	for i := 0; i < modelStruct.NumField(); i++ {
		field := modelStruct.Field(i)

		headerTag := field.Tag.Get("header")
		headers = append(headers, headerTag)
	}
	return headers
}

func PrintTable[T any](body []T) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	headers := GetHeaders(*new(T))
	var headerRow table.Row
	headerRow = append(headerRow, "#")
	for i := range headers {
		headerRow = append(headerRow, headers[i])
	}
	t.AppendHeader(headerRow)

	// body
	var rows []table.Row
	sequence := 1
	for i := range body {
		rowContent := reflect.ValueOf(body[i])

		var row []interface{}
		row = append(row, sequence)
		for i := 0; i < rowContent.NumField(); i++ {
			row = append(row, rowContent.Field(i).Interface())
		}

		rows = append(rows, row)
		sequence++
	}
	t.AppendRows(rows)
	t.AppendSeparator()
	t.Render()
}
