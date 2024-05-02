package pkg

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"text/template"
	"time"
)

const (
	fsPrefix       = "files/"
	pageTemplateFn = "template.html"
)

var (
	//go:embed files/template*
	templateFS embed.FS

	//go:embed files/readstat.css
	readstatCSS string
)

type chartTemplateData struct {
	Title       string
	Description string

	BooksFinished   []chartTemplateStat
	BooksRead       int
	BookReadingTime string
	BookReadCount   string

	ArticlesFinished   []chartTemplateStat
	ArticlesRead       int
	ArticleReadingTime string
	ArticleReadCount   string

	TotalReadingTime string

	ReadstatCSS string
}

type chartTemplateStat struct {
	Title    string
	URL      string
	Author   string
	Duration string
	Sessions int
	Month    string
}

const (
	titleFmt = "Reading Stats for %d"
)

// NewChart read the stats and write single page html to filename
func NewChart(stats Stats, year int, filename string) error {
	var err error

	months := []string{"", "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}

	data := chartTemplateData{
		ReadstatCSS: readstatCSS,

		Title: fmt.Sprintf(titleFmt, year),
		// Description: "reading stats todo description",

		BooksRead:    len(stats.BooksFinishedYear(year)),
		ArticlesRead: len(stats.ArticlesFinishedYear(year)),
	}

	bookReadCount := []int{}

	articleReadCount := []int{}

	totalReadTime := []float32{}

	for idx := 1; idx <= 12; idx++ {
		finBooks := stats.BooksFinishedMonth(year, idx)
		for jdx := range finBooks {
			totalReadTime = append(totalReadTime, float32(finBooks[jdx].ReadSeconds()))

			data.BooksFinished = append(data.BooksFinished, chartTemplateStat{
				Title:    finBooks[jdx].Title,
				Author:   finBooks[jdx].Author,
				Duration: HumanizeDurationShort(time.Second * time.Duration(finBooks[jdx].ReadSeconds())),
				Sessions: finBooks[jdx].NumSessions(),
				Month:    months[idx],
			})
		}

		finArts := stats.ArticlesFinishedMonth(year, idx)
		for jdx := range finArts {
			totalReadTime[idx-1] += float32(finArts[jdx].ReadSeconds())

			data.ArticlesFinished = append(data.ArticlesFinished, chartTemplateStat{
				Title:    finArts[jdx].Title,
				URL:      finArts[jdx].URL,
				Duration: HumanizeDurationShort(time.Second * time.Duration(finArts[jdx].ReadSeconds())),
				Sessions: finArts[jdx].NumSessions(),
				Month:    months[idx],
			})
		}

		bookReadCount = append(bookReadCount, len(stats.BooksFinishedMonth(year, idx)))
		articleReadCount = append(articleReadCount, len(stats.ArticlesFinishedMonth(year, idx)))
	}

	data.BookReadingTime = HumanizeDurationShort(time.Second * time.Duration(stats.BooksSecondsReadYear(year)))
	data.ArticleReadingTime = HumanizeDurationShort(time.Second * time.Duration(stats.ArticlesSecondsReadYear(year)))

	bC, _ := json.Marshal(bookReadCount)
	data.BookReadCount = string(bC)

	aC, _ := json.Marshal(articleReadCount)
	data.ArticleReadCount = string(aC)

	tC, _ := json.Marshal(totalReadTime)
	data.TotalReadingTime = string(tC)

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
