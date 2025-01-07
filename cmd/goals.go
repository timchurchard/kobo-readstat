package cmd

import (
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"

	"github.com/snabb/isoweek"
	"github.com/timchurchard/kobo-readstat/pkg"
)

func Goals(out io.Writer) int {
	const ()

	var (
		storageFn    string
		year         int
		showSessions bool
	)

	flag.StringVar(&storageFn, "storage", defaultStorage, usageStoragePath)
	flag.StringVar(&storageFn, "s", defaultStorage, usageStoragePath)

	flag.IntVar(&year, "year", defaultYear, usageYear)
	flag.IntVar(&year, "y", defaultYear, usageYear)

	flag.BoolVar(&showSessions, "showsessions", false, usageShowSessions)

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

		if showSessions {
			fmt.Println("|    |Monday|Tuesda|Wednes|Thursd|Friday|Saturd|Sunday|        |")

			for idx := 0; idx < 24; idx++ {
				fmt.Println(readsInWeekToTableLine(stats, weekStartTime, idx))
			}

			fmt.Println(readsInWeekDayLine(stats, year, week))
		}

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

// readsInWeekToTableLine takes the week stats and looks for reading sessions in hour
// returns a string to match the table with 6 characters per (hour) e.g |0800|------|      | (indicating read on Mon not Tues)
func readsInWeekToTableLine(stats pkg.Stats, weekStart time.Time, hour int) string {
	var secs int

	line := fmt.Sprintf("|%02d00|      |      |      |      |      |      |      |", hour)
	year, week := weekStart.ISOWeek()
	totalSeconds := 0

	for idx := range stats.Years[year].Weeks[week].Books {
		for jdx := range stats.Years[year].Weeks[week].Books[idx].Reads {
			line, secs = bookReadToLineUpdate(line, stats.Years[year].Weeks[week].Books[idx].Reads[jdx], year, week, hour)
			totalSeconds += secs
		}
	}

	for idx := range stats.Years[year].Weeks[week].Articles {
		for jdx := range stats.Years[year].Weeks[week].Articles[idx].Reads {
			line, secs = bookReadToLineUpdate(line, stats.Years[year].Weeks[week].Articles[idx].Reads[jdx], year, week, hour)
			totalSeconds += secs
		}
	}

	totalDuration := time.Duration(totalSeconds) * time.Second

	return line + fmt.Sprintf("%8s|", totalDuration.String())
}

func bookReadToLineUpdate(line string, read pkg.StatsRead, year, week, hour int) (string, int) {
	const readSymbol = "-"

	// Offsets in the line for the day of the week
	offsets := map[time.Weekday]int{
		1: 7 - 1, // Monday is day 1 and starts on character 7
		2: 7*2 - 1,
		3: 7*3 - 1,
		4: 7*4 - 1,
		5: 7*5 - 1,
		6: 7*6 - 1,
		0: 7*7 - 1, // Sunday is day 0 (from time.Time WeekDay()) but rightmost on the line
	}

	readTime, _ := time.Parse(pkg.StorageTimeFmt, read.Time)
	readYear, readWeek := readTime.ISOWeek()
	secs := 0

	if year == readYear && week == readWeek && hour == readTime.Hour() {
		startMinute := min(max(readTime.Minute()/10, 0), 6)
		startOffset := offsets[readTime.Weekday()] + startMinute

		count := int(math.Ceil(float64(read.Duration) / 600))

		for startMinute+count > 6 {
			count -= 1
		}

		line = line[0:startOffset] + strings.Repeat(readSymbol, count) + line[startOffset+count:]

		secs = min(read.Duration, 3600)
	}

	return line, secs
}

func readsInWeekDayLine(stats pkg.Stats, year int, week int) string {
	dailySecs := []int{
		readsInWeekDay(stats, year, week, 1), readsInWeekDay(stats, year, week, 2),
		readsInWeekDay(stats, year, week, 3), readsInWeekDay(stats, year, week, 4),
		readsInWeekDay(stats, year, week, 5), readsInWeekDay(stats, year, week, 6),
		readsInWeekDay(stats, year, week, 0),
	}

	dailyAvgDur := time.Duration(int(calculateAverage([]float32{
		float32(dailySecs[0]), float32(dailySecs[1]), float32(dailySecs[2]),
		float32(dailySecs[3]), float32(dailySecs[4]), float32(dailySecs[5]),
		float32(dailySecs[6]),
	}))) * time.Second

	return fmt.Sprintf("|    |%s|%s|%s|%s|%s|%s|%s|        | Daily average %s\n\n",
		fmtSecondsToSixChar(dailySecs[0]), fmtSecondsToSixChar(dailySecs[1]), fmtSecondsToSixChar(dailySecs[2]),
		fmtSecondsToSixChar(dailySecs[3]), fmtSecondsToSixChar(dailySecs[4]), fmtSecondsToSixChar(dailySecs[5]),
		fmtSecondsToSixChar(dailySecs[6]),
		dailyAvgDur.String(),
	)
}

func fmtSecondsToSixChar(secs int) string {
	dur := time.Duration(secs) * time.Second
	if len(dur.String()) > 6 {
		return fmt.Sprintf("%6s", dur.String()[0:5])
	}

	return fmt.Sprintf("%6s", dur.String())
}

func readsInWeekDay(read pkg.Stats, year, week, day int) int {
	secs := 0

	for idx := range read.Years[year].Weeks[week].Books {
		for jdx := range read.Years[year].Weeks[week].Books[idx].Reads {
			readTime, _ := time.Parse(pkg.StorageTimeFmt, read.Years[year].Weeks[week].Books[idx].Reads[jdx].Time)
			readYear, readWeek := readTime.ISOWeek()
			readDay := int(readTime.Weekday())

			if year == readYear && week == readWeek && day == readDay {
				secs += read.Years[year].Weeks[week].Books[idx].Reads[jdx].Duration
			}
		}
	}

	for idx := range read.Years[year].Weeks[week].Articles {
		for jdx := range read.Years[year].Weeks[week].Articles[idx].Reads {
			readTime, _ := time.Parse(pkg.StorageTimeFmt, read.Years[year].Weeks[week].Articles[idx].Reads[jdx].Time)
			readYear, readWeek := readTime.ISOWeek()
			readDay := int(readTime.Weekday())

			if year == readYear && week == readWeek && day == readDay {
				secs += read.Years[year].Weeks[week].Articles[idx].Reads[jdx].Duration
			}
		}
	}

	return secs
}
