//go:generate mockgen -package pkg -destination storage_mock.go -source storage.go
package pkg

import (
	"encoding/json"
	"os"
	"time"
)

type Storage interface {
	Save() error

	AddContent(fn, title, author, url string, words int, book, finished bool)
	AddDevice(device, model string)
	AddEvent(fn, device, name string, t time.Time, duration int)

	Contents() []StorageContent
	Events(cID string) []StorageEvents
}

type JSONStorage struct {
	DeviceMap  map[string]StorageDevice   `json:"devices"`
	ContentMap map[string]StorageContent  `json:"contents"`
	EventMap   map[string][]StorageEvents `json:"events"`

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

const (
	StorageTimeFmt = "2006-01-02T15:04:05.000"
)

func OpenStorageOrCreate(fn string) (Storage, error) {
	storage := JSONStorage{
		DeviceMap:  map[string]StorageDevice{},
		ContentMap: map[string]StorageContent{},
		EventMap:   map[string][]StorageEvents{},
		fn:         fn,
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

func (s *JSONStorage) AddContent(fn, title, author, url string, words int, book, finished bool) {
	s.ContentMap[fn] = StorageContent{
		ID:         fn,
		Title:      title,
		Author:     author,
		Words:      words,
		URL:        url,
		IsBook:     book,
		IsFinished: finished,
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
