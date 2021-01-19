package main

import "fmt"

const version = "0.1.0"

var conf *config

func main() {
	conf = loadConfig()

	trayIcon := createTrayIcon(conf, func() int {
		return checkPackages(conf)
	})
	trayIcon.processRefreshLoop()
	trayIcon.start()
}

func getVersion() string {
	return fmt.Sprintf("pkgtray %s\nhttps://github.com/tvrzna/pkgtray\n\nReleased under the MIT License.", version)
}
