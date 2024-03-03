package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/timchurchard/kobo-readstat/pkg"
)

// Sync command reads a Kobo database and creates/updates local storage
func Sync(out io.Writer) int {
	const (
		defaultEmpty   = ""
		defaultStorage = "./readstat.json"

		usageDatabasePath = "Path to /media/kobo/.kobo/KoboReader.sqlite"
		usageStoragePath  = "Path to local storage default: " + defaultStorage
	)
	var (
		databaseFn string
		storageFn  string
	)

	flag.StringVar(&databaseFn, "database", defaultEmpty, usageDatabasePath)
	flag.StringVar(&databaseFn, "d", defaultEmpty, usageDatabasePath)

	flag.StringVar(&storageFn, "storage", defaultStorage, usageStoragePath)
	flag.StringVar(&storageFn, "s", defaultStorage, usageStoragePath)

	flag.Usage = func() {
		fmt.Fprintf(out, "Usage of %s %s:\n", os.Args[0], os.Args[1])

		flag.PrintDefaults()
	}

	flag.Parse()

	if databaseFn == "" {
		fmt.Println("-d or --database /path/to/KoboReader.sqlite is required.")
		return 1
	}

	// Read data from Kobo DB
	db, err := pkg.NewKoboDatabase(databaseFn)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Create/Update Storage
	storage, err := pkg.OpenStorageOrCreate(storageFn)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := storage.Save(); err != nil {
			fmt.Printf("Error saving: %v\n", err)
		}
	}()

	// Do the sync!
	err = pkg.Sync(db, storage)
	if err != nil {
		fmt.Printf("Error syncing: %v\n", err)
		return 1
	}

	return 0
}
