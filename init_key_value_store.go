package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jszwec/csvutil"
)

func initializeKeyValueStore(cs *CSet, dataFolderPath string) error {
	err := filepath.Walk(dataFolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !isValidSuffix(path) {
			return nil
		}
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return nil
		}
		list := []Reason{}
		err = csvutil.Unmarshal(b, &list)
		if err != nil {
			log.Printf("Failed to read: %s : %v\n", path, err)
			return nil
		}
		log.Printf("Read %d lines in: %s\n", len(list), path)
		for _, reason := range list {
			cs.AddFromCSV(reason)
		}
		return nil
	})
	return err
}

var suffix = []string{".csv"}

func isValidSuffix(path string) bool {
	for _, s := range suffix {
		if strings.HasSuffix(path, s) {
			return true
		}
	}
	return false
}
