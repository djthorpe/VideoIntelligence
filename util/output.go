package util

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

// Output defines the table of output data
type Output struct {
	columns []string
	rows    []*row
}

type row struct {
	values map[string]string
}

// NewOutput returns an output object
func NewOutput(columns ...string) *Output {
	this := new(Output)
	this.columns = make([]string, len(columns))
	this.rows = make([]*row, 0, 0)
	for i, column := range columns {
		this.columns[i] = column
	}
	return this
}

// AddColumns
func (this *Output) AddColumns(columns ...string) {
	this.columns = append(this.columns, columns...)
}

// AppendMap appends a row to the table
func (this *Output) AppendMap(row map[string]interface{}) {
	r := this.newRow()
	for k, v := range row {
		r.set(k, fmt.Sprintf("%v", v))
	}
}

func (this *Output) RenderASCII() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(this.columns)
	for _, r := range this.rows {
		table.Append(r.row(this.columns))
	}
	table.Render()
}

////////////////////////

func (this *Output) newRow() *row {
	r := new(row)
	r.values = make(map[string]string, 0)
	this.rows = append(this.rows, r)
	return r
}

func (this *row) set(key string, value string) {
	this.values[key] = value
}

func (this *row) row(columns []string) []string {
	row := make([]string, len(columns))
	for i, k := range columns {
		v, exists := this.values[k]
		if exists {
			row[i] = v
		} else {
			row[i] = "<nil>"
		}
	}
	return row
}
