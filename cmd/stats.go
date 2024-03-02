package cmd

import (
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/timchurchard/readstat/internal"
	"github.com/timchurchard/readstat/pkg"
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

		usageShowBooks    = "Show book title and details"
		usageShowArticles = "Show pocket article title and details"
		usageShowSessions = "Show reading sessions"
	)
	var (
		storageFn    string
		mode         string
		year         int
		outFn        string
		showBooks    bool
		showArticles bool
		showSessions bool
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
	flag.BoolVar(&showSessions, "showsessions", false, usageShowSessions)

	flag.Usage = func() {
		fmt.Fprintf(out, "Usage of %s %s:\n", os.Args[0], os.Args[1])

		flag.PrintDefaults()
	}

	flag.Parse()

	if _, err := os.Stat(storageFn); err != nil {
		panic(fmt.Sprintf("storage not found: %v", err))
	}

	storage, err := internal.OpenStorageOrCreate(storageFn)
	if err != nil {
		panic(err)
	}

	stats := pkg.NewStatsForYear(storage, year)

	booksReadSeconds := stats.BooksSecondsReadYear()
	booksReadDuration, _ := time.ParseDuration(fmt.Sprintf("%ds", booksReadSeconds))
	articlesReadSeconds := stats.ArticlesSecondsReadYear()
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
		fmt.Printf("Finished books\t\t\t: %d\n", len(stats.BooksFinishedYear()))
		fmt.Printf("Finished articles\t\t: %d\n", len(stats.ArticlesFinishedYear()))
		fmt.Printf("Total finished words\t\t: %s\n", humanizeInt(stats.WordsFinishedYear()))
		fmt.Printf("Time reading books\t\t: %s\n", humanizeDuration(booksReadDuration))
		fmt.Printf("Time reading articles\t\t: %s\n", humanizeDuration(articlesReadDuration))
		fmt.Printf("Total time reading\t\t: %s\n", humanizeDuration(totalReadDuration))

		fmt.Println("\n----------")

		months := []string{"", "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
		for idx := 1; idx <= 12; idx++ {
			monthBookReadDuration, _ := time.ParseDuration(fmt.Sprintf("%ds", stats.BooksSecondsReadMonth(idx)))
			monthArticleReadDuration, _ := time.ParseDuration(fmt.Sprintf("%ds", stats.ArticlesSecondsReadMonth(idx)))

			fmt.Printf("\n%s %d - Finished books: %d, articles: %d, time spend reading books: %s and articles: %s\n", months[idx], year, len(stats.BooksFinishedMonth(idx)), len(stats.ArticlesFinishedMonth(idx)), humanizeDuration(monthBookReadDuration), humanizeDuration(monthArticleReadDuration))

			if showBooks {
				for _, finishedBook := range stats.BooksFinishedMonth(idx) {
					duration := time.Duration(finishedBook.ReadSeconds()) * time.Second
					fmt.Printf("\t finished book: %s - %s (Duration: %s over %d Sessions)\n", finishedBook.Title, finishedBook.Author, duration, finishedBook.NumSessions())

					if showSessions {
						for jdx := range finishedBook.Reads {
							duration = time.Duration(finishedBook.Reads[jdx].Duration) * time.Second
							fmt.Printf("\t\tAt %s for %s\n", finishedBook.Reads[jdx].Time, duration)
						}
					}
				}
			}

			if showArticles {
				for _, finishedArticle := range stats.ArticlesFinishedMonth(idx) {
					fmt.Printf("\t finished article: %s - %s (%s)\n", finishedArticle.Title, finishedArticle.Author, finishedArticle.URL)
				}
			}
		}
	}

	return 0
}

// humanizeDuration humanizes time.Duration output to a meaningful value,
// golang's default “time.Duration“ output is badly formatted and unreadable.
// From: https://gist.github.com/harshavardhana/327e0577c4fed9211f65
func humanizeDuration(duration time.Duration) string {
	if duration.Seconds() < 60.0 {
		return fmt.Sprintf("%d seconds", int64(duration.Seconds()))
	}
	if duration.Minutes() < 60.0 {
		remainingSeconds := math.Mod(duration.Seconds(), 60)
		return fmt.Sprintf("%d minutes %d seconds", int64(duration.Minutes()), int64(remainingSeconds))
	}
	if duration.Hours() < 24.0 {
		remainingMinutes := math.Mod(duration.Minutes(), 60)
		remainingSeconds := math.Mod(duration.Seconds(), 60)
		return fmt.Sprintf("%d hours %d minutes %d seconds",
			int64(duration.Hours()), int64(remainingMinutes), int64(remainingSeconds))
	}
	remainingHours := math.Mod(duration.Hours(), 24)
	remainingMinutes := math.Mod(duration.Minutes(), 60)
	remainingSeconds := math.Mod(duration.Seconds(), 60)
	return fmt.Sprintf("%d days %d hours %d minutes %d seconds",
		int64(duration.Hours()/24), int64(remainingHours),
		int64(remainingMinutes), int64(remainingSeconds))
}

// humanizeInt
// Based on https://github.com/dustin/go-humanize/blob/v1.0.1/comma.go#L15
func humanizeInt(num int) string {
	parts := []string{"", "", "", "", "", "", ""}
	j := len(parts) - 1

	for num > 999 {
		parts[j] = strconv.FormatInt(int64(num%1000), 10)
		num = num / 1000
		j--
	}

	parts[j] = strconv.Itoa(int(num))
	return strings.Join(parts[j:], ",")
}
