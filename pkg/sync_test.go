package pkg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSync(t *testing.T) {
	/*t.Run("real", func(t *testing.T) {
		// Test a corrupt database with mismatched start/end time events
		db, err := NewKoboDatabase("../testfiles/20240513/clara2e/KoboReader.sqlite")
		assert.NoError(t, err)

		db.Events()
	})*/

	t.Run("min", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := NewMockKoboDatabase(ctrl)
		db.EXPECT().Shelves() // TODO gomock inorder
		db.EXPECT().ShelfContents()
		db.EXPECT().Bookmarks()
		db.EXPECT().Contents()
		db.EXPECT().Events()
		db.EXPECT().Device().Return("aaa", "bbb")

		storage := NewMockStorage(ctrl)
		storage.EXPECT().AddDevice("aaa", "bbb")

		err := Sync(db, storage)
		assert.NoError(t, err)
	})

	t.Run("simple", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		contents := []KoboBook{
			{
				ID:              "aaa",
				Title:           "Title AAA",
				Author:          "Author AAA",
				URL:             "aaa.com",
				Finished:        false,
				ProgressPercent: 69,
				Parts: map[string]KoboBookPart{
					"aaa": {
						WordCount: 1234,
					},
				},
				IsBook: true,
			},
		}

		aaaEvents := []KoboEvent{
			{
				BookID:    "aaa",
				EventType: "Read",
				Time:      time.Time{},
				ReadingSessions: []KoboEventReadingSession{
					{
						UnixStart: 123,
						UnixEnd:   1234,
						Start:     time.Time{},
						End:       time.Time{},
					},
				},
			},
		}

		aaaShelves := []KoboShelf{
			{
				ID:           "AAA",
				Name:         "BBB",
				InternalName: "CCC",
				Type:         "DDD",
				IsDeleted:    false,
			},
		}

		aaaShelfContent := []KoboShelfContent{
			{
				ShelfName: "AAA",
				ContentID: "aaa",
				IsDeleted: false,
			},
		}

		aaaBookmarks := []KoboBookmark{
			{
				ID:          "AA",
				VolumeID:    "BB",
				ContentID:   "CC",
				BookPath:    "DD",
				Index:       1,
				StartOffset: 2,
				EndOffset:   3,
				Text:        "EE",
				Annotation:  "FF",
				Created:     time.Date(2001, 1, 2, 3, 4, 5, 0, time.UTC),
				Modified:    time.Date(2002, 1, 2, 3, 4, 5, 0, time.UTC),
				Type:        "GG",
			},
		}

		db := NewMockKoboDatabase(ctrl)
		db.EXPECT().Contents().Return(contents, nil)
		db.EXPECT().Events().Return(aaaEvents, nil)
		db.EXPECT().Device().Return("xxx", "yyy")
		db.EXPECT().Shelves().Return(aaaShelves, nil)
		db.EXPECT().ShelfContents().Return(aaaShelfContent, nil)
		db.EXPECT().Bookmarks().Return(aaaBookmarks, nil)

		storage := NewMockStorage(ctrl)
		storage.EXPECT().AddContent("aaa", "Title AAA", "Author AAA", "aaa.com", 1234, true, false, 69)
		storage.EXPECT().AddEvent("aaa", "xxx", "Read", time.Time{}, 1111)
		storage.EXPECT().AddDevice("xxx", "yyy")
		storage.EXPECT().AddShelf("AAA", "BBB", "CCC", "DDD", false)
		storage.EXPECT().AddShelfContent("AAA", "aaa", false)
		storage.EXPECT().AddBookmark("AA", "BB", "CC", "GG", "DD", 1, 2, 3, "EE", "FF",
			time.Date(2001, 1, 2, 3, 4, 5, 0, time.UTC), time.Date(2002, 1, 2, 3, 4, 5, 0, time.UTC))

		err := Sync(db, storage)
		assert.NoError(t, err)
	})
}
