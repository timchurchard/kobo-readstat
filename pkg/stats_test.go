package pkg

import (
	"fmt"
	"github.com/timchurchard/readstat/internal"
	"testing"
)

func TestNewStatsForYear(t *testing.T) {
	const (
		testDeviceAID = "test-device-a"

		testBookAID    = "books/test-book-a.epub"
		testBookBID    = "books/test-book-b.epub"
		testArticleAID = "articles/test-article-a.epub"
	)

	testStorage := internal.Storage{
		Devices: map[string]internal.StorageDevice{
			testDeviceAID: internal.StorageDevice{
				Device: "Test Device 123",
				Model:  "ABC123",
			},
		},
		Contents: map[string]internal.StorageContent{
			testBookAID: {
				ID:         testBookAID,
				Title:      "Test Book A",
				Author:     "AAA",
				URL:        "",
				Words:      123,
				IsBook:     true,
				IsFinished: true,
			},
			testBookBID: {
				ID:         testBookBID,
				Title:      "Test Book B",
				Author:     "BBB",
				URL:        "",
				Words:      9999,
				IsBook:     false,
				IsFinished: false,
			},
			testArticleAID: {
				ID:         testArticleAID,
				Title:      "Test Article A",
				Author:     "XXX",
				URL:        "test.com/article/a",
				Words:      1111,
				IsBook:     false,
				IsFinished: true,
			},
		},
		Events: map[string][]internal.StorageEvents{
			testBookAID: {
				{EventName: "Read", Time: "2020-02-01T01:02:03.000", Duration: 100, Device: testDeviceAID},
				{EventName: "Read", Time: "2020-02-01T02:02:03.000", Duration: 200, Device: testDeviceAID},
				{EventName: "Read", Time: "2020-02-01T03:02:03.000", Duration: 300, Device: testDeviceAID},
				{EventName: "Finish", Time: "2020-02-01T03:02:03.000", Duration: 0, Device: testDeviceAID},
			},
			testBookBID: {
				{EventName: "Read", Time: "2020-02-02T01:02:03.000", Duration: 10, Device: testDeviceAID},
				{EventName: "Read", Time: "2020-02-02T02:02:03.000", Duration: 20, Device: testDeviceAID},
				{EventName: "Read", Time: "2020-02-02T03:02:03.000", Duration: 30, Device: testDeviceAID},
			},
			testArticleAID: {
				{EventName: "Read", Time: "2020-02-03T02:02:03.000", Duration: 1000, Device: testDeviceAID},
				{EventName: "Finish", Time: "2020-02-03T03:02:03.000", Duration: 0, Device: testDeviceAID},
			},
		},
	}

	yearStats := NewStatsForYear(testStorage, 2020)

	finYear := yearStats.BooksFinishedYear()
	fmt.Println(finYear)
}
