//go:generate mockgen -package pkg -destination storage_mock.go -source storage.go
package pkg

import (
	"encoding/json"
	"os"
	"time"
)

type Storage interface {
	Save() error

	AddContent(fn, title, author, url string, words int, book, finished bool, percent int)
	AddDevice(device, model string)
	AddEvent(fn, device, name string, t time.Time, duration int)

	AddShelf(ID, name, internalName, shelfType string, isDeleted bool)
	AddShelfContent(shelfName, fn string, isDeleted bool)

	Contents() []StorageContent
	Events(cID string) []StorageEvents
	Bookmarks(cID string) []StorageBookmark

	Shelfs() []StorageShelf
	ShelfContents(shelfName string) []StorageShelfContent

	AddBookmark(bID, vID, cID, typeStr, path string, index, startOffset, endOffset int, text, annotation string, created, modified time.Time)
}

type JSONStorage struct {
	DeviceMap  map[string]StorageDevice   `json:"devices"`
	ContentMap map[string]StorageContent  `json:"contents"`
	EventMap   map[string][]StorageEvents `json:"events"`

	// Extra data that may be interesting. TODO need to think about relating content across devices e.g. books with different CIDs, fine for stats but not for shelve content
	Shelf        map[string]StorageShelf          `json:"shelf"`
	ShelfContent map[string][]StorageShelfContent `json:"shelf_content"`
	Bookmark     map[string][]StorageBookmark     `json:"bookmark"`

	fn string
}

type StorageDevice struct {
	Device string `json:"device"`
	Model  string `json:"model"`
}

type StorageContent struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	URL    string `json:"url"`

	Words int `json:"words"`

	IsBook     bool `json:"book"`
	IsFinished bool `json:"article_is_finished"`
}

type StorageEvents struct {
	EventName string `json:"event"`
	Time      string `json:"time"`
	Duration  int    `json:"duration"`
	Device    string `json:"device"`
}

type StorageShelf struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	InternalName string `json:"internal_name"`
	Type         string `json:"type"`
	IsDeleted    bool   `json:"is_deleted"`
}

type StorageShelfContent struct {
	ShelfID   string `json:"shelf_id"`
	ContentID string `json:"content_id"`
	IsDeleted bool   `json:"is_deleted"`
}

type StorageBookmark struct {
	ID          string `json:"id,omitempty"`
	VolumeID    string `json:"volume_id,omitempty"`
	ContentID   string `json:"content_id,omitempty"`
	BookPath    string `json:"book_path,omitempty"`
	Index       int    `json:"idx,omitempty"`
	StartOffset int    `json:"start_offset,omitempty"`
	EndOffset   int    `json:"end_offset,omitempty"`
	Text        string `json:"text,omitempty"`
	Annotation  string `json:"annotation,omitempty"`
	Created     string `json:"created,omitempty"`
	Modified    string `json:"modified,omitempty"`
	Type        string `json:"type,omitempty"`
}

const (
	StorageTimeFmt = "2006-01-02T15:04:05.000"
)

func OpenStorageOrCreate(fn string) (Storage, error) {
	storage := JSONStorage{
		DeviceMap:    map[string]StorageDevice{},
		ContentMap:   map[string]StorageContent{},
		EventMap:     map[string][]StorageEvents{},
		Shelf:        map[string]StorageShelf{},
		ShelfContent: map[string][]StorageShelfContent{},
		Bookmark:     map[string][]StorageBookmark{},
		fn:           fn,
	}

	if _, err := os.Stat(fn); err == nil {
		storageBytes, err := os.ReadFile(fn)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(storageBytes, &storage)
		if err != nil {
			return nil, err
		}
	}

	return &storage, nil
}

func (s *JSONStorage) Save() error {
	storageBytes, err := json.Marshal(s)
	if err != nil {
		return err
	}

	return os.WriteFile(s.fn, storageBytes, 0o644)
}

func (s *JSONStorage) AddContent(fn, title, author, url string, words int, book, finished bool, percent int) {
	previouslyFinished := false
	if _, ok := s.ContentMap[fn]; ok {
		previouslyFinished = s.ContentMap[fn].IsFinished
	}

	if !book && percent == 100 {
		// Pocket articles work around where finished column is false but progress is 100%
		finished = true
	}

	s.ContentMap[fn] = StorageContent{
		ID:         fn,
		Title:      title,
		Author:     author,
		Words:      words,
		URL:        url,
		IsBook:     book,
		IsFinished: finished || previouslyFinished, // Content cannot go from 'finished' to unfinished (e.g. duplicate content across multiple devices)
	}
}

