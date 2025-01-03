package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/snabb/isoweek"
	"github.com/timchurchard/kobo-readstat/pkg"
)

func Goals(out io.Writer) int {
	const ()

	var (
		storageFn string
		year      int
	)

	flag.StringVar(&storageFn, "storage", defaultStorage, usageStoragePath)
	flag.StringVar(&storageFn, "s", defaultStorage, usageStoragePath)

	flag.IntVar(&year, "year", defaultYear, usageYear)
	flag.IntVar(&year, "y", defaultYear, usageYear)

	flag.Usage = func() {
		fmt.Fprintf(out, "Usage of %s %s:\n", os.Args[0], os.Args[1])

		flag.PrintDefaults()
	}

	flag.Parse()

	if _, err := os.Stat(storageFn); err != nil {
		panic(fmt.Sprintf("storage not found: %v", err))
	}

	storage, err := pkg.OpenStorageOrCreate(storageFn)
	if err != nil {
		panic(err)
	}

	stats := pkg.NewStats(storage)

	hoursPerWeek := make([]float32, 0)
	totalBooks := 0
	totalBookSeconds := 0
	totalArticles := 0
	totalArticleSeconds := 0

	for week := 1; week <= 53; week++ {
		weekStartTime := isoweek.StartTime(year, week, time.UTC)

		weekStartYear, _ := weekStartTime.ISOWeek()
		if weekStartYear != year {
			break
		}

		totalBooks += len(stats.Years[year].Weeks[week].FinishedBooks)
		weekBooksRead := len(stats.Years[year].Weeks[week].Books)
		weekBookSeconds := stats.BooksSecondsReadWeek(year, week)
		totalBookSeconds += weekBookSeconds

		totalArticles += len(stats.Years[year].Weeks[week].FinishedArticles)
		weekArticleRead := len(stats.Years[year].Weeks[week].Articles)
		weekArticleSeconds := stats.ArticleSecondsReadWeek(year, week)
		totalArticleSeconds += weekArticleSeconds

		weekTotalSeconds := weekBookSeconds + weekArticleSeconds

		fmt.Printf("Week %02d (%s) - Total %s hours  \t- Books %d in %s hours  \t- Articles %d in %s hours\n",
			week, weekStartTime.Format("2006-01-02"), pkg.SecondsToHoursString(weekTotalSeconds),
			weekBooksRead, pkg.SecondsToHoursString(weekBookSeconds),
			weekArticleRead, pkg.SecondsToHoursString(weekArticleSeconds))

		if weekStartTime.Unix() > time.Now().Unix() {
			break
		}

		hoursPerWeek = append(hoursPerWeek, float32(weekTotalSeconds)/3600.0)
	}

	fmt.Println("---")
	fmt.Printf("Totals!\t\t     - Hours %s  \t\t- %d Books in %s hours   \t- %d Articles in %s hours",
		pkg.SecondsToHoursString(totalBookSeconds+totalArticleSeconds),
		totalBooks, pkg.SecondsToHoursString(totalBookSeconds), totalArticles, pkg.SecondsToHoursString(totalArticleSeconds))
	fmt.Printf("\t- Weekly average hours: %.02f\n", calculateAverage(hoursPerWeek))

	return 0
}

func calculateAverage(list []float32) float32 {
	total := float32(0.0)

	for _, num := range list {
		total += num
	}

	return total / float32(len(list))
}
