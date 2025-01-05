package cmd

import "time"

const (
	defaultEmpty = ""

	defaultStorage   = "./readstat.json"
	usageStoragePath = "Path to local storage default: " + defaultStorage

	usageDatabasePath = "Path to /media/kobo/.kobo/KoboReader.sqlite"

	usageYear = "Year to generate stats for (default this year)"

	usageShowSessions = "Show reading sessions"
)

var defaultYear int

func init() {
	defaultYear = time.Now().Year()
}
