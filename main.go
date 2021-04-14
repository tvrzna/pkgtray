package main

import (
	"fmt"

	"github.com/tvrzna/pkgtray/checker"
)

const version = "0.1.1"

var conf *config

func main() {
	conf = loadConfig()

	trayIcon := createTrayIcon(conf, func() int {
		return checker.CheckPackages(conf.pkgmanager.pkgType, conf.pkgmanagerPath)
	})
	trayIcon.processRefreshLoop()
	trayIcon.start()
}

func getVersion() string {
	return fmt.Sprintf("pkgtray %s\nhttps://github.com/tvrzna/pkgtray\n\nReleased under the MIT License.", version)
}
