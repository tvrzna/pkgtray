package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type enPkgManager int

const (
	xbps enPkgManager = iota
	pacman
	apk
	apt_get
)

type pkgManager struct {
	name    string
	command string
	pkgType enPkgManager
}

type config struct {
	intervalSeconds int
	pkgmanager      pkgManager
	pkgmanagerPath  string
}

func loadConfig() *config {
	xbpsPkgMngr := pkgManager{"xbps", "xbps-install", xbps}
	pacmanPkgMngr := pkgManager{"pacman", "checkupdates", pacman}
	apkPkgMngr := pkgManager{"apk", "apk", apk}
	aptGetPkgManager := pkgManager{"apt-get", "apt-get", apt_get}

	pkgManagers := []pkgManager{xbpsPkgMngr, pacmanPkgMngr, apkPkgMngr}

	conf := &config{}

	var useXbps bool
	var usePacman bool
	var useApk bool
	var useAptGet bool
	var interval int
	var printHelp bool
	var printVersion bool

	flag.BoolVar(&printHelp, "h", false, "print this help")
	flag.BoolVar(&printHelp, "-help", false, "print this help")

	flag.BoolVar(&printVersion, "v", false, "print version")
	flag.BoolVar(&printVersion, "-version", false, "print version")

	flag.BoolVar(&useXbps, "x", false, "use XBPS as desired package manager")
	flag.BoolVar(&useXbps, "-xbps", false, "use XBPS as desired package manager")

	flag.BoolVar(&usePacman, "p", false, "use Pacman as desired package manager")
	flag.BoolVar(&usePacman, "-pacman", false, "use Pacman as desired package manager")

	flag.BoolVar(&useApk, "a", false, "use Apk as desired package manager")
	flag.BoolVar(&useApk, "-apk", false, "use Apk as desired package manager")

	flag.BoolVar(&useAptGet, "g", false, "use apt-get as desired package manager")
	flag.BoolVar(&useAptGet, "-apt-get", false, "use apt-get as desired package manager")

	flag.IntVar(&interval, "t", 3600, "time interval to check package manager for new updates")
	flag.IntVar(&interval, "-time", 3600, "time interval to check package manager for new updates")

	flag.Parse()

	if printHelp {
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

	if printVersion {
		fmt.Println(getVersion())
		os.Exit(0)
	}

	if useXbps {
		conf.pkgmanager = xbpsPkgMngr
	} else if usePacman {
		conf.pkgmanager = pacmanPkgMngr
	} else if useApk {
		conf.pkgmanager = apkPkgMngr
	} else if useAptGet {
		conf.pkgmanager = aptGetPkgManager
	} else {
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
	}

	if conf.pkgmanagerPath == "" {
		path, err := exec.LookPath(conf.pkgmanager.command)
		if err == nil {
			conf.pkgmanagerPath = path
		} else {
			conf.pkgmanagerPath = conf.pkgmanager.command
		}
	}

	if interval > 0 {
		conf.intervalSeconds = interval
	} else {
		conf.intervalSeconds = 3600
	}

	return conf
}
