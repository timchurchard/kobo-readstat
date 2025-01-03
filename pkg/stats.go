//go:generate mockgen -package pkg -destination stats_mock.go -source stats.go
package pkg

import (
	"time"
)

type Statter interface {
	BooksFinishedYear(year int) []*StatsBook
	BooksFinishedMonth(year, month int) []*StatsBook

	ArticlesFinishedYear(year int) []*StatsBook
	ArticlesFinishedMonth(year, month int) []*StatsBook

	BooksSecondsReadYear(year int) int
	BooksSecondsReadMonth(year, month int) int

	ArticlesSecondsReadYear(year int) int
	ArticlesSecondsReadMonth(year, month int) int
}

type Stats struct {
	Years map[int]YearStats `json:"years"`

	Content map[string]StatsBook `json:"content"`
}

type YearStats struct {
	Months map[int]MonthStats `json:"months"`
	Weeks  map[int]MonthStats `json:"weeks"`
}

type MonthStats struct {
	FinishedBooks    map[string]*StatsBook `json:"finished"`
	FinishedArticles map[string]*StatsBook `json:"finished_articles"`

	Books    map[string]*StatsBook `json:"progressed"`
	Articles map[string]*StatsBook `json:"progressed_articles"`
}

// NewStats take a Storage and calculate the full stats
func NewStats(storage Storage) Stats {
	result := Stats{
		Years:   make(map[int]YearStats),
		Content: make(map[string]StatsBook),
	}

	// Collect all info per content
	for _, content := range storage.Contents() {
		book := StatsBook{
			BookID:     content.ID,
			Title:      content.Title,
			Author:     content.Author,
			URL:        content.URL,
			Words:      content.Words,
			IsBook:     content.IsBook,
			IsFinished: content.IsFinished,
			Reads:      make([]StatsRead, 0),
			Bookmarks:  make([]StatsBookmark, 0),
		}

		for _, event := range storage.Events(content.ID) {
			eventTime, err := time.Parse(StorageTimeFmt, event.Time)
			if err != nil {
				panic(err)
			}

			switch event.EventName {
			case FinishEvent.String():
				book.IsFinished = true
				book.FinishedTime = event.Time

			case ReadEvent.String():
				book.Reads = append(book.Reads, StatsRead{
					Time:     event.Time,
					Duration: event.Duration,
				})

				if book.IsFinished && book.FinishedTime == "" {
					book.FinishedTime = event.Time // TODO!
				}
			}

			if _, exists := result.Years[eventTime.Year()]; !exists {
				result.Years[eventTime.Year()] = YearStats{
					Months: make(map[int]MonthStats),
					Weeks:  make(map[int]MonthStats),
				}

				for idx := 1; idx <= 12; idx++ {
					result.Years[eventTime.Year()].Months[idx] = MonthStats{
						FinishedBooks:    make(map[string]*StatsBook),
						FinishedArticles: make(map[string]*StatsBook),
						Books:            make(map[string]*StatsBook),
						Articles:         make(map[string]*StatsBook),
					}
				}

				for idx := 1; idx <= 53; idx++ {
					result.Years[eventTime.Year()].Weeks[idx] = MonthStats{
						FinishedBooks:    make(map[string]*StatsBook),
						FinishedArticles: make(map[string]*StatsBook),
						Books:            make(map[string]*StatsBook),
						Articles:         make(map[string]*StatsBook),
					}
				}
			}
		}

		bookmarks := storage.Bookmarks(content.ID)
		for idx := range bookmarks {
			book.Bookmarks = append(book.Bookmarks, StatsBookmark{
				ID:          bookmarks[idx].ID,
				Index:       bookmarks[idx].Index,
				StartOffset: bookmarks[idx].StartOffset,
				EndOffset:   bookmarks[idx].EndOffset,
				Text:        bookmarks[idx].Text,
				Annotation:  bookmarks[idx].Annotation,
				Created:     bookmarks[idx].Created,
				Modified:    bookmarks[idx].Modified,
				Type:        bookmarks[idx].Type,
			})
		}

		result.Content[content.ID] = book
	}

	// Content to Months
	for cid := range result.Content {
		book := result.Content[cid]

		if result.Content[cid].IsFinished && result.Content[cid].FinishedTime != "" {
			finishedTime, err := time.Parse(StorageTimeFmt, result.Content[cid].FinishedTime)
			if err != nil {
				panic(err)
			}

			finishedYear := finishedTime.Year()
			finishedMonth := int(finishedTime.Month())
			finishedYearWeek, finishedWeek := finishedTime.ISOWeek() // Note: readYearWeek may differ from readYear e.g. the timestamp is 2024-12-30 but that is in week 1 of 2025

			if result.Content[cid].IsBook {
				if _, exists := result.Years[finishedYear].Months[finishedMonth].FinishedBooks[cid]; !exists {
					result.Years[finishedYear].Months[finishedMonth].FinishedBooks[cid] = &book
				}

				if _, exists := result.Years[finishedYearWeek].Weeks[finishedWeek].FinishedBooks[cid]; !exists {
					result.Years[finishedYearWeek].Weeks[finishedWeek].FinishedBooks[cid] = &book
				}
			} else {
				if _, exists := result.Years[finishedYear].Months[finishedMonth].FinishedArticles[cid]; !exists {
					result.Years[finishedYear].Months[finishedMonth].FinishedArticles[cid] = &book
				}

				if _, exists := result.Years[finishedYearWeek].Weeks[finishedWeek].FinishedArticles[cid]; !exists {
					result.Years[finishedYearWeek].Weeks[finishedWeek].FinishedArticles[cid] = &book
				}
			}
		}

		for idx := range result.Content[cid].Reads {
			readTime, err := time.Parse(StorageTimeFmt, result.Content[cid].Reads[idx].Time)
			if err != nil {
				panic(err)
			}

			readYear := readTime.Year()
			readMonth := int(readTime.Month())
			readYearWeek, readWeek := readTime.ISOWeek() // Note: readYearWeek may differ from readYear e.g. the timestamp is 2024-12-30 but that is in week 1 of 2025

			if result.Content[cid].IsBook {
				if _, exists := result.Years[readYear].Months[readMonth].Books[cid]; !exists {
					result.Years[readYear].Months[readMonth].Books[cid] = &book
				}

				if _, exists := result.Years[readYearWeek].Weeks[readWeek].Books[cid]; !exists {
					result.Years[readYearWeek].Weeks[readWeek].Books[cid] = &book
				}
			} else {
				if _, exists := result.Years[readYear].Months[readMonth].Articles[cid]; !exists {
					result.Years[readYear].Months[readMonth].Articles[cid] = &book
				}

				if _, exists := result.Years[readYearWeek].Weeks[readWeek].Articles[cid]; !exists {
					result.Years[readYearWeek].Weeks[readWeek].Articles[cid] = &book
				}
			}
		}
	}

	return result
}

