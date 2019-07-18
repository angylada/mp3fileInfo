package main

import (
	"fmt"
	"github.com/bogem/id3v2"
	id3v22 "github.com/bogem/id3v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

func main() {

	wd, _ := os.Getwd()

	files, err := ioutil.ReadDir(".")

	if err != nil {
		log.Fatal(err)
	}

	var filteredFiles []string


	for _, f := range files {
		if filepath.Ext(f.Name()) == ".mp3" {
			filteredFiles = append(filteredFiles, f.Name())
		}
	}

	if len(filteredFiles) < 1 {
		fmt.Println("Nothing to be done")
		return
	}
	worker := New("([^ ]| [^-]| -[^ ])+[^\\.mp3]")

	var wg sync.WaitGroup


	wg.Add(len(filteredFiles))
	for _, mp3 := range filteredFiles {
		go worker.addMetaData(mp3, wd, &wg)
	}

	wg.Wait()
	fmt.Println("DONE")
}

type Worker struct {
	rgxp regexp.Regexp
}

func New(regex string) *Worker {
	r := regexp.MustCompile(regex)
	w := Worker{
		rgxp: *r,
	}

	return &w
}

func (w Worker) addMetaData(fileName string, wd string, wg *sync.WaitGroup) {
	defer wg.Done()
	groups := w.rgxp.FindAllString(fileName, 2)
	if groups == nil || len(groups) < 2  {
		return
	}

	fmt.Println(wd + "/" + fileName)
	mp3File, mp3Err := id3v22.Open(wd + "/" + fileName, id3v2.Options{Parse: true})

	if mp3Err != nil {
		fmt.Println(mp3Err)
		return
	}
	defer mp3File.Close()
	fmt.Println(mp3File.Artist())

	if len(mp3File.Artist()) < 1 {
		fmt.Println("'" + mp3File.Artist() + "'" + groups[0])
		mp3File.SetArtist(groups[0])
	}

	if len(mp3File.Title()) < 1 {
		mp3File.SetTitle(strings.Trim(groups[1], "- "))
	}
	_ = mp3File.Save()
}

/*
((?! - ).)+[^\.mp3] !!!!
=> ([^ ]| [^-]| -[^ ])+[^\.mp3]

 */