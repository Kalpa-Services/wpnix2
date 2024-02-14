package main

import (
	"fmt"
	"os/exec"
)

func checkAndInstallPerl() {
	if _, err := exec.LookPath("perl"); err != nil {
		fmt.Println("Perl is not installed. Installing Perl...")
		exec.Command("apt", "install", "perl", "-y").Run()
	} else {
		fmt.Println("Perl is already installed.")
	}
}
