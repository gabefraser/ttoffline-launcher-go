package patcher

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
	"time"

	"toontown-offline-launcher/utils"

	"github.com/mholt/archiver/v3"
)

var patcher map[string]interface{}
var files Files
var patchingURL = "https://releases.toontownoffline.net"
var patcherURL = patchingURL + "/%s.json"
var executables = []string{"ToontownOffline", "astrond-linux", "astrond-darwin", "offline", "offline.exe"}

func PatchFiles() {
	for _, file := range files.GetFiles() {
		if _, err := os.Stat(file.GetFullFilePath()); os.IsNotExist(err) {
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

func ParsePatcher() {
	url := fmt.Sprintf(patcherURL, utils.GetRuntimePlatform())

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

					files.AddFile(*file)
				}
			}
		}
	}
}

func downloadFile(file File) error {
	if _, err := os.Stat(file.GetFilePath()); os.IsNotExist(err) {
		os.MkdirAll(file.GetFilePath(), os.ModePerm)
	}

	out, err := os.Create(file.GetFullFilePath() + ".bz2")
	if err != nil {
		return err
	}

	resp, err := http.Get(file.GetURLPath())
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

	decompressBzip2(file.GetFullFilePath()+".bz2", file.GetFullFilePath())

	if utils.Contains(executables, file.Name) {
		err := os.Chmod(file.GetFullFilePath(), 755)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func decompressBzip2(filePath string, fileName string) {
	err := archiver.DecompressFile(filePath, fileName)
	if err != nil {
		panic(err)
	}

	err = os.Remove(filePath)
	if err != nil {
		panic(err)
	}
}

func getFileHash(file File) (string, error) {
	var md5Hash string

	curFile, err := os.Open(file.GetFullFilePath())
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
