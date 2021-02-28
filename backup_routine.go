package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"time"
)

func timeStamp() string {
	now := time.Now()
	return fmt.Sprintf("%d-%02d-%02d_%02d-%02d-%02d",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
	)
}

func backupDatabase(ctx context.Context, cs *CSet, cfg *Config) {
	// database backups
	ticker := time.NewTicker(cfg.BackupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Closing backup routine...")

			if time.Now().Before(startupTime.Add(cfg.DurationBeforeFirstBackup)) {
				return
			}

			log.Println("Creating backup...")
			filename := timeStamp() + ".csv"
			filename = path.Join(cfg.DataPath, filename)
			data, err := cs.DumpCSV()
			if err != nil {
				log.Printf("Failed to retrieve data for backup: %v\n", err)
				continue
			}
			err = ioutil.WriteFile(filename, data, 0660)
			if err != nil {
				log.Printf("Failed to write data to file '%s': %v", filename, err)
				continue
			}
			log.Printf("Created backup: %s\n", filename)

			return
		case <-ticker.C:
			filename := timeStamp() + ".csv"
			filename = path.Join(cfg.DataPath, filename)
			data, err := cs.DumpCSV()
			if err != nil {
				log.Printf("Failed to retrieve data for backup: %v\n", err)
				continue
			}
			err = ioutil.WriteFile(filename, data, 0660)
			if err != nil {
				log.Printf("Failed to write data to file '%s': %v", filename, err)
				continue
			}
			log.Printf("Created backup: %s\n", filename)

		}
	}

}
