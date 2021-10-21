package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"runtime"
	"time"
)

type File struct {
	Name string
	Hash string
	Path string
}

type Files struct {
	List []File
}

func (f *Files) addFile(file File) []File {
	f.List = append(f.List, file)
	return f.List
}

func (f *Files) getFiles() []File {
	return f.List
}

func (f *File) getFilePath() string {
	return fmt.Sprintf("%s/%s", f.Path, f.Name)
}

var patcher map[string]interface{}
var files Files
var patchingURL = "https://releases.toontownoffline.net/%s.json"

func main() {
	url := fmt.Sprintf(patchingURL, getRuntimePlatform())

	parsePatcher(url)

	for _, file := range files.getFiles() {
		fmt.Println(file.getFilePath())
	}
}

func parsePatcher(url string) {
	patchClient := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Toontown Offline Launcher")

	res, getErr := patchClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	jsonErr := json.Unmarshal(body, &patcher)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	generateFiles(patcher)
}

func generateFiles(m map[string]interface{}) {
	for _, v := range m {
		if mv, ok := v.(map[string]interface{}); ok {
			for filename, s := range mv {
				if mv, ok := s.(map[string]interface{}); ok {
					fn, _ := regexp.Compile("([^\\/]+$)")

					fileName := fn.FindString(filename)
					path := fmt.Sprintf("%v", mv["path"])
					hash := fmt.Sprintf("%v", mv["hash"])

					file := &File{
						Name: fileName,
						Path: path,
						Hash: hash,
					}

					files.addFile(*file)
				}
			}
		}
	}
}

func getRuntimePlatform() string {
	patchingPlatform := runtime.GOOS
	if patchingPlatform == "darwin" {
		return "mac"
	}

	return patchingPlatform
}
