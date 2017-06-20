package main

import (
	"encoding/json"
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

type album struct {
	Name      string `json:"name"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

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

func match(filePath string) bool {
	exts := []string{"cr2", "jpg"}

	for _, ext := range exts {
		if strings.ToLower(filepath.Ext(filePath)) == "."+strings.ToLower(ext) {
			return true
		}
	}

	return false
}

func action(srcFilePath string, dstDirRoot string, albums []album) filesmurf.ActionFunc {
	return func(filePath string) {
		dstPath := getDstPath(srcFilePath, dstDirRoot, albums)
		if dstPath != "" {
			pathTree := strings.Split(dstPath, "/")
			os.MkdirAll(strings.Join(pathTree[:len(pathTree)-1], "/"), os.ModePerm)
			os.Rename(filePath, dstPath)
		}
	}
}

func main() {
	srcPathRoot := os.Args[1]
	dstPathRoot := os.Args[2]
	albumPathRoot := "./albums.json"

	if len(os.Args) == 4 {
		albumPathRoot = os.Args[3]
	}

	albums := parseAlbumConf(albumPathRoot)
	filesmurf.Run(srcPathRoot, match, action(srcPathRoot, dstPathRoot, albums))
}