func (s *JSONStorage) AddDevice(device, model string) {
	s.DeviceMap[device] = StorageDevice{
		Device: device,
		Model:  model,
	}
}

func (s *JSONStorage) AddEvent(fn, device, name string, t time.Time, duration int) {
	timeStr := t.Format(StorageTimeFmt)

	found := false
	for eIdx := range s.EventMap[fn] {
		if s.EventMap[fn][eIdx].EventName == name && s.EventMap[fn][eIdx].Time == timeStr {
			found = true
		}
	}

	if !found {
		s.EventMap[fn] = append(s.EventMap[fn], StorageEvents{
			EventName: name,
			Time:      timeStr,
			Duration:  duration,
			Device:    device,
		})
	}
}

func (s *JSONStorage) Contents() []StorageContent {
	result := make([]StorageContent, len(s.ContentMap))
	idx := 0

	for _, content := range s.ContentMap {
		result[idx] = content
		idx += 1
	}

	return result
}

func (s *JSONStorage) Events(cID string) []StorageEvents {
	result := make([]StorageEvents, len(s.EventMap[cID]))
	idx := 0

	for _, event := range s.EventMap[cID] {
		result[idx] = event
		idx += 1
	}

	return result
}

func (s *JSONStorage) AddShelf(ID, name, internalName, shelfType string, isDeleted bool) {
	if s.Shelf == nil {
		s.Shelf = map[string]StorageShelf{} // TODO/FIXME ! panic without but is initialised in Open function ??
	}

	s.Shelf[ID] = StorageShelf{
		ID:           ID,
		Name:         name,
		InternalName: internalName,
		Type:         shelfType,
		IsDeleted:    isDeleted,
	}
}

func (s *JSONStorage) AddShelfContent(shelfName, fn string, isDeleted bool) {
	if s.ShelfContent == nil {
		s.ShelfContent = map[string][]StorageShelfContent{} // TODO/FIXME ! panic without but is initialised in Open function ??
	}

	if _, exists := s.Shelf[shelfName]; !exists {
		s.ShelfContent[shelfName] = make([]StorageShelfContent, 0)
	}

	found := false
	for idx := range s.ShelfContent[shelfName] {
		if s.ShelfContent[shelfName][idx].ContentID == fn {
			s.ShelfContent[shelfName][idx].IsDeleted = isDeleted

			found = true
			break
		}
	}

	if !found {
		s.ShelfContent[shelfName] = append(s.ShelfContent[shelfName], StorageShelfContent{
			ShelfID:   shelfName,
			ContentID: fn,
			IsDeleted: isDeleted,
		})
	}
}

func (s *JSONStorage) Shelfs() []StorageShelf {
	result := make([]StorageShelf, len(s.Shelf))

	for idx := range s.Shelf {
		result = append(result, s.Shelf[idx])
	}

	return result
}

func (s *JSONStorage) ShelfContents(shelfName string) []StorageShelfContent {
	result := make([]StorageShelfContent, len(s.ShelfContent[shelfName]))
	result = append(result, s.ShelfContent[shelfName]...)

	return result
}

func (s *JSONStorage) AddBookmark(bID, vID, cID, typeStr, path string, index, startOffset, endOffset int, text, annotation string, created, modified time.Time) {
	createdStr := created.Format(StorageTimeFmt)
	modifiedStr := modified.Format(StorageTimeFmt)

	found := false
	for bIdx := range s.Bookmark[cID] {
		if s.Bookmark[cID][bIdx].ID == bID && s.Bookmark[cID][bIdx].Modified == modifiedStr {
			found = true
		}
	}

	if !found {
		s.Bookmark[cID] = append(s.Bookmark[cID], StorageBookmark{
			ID:          bID,
			VolumeID:    vID,
			ContentID:   cID,
			BookPath:    path,
			Index:       index,
			StartOffset: startOffset,
			EndOffset:   endOffset,
			Text:        text,
			Annotation:  annotation,
			Created:     createdStr,
			Modified:    modifiedStr,
			Type:        typeStr,
		})
	}
}

func (s *JSONStorage) Bookmarks(cID string) []StorageBookmark {
	result := make([]StorageBookmark, 0)
	result = append(result, s.Bookmark[cID]...)

	return result
}
