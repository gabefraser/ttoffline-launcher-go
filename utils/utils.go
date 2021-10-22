package utils

import (
	"flag"
	"fmt"
	"os/exec"
	"runtime"
	"sort"
	"strings"
)

func GetRuntimePlatform() string {
	patchingPlatform := runtime.GOOS
	if patchingPlatform == "darwin" {
		return "mac"
	}

	return patchingPlatform
}

func BootGame(startCommand string) {
	fmt.Println("Booting the game...")

	var args []string

	dedicatedPtr := flag.Bool("dedicated", false, "Runs the dedicated server without needing to open offline(.exe).")
	flag.Parse()

	if *dedicatedPtr {
		args = append(args, "--dedicated")
	}

	cmd := exec.Command(startCommand, strings.Join(args[:], " "))

	cmd.Start()
}

func Contains(s []string, searchTerm string) bool {
	sort.Strings(s)
	i := sort.SearchStrings(s, searchTerm)
	return i < len(s) && s[i] == searchTerm
}
