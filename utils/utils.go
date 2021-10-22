package utils

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func GetRuntimePlatform() string {
	patchingPlatform := runtime.GOOS
	if patchingPlatform == "darwin" {
		return "mac"
	}

	return patchingPlatform
}

func BootGame(args ...string) (p *os.Process, err error) {
	fmt.Println("Booting the game...")

	dedicatedPtr := flag.Bool("dedicated", false, "Runs the dedicated server without needing to open offline(.exe).")
	flag.Parse()

	if *dedicatedPtr {
		args = append(args, "--dedicated")
	}

	if args[0], err = exec.LookPath(args[0]); err == nil {
		var procAttr os.ProcAttr

		procAttr.Files = []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		}

		p, err := os.StartProcess(args[0], args, &procAttr)
		if err != nil {
			return p, nil
		}
	}

	return nil, err
}
