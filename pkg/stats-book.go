package pkg

import "time"

type StatsBook struct {
	BookID string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	URL    string `json:"url"`

	Words int `json:"words"`

	IsBook       bool   `json:"is_book"`
	IsFinished   bool   `json:"is_finished"`
	FinishedTime string `json:"finished_time"`

	Reads     []StatsRead     `json:"reads"`
	Bookmarks []StatsBookmark `json:"bookmarks"`
}

type StatsRead struct {
	Time     string `json:"time"`
	Duration int    `json:"duration"`
}

type StatsBookmark struct {
	ID          string `json:"id,omitempty"`
	Index       int    `json:"idx,omitempty"`
	StartOffset int    `json:"start_offset,omitempty"`
	EndOffset   int    `json:"end_offset,omitempty"`
	Text        string `json:"text,omitempty"`
	Annotation  string `json:"annotation,omitempty"`
	Created     string `json:"created,omitempty"`
	Modified    string `json:"modified,omitempty"`
	Type        string `json:"type,omitempty"`
}

func (b StatsBook) FirstReadTime() string {
	result := b.FinishedTime

	for idx := range b.Reads {
		if b.Reads[idx].Time < result {
			result = b.Reads[idx].Time
		}
	}

	return result
}

func (b StatsBook) ReadSeconds() int {
	result := 0

	for idx := range b.Reads {
		result += b.Reads[idx].Duration
	}

	return result
}

func (b StatsBook) ReadSecondsInYear(year int) int {
	result := 0

	for idx := range b.Reads {
		readTime, _ := time.Parse(StorageTimeFmt, b.Reads[idx].Time)

		if year == readTime.Year() {
			result += b.Reads[idx].Duration
		}
	}

	return result
}

func (b StatsBook) ReadSecondsInWeek(year int, week int) int {
	result := 0

	for idx := range b.Reads {
		readTime, _ := time.Parse(StorageTimeFmt, b.Reads[idx].Time)

		_, readWeek := readTime.ISOWeek()
		if year == readTime.Year() && week == readWeek {
			result += b.Reads[idx].Duration
		}
	}

	return result
}

func (b StatsBook) NumSessions() int {
	return len(b.Reads)
}

func (b StatsBook) NumSessionsInYear(year int) int {
	result := 0

	for idx := range b.Reads {
		readTime, _ := time.Parse(StorageTimeFmt, b.Reads[idx].Time)

		if year == readTime.Year() {
			result++
		}
	}

	return result
}
