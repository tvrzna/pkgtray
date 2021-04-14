package checker

import (
	"bufio"
	"errors"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type EnPkgManager int

const (
	Xbps EnPkgManager = iota
	Pacman
	Apk
	Apt_get
)

func CheckPackages(pkgManager EnPkgManager, pkgManagerPath string) int {
	switch pkgManager {
	case Apk:
		return checkApk(pkgManagerPath)
	case Apt_get:
		return checkAptGet(pkgManagerPath)
	case Pacman:
		return checkPacman(pkgManagerPath)
	case Xbps:
		return checkXbps(pkgManagerPath)
	}
	return -1
}

func checkByLineCount(command string, method func(string)) error {
	arrExec := strings.Split(command, " ")

	cmd := exec.Command(arrExec[0], arrExec[1:]...)

	stdout, err := cmd.StdoutPipe()
	scanner := bufio.NewScanner(stdout)
	cmd.Start()
	for scanner.Scan() {
		method(scanner.Text())
	}
	cmd.Wait()

	if err != nil {
		return err
	}
	if scanner.Err() != nil {
		return errors.New("Could not perform check for updates")
	}
	return nil
}

func checkApk(path string) int {
	result := 0

	_, err := exec.Command(path, "update").Output()
	if err != nil {
		log.Print(err)
		return -1
	}

	err = checkByLineCount(path+" upgrade -s", func(line string) {
		if strings.HasPrefix(line, "(") && strings.HasSuffix(line, ")") {
			result++
		}
	})
	if err != nil {
		log.Print(err)
		return -1
	}

	return result
}

func checkAptGet(path string) int {
	result := 0

	_, err := exec.Command(path, "update").Output()
	if err != nil {
		log.Print(err)
		return -1
	}

	err = checkByLineCount(path+" -u upgrade --assume-no", func(line string) {
		if strings.Contains(line, "upgraded,") && strings.Contains(line, "newly installed,") {
			data := strings.Split(line, " ")
			count, err := strconv.Atoi(data[0])
			if err == nil {
				result = count
			} else {
				result = -1
			}
		}
	})
	if err != nil {
		log.Print(err)
		return -1
	}

	return result
}

func checkPacman(path string) int {
	result := 0

	err := checkByLineCount(path, func(line string) {
		result++
	})
	if err != nil {
		log.Print(err)
		return -1
	}

	return result
}

func checkXbps(path string) int {
	result := 0

	err := checkByLineCount(path+" -Munv", func(line string) {
		if strings.Contains(line, "Found") {
			result++
		}
	})
	if err != nil {
		log.Print(err)
		return -1
	}

	return result
}
