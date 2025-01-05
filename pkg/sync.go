package pkg

import (
	"time"
)

func Sync(db KoboDatabase, storage Storage) error {
	const readingSpeedWPM = 180 / 60 // reading speed words-per-minute (in seconds) to guess reading time by word count (should be ~200)

	contents, err := db.Contents()
	if err != nil {
		return err
	}

	events, err := db.Events()
	if err != nil {
		return err
	}

	shelf, err := db.Shelves()
	if err != nil {
		return err
	}

	shelfContent, err := db.ShelfContents()
	if err != nil {
		return err
	}

	bookmarks, err := db.Bookmarks()
	if err != nil {
		return err
	}

	device, model := db.Device()
	storage.AddDevice(device, model)

	for cIdx := range contents {
		storage.AddContent(contents[cIdx].ID, contents[cIdx].Title, contents[cIdx].Author,
			contents[cIdx].URL, contents[cIdx].TotalWords(), contents[cIdx].IsBook,
			contents[cIdx].Finished, contents[cIdx].ProgressPercent)
	}

	for sIdx := range shelf {
		storage.AddShelf(shelf[sIdx].ID, shelf[sIdx].Name, shelf[sIdx].InternalName, shelf[sIdx].Type, shelf[sIdx].IsDeleted)
	}

	for lIdx := range shelfContent {
		storage.AddShelfContent(shelfContent[lIdx].ShelfName, shelfContent[lIdx].ContentID, shelfContent[lIdx].IsDeleted)
	}

	for bIdx := range bookmarks {
		storage.AddBookmark(bookmarks[bIdx].ID, bookmarks[bIdx].VolumeID, bookmarks[bIdx].ContentID, bookmarks[bIdx].Type,
			bookmarks[bIdx].BookPath, bookmarks[bIdx].Index, bookmarks[bIdx].StartOffset, bookmarks[bIdx].EndOffset,
			bookmarks[bIdx].Text, bookmarks[bIdx].Annotation, bookmarks[bIdx].Created, bookmarks[bIdx].Modified)
	}

	for eIdx := range events {
		if events[eIdx].EventType == GuessReadEvent {
			// Guess Read Events (for pocket articles that are finished but do not have reading seconds)
			startUnix := int(events[eIdx].Time.Unix())

			secondsRead := 0
			for cIdx := range contents {
				if contents[cIdx].ID == events[eIdx].BookID {
					secondsRead = contents[cIdx].TotalWords() / readingSpeedWPM

					break
				}
			}

			events[eIdx].EventType = ReadEvent
			events[eIdx].ReadingSessions = []KoboEventReadingSession{
				{
					UnixStart: startUnix,
					UnixEnd:   startUnix + secondsRead,
					Start:     time.Unix(int64(startUnix), 0),
					End:       time.Unix(int64(startUnix+secondsRead), 0),
				},
			}
		}

		if events[eIdx].EventType == ReadEvent {
			for sIdx := range events[eIdx].ReadingSessions {
				durationSecs := events[eIdx].ReadingSessions[sIdx].UnixEnd - events[eIdx].ReadingSessions[sIdx].UnixStart
				storage.AddEvent(events[eIdx].BookID, device, events[eIdx].EventType.String(),
					events[eIdx].ReadingSessions[sIdx].Start, durationSecs)
			}
		} else {
			storage.AddEvent(events[eIdx].BookID, device, events[eIdx].EventType.String(), events[eIdx].Time, 0)
		}
	}

	return nil
}
