package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
	"fmt"
	"io/ioutil"
)

func checkError(e error) {
    if e != nil {
        log.Fatal(e)
    }
}

func crossBuildStart() {
	if _, err := os.Stat("/bin/sh.real"); os.IsNotExist(err) {
		err = os.Link("/bin/sh", "/bin/sh.real")
		checkError(err)
	}

	err := os.Remove("/bin/sh")
	checkError(err)

	err = os.Link("/usr/bin/resin-xbuild", "/bin/sh")
	checkError(err)

	if _, err := os.Stat("/.balena/image-info"); err == nil {
		info, err := ioutil.ReadFile("/.balena/image-info")
		checkError(err)

		fmt.Print(string(info))

		err = os.Rename("/.balena/image-info", "/.balena/image-info_")
		checkError(err)
	}
}

func crossBuildEnd() {
	err := os.Remove("/bin/sh")
	checkError(err)

	err = os.Link("/bin/sh.real", "/bin/sh")
	checkError(err)
}

func runShell() error {
	var cmd *exec.Cmd

	if _, err := os.Stat("/usr/bin/qemu-arm-static"); err == nil {
		cmd = exec.Command("/usr/bin/qemu-arm-static", append([]string{"-execve", "-0", os.Args[0], "/bin/sh"}, os.Args[1:]...)...)
	} else {
		cmd = exec.Command("/usr/bin/qemu-aarch64-static", append([]string{"-execve", "-0", os.Args[0], "/bin/sh"}, os.Args[1:]...)...)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	switch os.Args[0] {
	case "cross-build-start":
		crossBuildStart()
	case "cross-build-end":
		crossBuildEnd()
	default:
		code := 0
		crossBuildEnd()

		if err := runShell(); err != nil {
			code = 1
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					code = status.ExitStatus()
				}
			}
		}

		crossBuildStart()

		os.Exit(code)
	}
}
