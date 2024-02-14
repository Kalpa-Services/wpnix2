package main

import (
	"fmt"
	"os"
	"os/exec"
)

func checkAndInstallPerl() {
	if _, err := exec.LookPath("perl"); err != nil {
		fmt.Println("Perl is not installed. Installing Perl...")
		if err := exec.Command("apt", "install", "perl", "-y").Run(); err != nil {
			fmt.Fprintln(os.Stderr, "\x1b[31mError installing Perl:", err, "\x1b[0m")
		} else {
			fmt.Println("\x1b[32mPerl installed successfully.\x1b[0m")
		}
	} else {
		fmt.Println("\x1b[33mPerl is already installed.\x1b[0m")
	}
}
