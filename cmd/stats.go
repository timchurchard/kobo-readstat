package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/timchurchard/kobo-readstat/pkg"
)

// Stats command reads local storage and produces stats
func Stats(out io.Writer) int {
	const (
		defaultEmpty   = ""
		defaultStorage = "./readstat.json"
		defaultYear    = 2024

		usageStoragePath = "Path to local storage default: " + defaultStorage
		usageMode        = "Mode html or text (default text)"
		usageYear        = "Year to generate stats for (default 2023)"
		usageOutPath     = "Path to output file (required for mode html)"

		usageShowBooks     = "Show book title and details"
		usageShowArticles  = "Show pocket article title and details"
		usageHideArticles  = "Hide pocket articles even time spent"
		usageShowSessions  = "Show reading sessions"
		usageShowBookEnds  = "Show book started and finished reading"
		usageShowBookmarks = "Show book annotations, notes and highlights"
	)
	var (
		storageFn     string
		mode          string
		year          int
		outFn         string
		showBooks     bool
		showArticles  bool
		hideArticles  bool
		showSessions  bool
		showBookEnds  bool
		showBookmarks bool
	)

	flag.StringVar(&mode, "mode", defaultEmpty, usageMode)
	flag.StringVar(&mode, "m", defaultEmpty, usageMode)

	flag.StringVar(&storageFn, "storage", defaultStorage, usageStoragePath)
	flag.StringVar(&storageFn, "s", defaultStorage, usageStoragePath)

	flag.StringVar(&outFn, "out", defaultEmpty, usageOutPath)
	flag.StringVar(&outFn, "o", defaultEmpty, usageOutPath)

	flag.IntVar(&year, "year", defaultYear, usageYear)
	flag.IntVar(&year, "y", defaultYear, usageYear)

	flag.BoolVar(&showBooks, "showbooks", true, usageShowBooks)

	flag.BoolVar(&showArticles, "showarticles", false, usageShowArticles)
	flag.BoolVar(&hideArticles, "hidearticles", false, usageHideArticles)

	flag.BoolVar(&showSessions, "showsessions", false, usageShowSessions)
	flag.BoolVar(&showBookEnds, "showbookends", false, usageShowBookEnds)
	flag.BoolVar(&showBookmarks, "showbookmarks", false, usageShowBookmarks)

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

	booksReadSeconds := stats.BooksSecondsReadYear(year)
	booksReadDuration, _ := time.ParseDuration(fmt.Sprintf("%ds", booksReadSeconds))
	articlesReadSeconds := stats.ArticlesSecondsReadYear(year)
	articlesReadDuration, _ := time.ParseDuration(fmt.Sprintf("%ds", articlesReadSeconds))

	totalReadSeconds := booksReadSeconds + articlesReadSeconds
	totalReadDuration, _ := time.ParseDuration(fmt.Sprintf("%ds", totalReadSeconds))

	switch strings.ToLower(mode) {
	case "html":
		if outFn == "" {
			fmt.Fprintf(out, "--out -o is required for mode html")
			return 1
		}

		err = pkg.NewChart(stats, year, outFn)
		if err != nil {
			panic(err)
		}

	case "text":
		fallthrough
	default:
		fmt.Printf("Year: %d\n", year)
		fmt.Printf("Finished books\t\t\t: %d\n", len(stats.BooksFinishedYear(year)))
		if !hideArticles {
			fmt.Printf("Finished articles\t\t: %d\n", len(stats.ArticlesFinishedYear(year)))
		}

		fmt.Printf("Time reading books\t\t: %s (hours: %s)\n", pkg.HumanizeDuration(booksReadDuration), pkg.SecondsToHoursString(booksReadSeconds))
		if !hideArticles {
			fmt.Printf("Time reading articles\t\t: %s (hours: %s)\n", pkg.HumanizeDuration(articlesReadDuration), pkg.SecondsToHoursString(articlesReadSeconds))
		}

		if !hideArticles {
			fmt.Printf("Total time reading\t\t: %s (hours: %s)\n", pkg.HumanizeDuration(totalReadDuration), pkg.SecondsToHoursString(totalReadSeconds))
		}

		fmt.Println("\n----------")

		months := []string{"", "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
		for idx := 1; idx <= 12; idx++ {
			monthBookReadDuration, _ := time.ParseDuration(fmt.Sprintf("%ds", stats.BooksSecondsReadMonth(year, idx)))
			monthArticleReadDuration, _ := time.ParseDuration(fmt.Sprintf("%ds", stats.ArticlesSecondsReadMonth(year, idx)))

			if !hideArticles {
				fmt.Printf("\n%s %d - Finished books: %d, articles: %d, time spend reading books: %s (hours: %s) and articles: %s (hours: %s)\n",
					months[idx], year, len(stats.BooksFinishedMonth(year, idx)), len(stats.ArticlesFinishedMonth(year, idx)),
					pkg.HumanizeDuration(monthBookReadDuration), pkg.SecondsToHoursString(stats.BooksSecondsReadMonth(year, idx)),
					pkg.HumanizeDuration(monthArticleReadDuration), pkg.SecondsToHoursString(stats.ArticlesSecondsReadMonth(year, idx)))
			} else {
				fmt.Printf("\n%s %d - Finished books: %d, time spend reading books: %s (hours: %s)\n",
					months[idx], year, len(stats.BooksFinishedMonth(year, idx)),
					pkg.HumanizeDuration(monthBookReadDuration), pkg.SecondsToHoursString(stats.BooksSecondsReadMonth(year, idx)))
			}

			if showBooks {
				for _, finishedBook := range stats.BooksFinishedMonth(year, idx) {
					duration := time.Duration(finishedBook.ReadSeconds()) * time.Second
					fmt.Printf("\t finished book: %s - %s (Duration: %s over %d Sessions)", finishedBook.Title, finishedBook.Author, duration, finishedBook.NumSessions())

					if showBookEnds {
						fmt.Printf(" Started: %s Finished: %s\n", formatTime(finishedBook.FirstReadTime()), formatTime(finishedBook.FinishedTime))
					} else {
						fmt.Printf("\n")
					}

					if showSessions {
						for jdx := range finishedBook.Reads {
							duration = time.Duration(finishedBook.Reads[jdx].Duration) * time.Second
							fmt.Printf("\t\tAt %s for %s\n", finishedBook.Reads[jdx].Time, duration)
						}
					}

					if showBookmarks {
						for jdx := range finishedBook.Bookmarks {
							fmt.Printf("\t\t%s: %s\n", finishedBook.Bookmarks[jdx].Type, finishedBook.Bookmarks[jdx].Text)
						}
					}
				}
			}

			if showArticles && !hideArticles {
				for _, finishedArticle := range stats.ArticlesFinishedMonth(year, idx) {
					fmt.Printf("\t finished article: %s - %s (%s)\n", finishedArticle.Title, finishedArticle.Author, finishedArticle.URL)
				}
			}
		}
	}

	return 0
}

func formatTime(ts string) string {
	return strings.Replace(strings.Replace(ts, "T", " ", 1), ".000", "", 1)
}
