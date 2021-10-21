package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

type File struct {
	Name string
	Hash string
	Path string
}

type Files struct {
	List []File
}

type DownloadTotal struct {
	Progress uint64
	Filename string
}

func (dt *DownloadTotal) Write(p []byte) (int, error) {
	n := len(p)
	dt.Progress += uint64(n)
	dt.PrintProgress(dt.Filename)
	return n, nil
}

func (dt DownloadTotal) PrintProgress(filename string) {
	fmt.Printf("\r%s", strings.Repeat(" ", 35))
	fmt.Printf("\rDownloading %s... %s complete", filename, humanize.Bytes(dt.Progress))
}

func (f *Files) addFile(file File) []File {
	f.List = append(f.List, file)
	return f.List
}

func (f *Files) getFiles() []File {
	return f.List
}

func (f *File) getFullFilePath() string {
	if len(f.Path) > 0 {
		sprintf := "%s/%s"
		path := f.Path

		if getRuntimePlatform() == "windows" {
			sprintf = "%s\\%s"
			path = strings.Replace(path, "/", "\\", -1)
		}

		return fmt.Sprintf(sprintf, path, f.Name)
	} else {
		return f.Name
	}
}

func (f *File) getFilePath() string {
	if len(f.Path) > 0 {
		path := f.Path

		if getRuntimePlatform() == "windows" {
			path = strings.Replace(path, "/", "\\", -1)
		}

		return path
	}

	return ""
}

func (f *File) getURLPath() string {
	if len(f.Path) > 0 {
		path := f.Path

		if getRuntimePlatform() == "windows" {
			path = strings.Replace(path, "/", "\\", -1)
		}

		return fmt.Sprintf("%s/%s/%s.bz2", patchingURL, f.Path, f.Name)
	}

	return fmt.Sprintf("%s/%s.bz2", patchingURL, f.Name)
}

var patcher map[string]interface{}
var files Files
var patchingURL = "https://releases.toontownoffline.net"
var patcherURL = patchingURL + "/%s.json"

func main() {
	url := fmt.Sprintf(patcherURL, getRuntimePlatform())

	parsePatcher(url)

	patchFiles()

	bootGame()
}

func bootGame() {
	fmt.Println("Booting the game...")
}

func patchFiles() {
	for _, file := range files.getFiles() {
		if _, err := os.Stat(file.getFullFilePath()); os.IsNotExist(err) {
			err := downloadFile(file)
			if err != nil {
				panic(err)
			}
		} else if os.IsExist(err) {
			hash, err := getFileHash(file)
			if err != nil {
				panic(err)
			}

			fmt.Println(fmt.Sprintf("File: %s, current hash: %s, patcher hash: %s", file.Name, hash, file.Hash))

			if hash != file.Hash {
				err := downloadFile(file)
				if err != nil {
					panic(err)
				}
			}
		}
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

func downloadFile(file File) error {
	if _, err := os.Stat(file.getFilePath()); os.IsNotExist(err) {
		os.MkdirAll(file.getFilePath(), os.ModePerm)
	}

	out, err := os.Create(file.getFullFilePath() + ".tmp")
	if err != nil {
		return err
	}

	resp, err := http.Get(file.getURLPath())
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()

	counter := &DownloadTotal{
		Filename: file.Name,
	}
	if _, err := io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}

	fmt.Println("")

	out.Close()

	if err := os.Rename(file.getFullFilePath()+".tmp", file.getFullFilePath()); err != nil {
		return err
	}

	return nil
}

func getFileHash(file File) (string, error) {
	var md5Hash string

	curFile, err := os.Open(file.getFullFilePath())
	if err != nil {
		return md5Hash, err
	}

	defer curFile.Close()

	hash := md5.New()

	if _, err := io.Copy(hash, curFile); err != nil {
		return md5Hash, err
	}

	hashInBytes := hash.Sum(nil)[:16]
	md5Hash = hex.EncodeToString(hashInBytes)

	return md5Hash, nil
}

func getRuntimePlatform() string {
	patchingPlatform := runtime.GOOS
	if patchingPlatform == "darwin" {
		return "mac"
	}

	return patchingPlatform
}
