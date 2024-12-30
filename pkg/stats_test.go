package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewStatsForYear(t *testing.T) {
	const (
		testDeviceAID = "test-device-a"

		testBookAID    = "books/test-book-a.epub"
		testBookBID    = "books/test-book-b.epub"
		testArticleAID = "articles/test-article-a.epub"
	)

	t.Run("", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		testContents := []StorageContent{
			{
				ID:         testBookAID,
				Title:      "Test Book A",
				Author:     "AAA",
				URL:        "",
				Words:      123,
				IsBook:     true,
				IsFinished: true,
			},
			{
				ID:         testBookBID,
				Title:      "Test Book B",
				Author:     "BBB",
				URL:        "",
				Words:      9999,
				IsBook:     false,
				IsFinished: false,
			},
			{
				ID:         testArticleAID,
				Title:      "Test Article A",
				Author:     "XXX",
				URL:        "test.com/article/a",
				Words:      1111,
				IsBook:     false,
				IsFinished: true,
			},
		}

		testStorage := NewMockStorage(ctrl)
		testStorage.EXPECT().Contents().Return(testContents)
		testStorage.EXPECT().Events(testBookAID).Return([]StorageEvents{
			{EventName: "Read", Time: "2020-02-01T01:02:03.000", Duration: 100, Device: testDeviceAID},
			{EventName: "Read", Time: "2020-02-01T02:02:03.000", Duration: 200, Device: testDeviceAID},
			{EventName: "Read", Time: "2020-02-01T03:02:03.000", Duration: 300, Device: testDeviceAID},
			{EventName: "Finish", Time: "2020-02-01T03:02:03.000", Duration: 0, Device: testDeviceAID},
		})
		testStorage.EXPECT().Events(testBookBID).Return([]StorageEvents{
			{EventName: "Read", Time: "2020-02-02T01:02:03.000", Duration: 10, Device: testDeviceAID},
			{EventName: "Read", Time: "2020-02-02T02:02:03.000", Duration: 20, Device: testDeviceAID},
			{EventName: "Read", Time: "2020-02-02T03:02:03.000", Duration: 30, Device: testDeviceAID},
		})
		testStorage.EXPECT().Events(testArticleAID).Return([]StorageEvents{
			{EventName: "Read", Time: "2020-02-03T02:02:03.000", Duration: 1000, Device: testDeviceAID},
			{EventName: "Finish", Time: "2020-02-03T03:02:03.000", Duration: 0, Device: testDeviceAID},
		})
		testStorage.EXPECT().Bookmarks(testBookAID).Return([]StorageBookmark{
			{ID: "aaa", VolumeID: "bbb", ContentID: testBookAID, BookPath: "ccc", Index: 123, StartOffset: 0, EndOffset: 0, Text: "text"},
		})
		testStorage.EXPECT().Bookmarks(testBookBID).Return([]StorageBookmark{})
		testStorage.EXPECT().Bookmarks(testArticleAID).Return([]StorageBookmark{})

		yearStats := NewStats(testStorage)

		finYear := yearStats.BooksFinishedYear(2020)
		assert.Len(t, finYear, 1)
		assert.Equal(t, testBookAID, finYear[0].BookID)
		assert.Len(t, finYear[0].Reads, 3)
	})
}
