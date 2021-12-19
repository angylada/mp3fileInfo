package mp3fileInfo

import (
	"fmt"
	id3v22 "github.com/bogem/id3v2/v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

var pathSeparator = "/"

func main() {
	osArgs := os.Args[1:]
	wd, _ := os.Getwd()

	if len(osArgs) > 0 {
		if osArgs[0] == "help" || osArgs[0] == "-h" || osArgs[0] == "--help" {
			fmt.Println("Usage: mp3fileInfo [ <PATH-TO-DIR> ]")
			fmt.Println("if PATH-TO-DIR is omitted, the current working directory will be used (!== location of executable)")
		}

		file, err := os.Open(osArgs[0])
		if err != nil {
			fmt.Printf("Directory '%s' does not exist\n", osArgs[0])
			os.Exit(1)
		}
		fileInfo, err := file.Stat()
		if ! fileInfo.IsDir() {
			fmt.Printf("Provided path '%s' is not a directory\n", osArgs[0])
			os.Exit(1)
		}
		wd = osArgs[0]
	}
	if runtime.GOOS == "windows" {
		pathSeparator = "\\"
	}

	files, err := ioutil.ReadDir(wd)

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
	counter := 0
	for _, mp3 := range filteredFiles {
		go worker.addMetaData(mp3, wd, &wg, counter)
		counter++
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

func (w Worker) addMetaData(fileName string, wd string, wg *sync.WaitGroup, counter int) {
	defer wg.Done()
	groups := w.rgxp.FindAllString(fileName, 2)
	if groups == nil || len(groups) < 2  {
		return
	}

	//fmt.Println(counter + "'" + wd + pathSeparator + fileName + "'")
	fmt.Printf("%d '%s%s%s'\n", counter, wd, pathSeparator, fileName)
	mp3File, mp3Err := id3v22.Open(wd + "/" + fileName, id3v22.Options{Parse: true})

	if mp3Err != nil {
		fmt.Println(mp3Err)
		return
	}
	defer mp3File.Close()

	if len(mp3File.Artist()) < 1 {
		fmt.Printf("%d Artist: '%s%s'\n", counter, mp3File.Artist(), strings.TrimSpace(groups[0]))
		mp3File.SetArtist(strings.TrimSpace(groups[0]))
	}

	if len(mp3File.Title()) < 1 {
		fmt.Printf("%d Title: '%s'\n", counter, strings.Trim(groups[1], "- "))
		mp3File.SetTitle(strings.TrimSpace(strings.Trim(groups[1], "- ")))
	}
	_ = mp3File.Save()
}
