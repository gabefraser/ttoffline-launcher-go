package main

import (
	"toontown-offline-launcher/patcher"
	"toontown-offline-launcher/utils"
)

func main() {
	patcher.ParsePatcher()

	patcher.PatchFiles()

	if utils.GetRuntimePlatform() == "windows" {
		utils.BootGame("offline.exe")
	} else if utils.GetRuntimePlatform() == "mac" {
		utils.BootGame("ToontownOffline")
	} else if utils.GetRuntimePlatform() == "linux" {
		utils.BootGame("offline")
	}
}
