package utils

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	startCmd, _ := filepath.Abs("./" + startCommand)

	dedicatedPtr := flag.Bool("dedicated", false, "Runs the dedicated server without needing to open offline(.exe).")
	flag.Parse()

	if *dedicatedPtr {
		args = append(args, "--dedicated")
	}

	cmd := exec.Command(startCmd)
	if args != nil {
		cmd = exec.Command(startCmd, strings.Join(args[:], " "))
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
	}
}

func Contains(s []string, searchTerm string) bool {
	sort.Strings(s)
	i := sort.SearchStrings(s, searchTerm)
	return i < len(s) && s[i] == searchTerm
}