func (s Stats) BooksFinishedYear(year int) []*StatsBook {
	result := make([]*StatsBook, 0)

	for month := range s.Years[year].Months {
		result = append(result, s.BooksFinishedMonth(year, month)...)
	}

	return result
}

func (s Stats) BooksFinishedMonth(year, month int) []*StatsBook {
	result := make([]*StatsBook, 0)

	for idx := range s.Years[year].Months[month].FinishedBooks {
		result = append(result, s.Years[year].Months[month].FinishedBooks[idx])
	}

	return result
}

func (s Stats) ArticlesFinishedYear(year int) []*StatsBook {
	result := make([]*StatsBook, 0)

	for month := range s.Years[year].Months {
		result = append(result, s.ArticlesFinishedMonth(year, month)...)
	}

	return result
}

func (s Stats) ArticlesFinishedMonth(year, month int) []*StatsBook {
	result := make([]*StatsBook, 0)

	for idx := range s.Years[year].Months[month].FinishedArticles {
		result = append(result, s.Years[year].Months[month].FinishedArticles[idx])
	}

	return result
}

func (s Stats) BooksSecondsReadYear(year int) int {
	result := 0

	for month := 1; month < 13; month++ {
		result += s.BooksSecondsReadMonth(year, month)
	}

	return result
}

func (s Stats) BooksSecondsReadMonth(year, month int) int {
	result := 0

	for idx := range s.Years[year].Months[month].Books {
		for jdx := range s.Years[year].Months[month].Books[idx].Reads {
			result += s.Years[year].Months[month].Books[idx].Reads[jdx].Duration
		}
	}

	return result
}

func (s Stats) BooksSecondsReadWeek(year, week int) int {
	result := 0

	for idx := range s.Years[year].Weeks[week].Books {
		result += s.Years[year].Weeks[week].Books[idx].ReadSecondsInWeek(year, week)
	}

	return result
}

func (s Stats) ArticlesSecondsReadYear(year int) int {
	result := 0

	for month := 1; month < 13; month++ {
		result += s.ArticlesSecondsReadMonth(year, month)
	}

	return result
}

func (s Stats) ArticlesSecondsReadMonth(year, month int) int {
	result := 0

	for idx := range s.Years[year].Months[month].Articles {
		for jdx := range s.Years[year].Months[month].Articles[idx].Reads {
			result += s.Years[year].Months[month].Articles[idx].Reads[jdx].Duration
		}
	}

	return result
}

func (s Stats) ArticleSecondsReadWeek(year, week int) int {
	result := 0

	for idx := range s.Years[year].Weeks[week].Articles {
		result += s.Years[year].Weeks[week].Articles[idx].ReadSecondsInWeek(year, week)
	}

	return result
}
