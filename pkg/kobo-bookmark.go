package pkg

import "time"

type KoboBookmark struct {
	ID          string
	VolumeID    string
	ContentID   string
	BookPath    string
	Index       int
	StartOffset int
	EndOffset   int
	Text        string
	Annotation  string
	Created     time.Time
	Modified    time.Time
	Type        string
}
