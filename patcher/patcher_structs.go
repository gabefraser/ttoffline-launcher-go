package patcher

import (
	"fmt"
	"strings"

	"toontown-offline-launcher/utils"

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

func (f *Files) AddFile(file File) []File {
	f.List = append(f.List, file)
	return f.List
}

func (f *Files) GetFiles() []File {
	return f.List
}

func (f *File) GetFullFilePath() string {
	if len(f.Path) > 0 {
		sprintf := "%s/%s"
		path := f.Path

		if utils.GetRuntimePlatform() == "windows" {
			sprintf = "%s\\%s"
			path = strings.Replace(path, "/", "\\", -1)
		}

		return fmt.Sprintf(sprintf, path, f.Name)
	} else {
		return f.Name
	}
}

func (f *File) GetFilePath() string {
	if len(f.Path) > 0 {
		path := f.Path

		if utils.GetRuntimePlatform() == "windows" {
			path = strings.Replace(path, "/", "\\", -1)
		}

		return path
	}

	return ""
}

func (f *File) GetURLPath() string {
	if len(f.Path) > 0 {
		path := f.Path

		if utils.GetRuntimePlatform() == "windows" {
			path = strings.Replace(path, "/", "\\", -1)
		}

		return fmt.Sprintf("%s/%s/%s.bz2", patchingURL, f.Path, f.Name)
	}

	return fmt.Sprintf("%s/%s.bz2", patchingURL, f.Name)
}
