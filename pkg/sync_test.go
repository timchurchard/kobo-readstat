package pkg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSync(t *testing.T) {
	t.Run("min", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := NewMockKoboDatabase(ctrl)
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

		db := NewMockKoboDatabase(ctrl)
		db.EXPECT().Contents().Return(contents, nil)
		db.EXPECT().Events().Return(aaaEvents, nil)
		db.EXPECT().Device().Return("xxx", "yyy")

		storage := NewMockStorage(ctrl)
		storage.EXPECT().AddContent("aaa", "Title AAA", "Author AAA", "aaa.com", 1234, true, false)
		storage.EXPECT().AddEvent("aaa", "xxx", "Read", time.Time{}, 1111)
		storage.EXPECT().AddDevice("xxx", "yyy")

		err := Sync(db, storage)
		assert.NoError(t, err)
	})
}
