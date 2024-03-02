package pkg

import (
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	"os"

	"text/template"
)

const (
	fsPrefix       = "files/"
	pageTemplateFn = "template.html"

	tableTemplateFn       = "template-table.html"
	columnChartTemplateFn = "template-column.html"
)

var (
	//go:embed files/template*
	templateFS embed.FS

	//go:embed files/charts.min.css
	chartMinCSS string

	//go:embed files/readstat.css
	readstatCSS string
)

type chartTemplateData struct {
	Title string

	ChartCSS    string
	ReadstatCSS string

	BookReadCountTable string
	BookReadCountChart string
}

type columnChartData struct {
	ID        string
	Type      string // e.g. column or bar
	Caption   string
	RowLabel  string
	DataLabel string
	Rows      []columnChartDataRow
}

type tableData struct {
	ID        string
	Type      string // e.g. column or bar
	Caption   string
	RowLabel  string
	DataLabel string
	Rows      []tableDataRow
}

type columnChartDataRow struct {
	Label string
	Data  int
}

type tableDataRow struct {
	Label string
	Data  string
}

const (
	titleFmt = "Reading Stats for %d"
)

// NewChart read the stats and write single page html to filename
func NewChart(stats YearStats, year int, filename string) error {
	var err error

	// Build data!
	data := chartTemplateData{
		ChartCSS:    chartMinCSS,
		ReadstatCSS: readstatCSS,

		Title: fmt.Sprintf(titleFmt, year),
	}

	// Books finished month
	bookReadCount := []int{}
	bookReadNames := []string{}

	months := []string{"", "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	for idx := 1; idx <= 12; idx++ {
		bookReadCount = append(bookReadCount, len(stats.BooksFinishedMonth(idx)))

		names := ""
		for _, finishedBook := range stats.BooksFinishedMonth(idx) {
			if names != "" {
				names = names + ", "
			}

			names += fmt.Sprintf("%s: %s", finishedBook.Title, finishedBook.Author)
		}

		bookReadNames = append(bookReadNames, names)
	}

	data.BookReadCountTable, err = renderTable("book-read-table", "Finished Books by Month", "Month", "Count", months[1:], bookReadNames)
	if err != nil {
		return err
	}

	data.BookReadCountChart, err = renderColumnChart("book-read-count", "Finished Books by Month", "Month", "Count", months[1:], bookReadCount)
	if err != nil {
		return err
	}

	// Write template!
	fp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}

	defer func() {
		_ = fp.Close()
	}()

	tmpl, err := template.New(pageTemplateFn).ParseFS(templateFS, fsPrefix+pageTemplateFn)
	if err != nil {
		return err
	}

	err = tmpl.Execute(fp, data)
	if err != nil {
		return err
	}

	return nil
}

func renderTable(id, title, rowLabel, dataLabel string, rows, datas []string) (string, error) {
	if len(rows) != len(datas) {
		return "", fmt.Errorf("length of rows must match length of datas %d != %d", len(rows), len(datas))
	}

	data := tableData{
		ID:        id,
		Caption:   title,
		RowLabel:  rowLabel,
		DataLabel: dataLabel,
		Rows:      make([]tableDataRow, 0),
	}

	for idx := range rows {
		data.Rows = append(data.Rows, tableDataRow{Label: rows[idx], Data: datas[idx]})
	}

	buf := []byte{}
	bufWriter := bytes.NewBuffer(buf)

	tmpl, err := template.New(tableTemplateFn).ParseFS(templateFS, fsPrefix+tableTemplateFn)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(bufWriter, data)
	if err != nil {
		return "", err
	}

	return bufWriter.String(), nil
}

func renderColumnChart(id, title, rowLabel, dataLabel string, rows []string, datas []int) (string, error) {
	return renderTableOrChart(columnChartTemplateFn, "bar", id, title, rowLabel, dataLabel, rows, datas)
}

func renderTableOrChart(templateFn, typeStr, id, title, rowLabel, dataLabel string, rows []string, datas []int) (string, error) {
	if len(rows) != len(datas) {
		return "", fmt.Errorf("length of rows must match length of datas %d != %d", len(rows), len(datas))
	}

	data := columnChartData{
		Type:      typeStr,
		ID:        id,
		Caption:   title,
		RowLabel:  rowLabel,
		DataLabel: dataLabel,
		Rows:      make([]columnChartDataRow, 0),
	}

	for idx := range rows {
		data.Rows = append(data.Rows, columnChartDataRow{Label: rows[idx], Data: datas[idx]})
	}

	buf := []byte{}
	bufWriter := bytes.NewBuffer(buf)

	tmpl, err := template.New(templateFn).ParseFS(templateFS, fsPrefix+templateFn)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(bufWriter, data)
	if err != nil {
		return "", err
	}

	return bufWriter.String(), nil
}
