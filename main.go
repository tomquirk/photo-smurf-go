package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/djherbis/times"
	"github.com/tomquirk/filesmurf"
)

// album structs represent a single album (collection of images)
type album struct {
	Name      string `json:"name"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// parseAlbumConf parses an album configuration file at a given location and returns
// an array of album structs
func parseAlbumConf(confFilePath string) []album {
	var albums []album

	file, e := ioutil.ReadFile(confFilePath)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	json.Unmarshal(file, &albums)

	return albums
}

// getDstPath returns the appropriate destination path for a file at a given file path
func getDstPath(filePath string, dstDirRoot string, albums []album) string {
	fileStat, err := times.Stat(filePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	birthTime := fileStat.BirthTime()
	filePathTree := strings.Split(filePath, "/")

	for _, album := range albums {
		start, _ := time.Parse(time.RFC822, album.StartTime)
		end, _ := time.Parse(time.RFC822, album.EndTime)

		if birthTime.After(start) && birthTime.Before(end) {
			return fmt.Sprintf("%s%s/%s", dstDirRoot, album.Name, filePathTree[len(filePathTree)-1])
		}
	}

	return ""
}

// matchImage determines whether the file at a given file path is an image
func matchImage(filePath string) bool {
	exts := []string{"cr2", "jpg"}

	for _, ext := range exts {
		if strings.ToLower(filepath.Ext(filePath)) == "."+strings.ToLower(ext) {
			return true
		}
	}

	return false
}

// moveImage moves a file at a given file path to a given destination
func moveImage(srcDirPath string, dstDirPath string, albums []album) filesmurf.ActionFunc {
	return func(filePath string) error {
		dstFilePath := getDstPath(filePath, dstDirPath, albums)
		if dstFilePath == "" {
			return errors.New("dest path cannot be empty")
		}

		filesmurf.Move(filePath, dstFilePath)
		return nil
	}
}

func main() {
	srcDirPath := os.Args[1]
	dstDirPath := os.Args[2]
	albumPathRoot := "./albums.json"

	if len(os.Args) == 4 {
		albumPathRoot = os.Args[3]
	}

	albums := parseAlbumConf(albumPathRoot)
	filesmurf.Run(srcDirPath, matchImage, moveImage(srcDirPath, dstDirPath, albums))
}
