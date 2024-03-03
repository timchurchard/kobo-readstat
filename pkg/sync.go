package pkg

func Sync(db KoboDatabase, storage Storage) error {
	contents, err := db.Contents()
	if err != nil {
		return err
	}

	events, err := db.Events()
	if err != nil {
		return err
	}

	device, model := db.Device()
	storage.AddDevice(device, model)

	for cIdx := range contents {
		storage.AddContent(contents[cIdx].ID, contents[cIdx].Title, contents[cIdx].Author,
			contents[cIdx].URL, contents[cIdx].TotalWords(), contents[cIdx].IsBook, contents[cIdx].Finished)
	}

	for eIdx := range events {
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
