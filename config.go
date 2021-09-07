package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/tvrzna/go-utils/args"
	"github.com/tvrzna/pkgtray/checker"
)

type pkgManager struct {
	name    string
	command string
	pkgType checker.EnPkgManager
}

type config struct {
	intervalSeconds int
	pkgmanager      pkgManager
	pkgmanagerPath  string
}

func loadConfig() *config {
	xbpsPkgMngr := pkgManager{"xbps", "xbps-install", checker.Xbps}
	pacmanPkgMngr := pkgManager{"pacman", "checkupdates", checker.Pacman}
	apkPkgMngr := pkgManager{"apk", "apk", checker.Apk}
	aptGetPkgManager := pkgManager{"apt-get", "apt-get", checker.Apt_get}

	pkgManagers := []pkgManager{xbpsPkgMngr, pacmanPkgMngr, apkPkgMngr}

	conf := &config{intervalSeconds: 3600}

	if args.ContainsArg(os.Args, "-h", "--help") {
		fmt.Println("Usage: pkgtray [options]")
		fmt.Println("Options:")
		fmt.Printf("  -h, --help\t\tprint this help\n")
		fmt.Printf("  -v, --version\t\tprint version\n")
		fmt.Printf("  -t, --time NUMBER\ttime interval in seconds to check package manager for new updates (default 3600)\n")
		fmt.Printf("  -x, --xbps\t\tuse XBPS as desired package manager\n")
		fmt.Printf("  -p, --pacman\t\tuse Pacman as desired package manager\n")
		fmt.Printf("  -a, --apk\t\tuse Apk as desired package manager\n")
		fmt.Printf("  -g, --apt-get\t\tuse apt-get desired package manager\n")
		os.Exit(0)
	}

	if args.ContainsArg(os.Args, "-v", "--version") {
		fmt.Println(getVersion())
		os.Exit(0)
	}

	args.ParseArgs(os.Args, func(arg, nextArg string) {
		switch arg {
		case "-x", "--xbps":
			conf.pkgmanager = xbpsPkgMngr
		case "-p", "--pacman":
			conf.pkgmanager = pacmanPkgMngr
		case "-a", "--apk":
			conf.pkgmanager = apkPkgMngr
		case "-g", "--apt-get":
			conf.pkgmanager = aptGetPkgManager
		case "-t", "--time":
			i, err := strconv.Atoi(nextArg)
			if err != nil {
				log.Printf("Unexpected time interval %s", nextArg)
			} else {
				conf.intervalSeconds = i
			}
		}
	})

	if conf.pkgmanager.command == "" {
		log.Print("No package manager selected.")
		for _, pkgManager := range pkgManagers {
			path, err := exec.LookPath(pkgManager.command)
			if err == nil {
				log.Printf("Autoselected %s.", pkgManager.name)
				conf.pkgmanager = pkgManager
				conf.pkgmanagerPath = path
				break
			}
		}
		if conf.pkgmanager.name == "" {
			log.Println("No supported package manager found!")
			os.Exit(1)
		}
	} else {
		log.Printf("Selected '%s' as package manager.", conf.pkgmanager.name)
	}

	if conf.pkgmanagerPath == "" {
		path, err := exec.LookPath(conf.pkgmanager.command)
		if err == nil {
			conf.pkgmanagerPath = path
		} else {
			conf.pkgmanagerPath = conf.pkgmanager.command
		}
	}

	return conf
}
