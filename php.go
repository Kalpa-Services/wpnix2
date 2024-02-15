package main

import (
	"fmt"
	"os"
	"os/exec"
)

func checkAndInstallPHP() {
	if _, err := exec.LookPath("php8.2"); err != nil {
		fmt.Println("PHP 8.2 is not installed. Installing PHP 8.2 and PHP 8.2-FPM...")
		if err := exec.Command("apt", "update").Run(); err != nil {
			fmt.Fprintln(os.Stderr, "\x1b[31mError updating package lists:", err, "\x1b[0m")
			return
		}
		if err := exec.Command("apt", "install", "software-properties-common", "-y").Run(); err != nil {
			fmt.Fprintln(os.Stderr, "\x1b[31mError installing software-properties-common:", err, "\x1b[0m")
			return
		}
		if err := exec.Command("add-apt-repository", "ppa:ondrej/php", "-y").Run(); err != nil {
			fmt.Fprintln(os.Stderr, "\x1b[31mError adding PHP repository:", err, "\x1b[0m")
			return
		}
		if err := exec.Command("apt", "update").Run(); err != nil {
			fmt.Fprintln(os.Stderr, "\x1b[31mError updating package lists after adding repository:", err, "\x1b[0m")
			return
		}
		if err := exec.Command("apt", "install", "php8.2", "php8.2-fpm", "php8.2-mysql", "php8.2-xml", "php8.2-mbstring", "php8.2-curl", "php8.2-zip", "-y").Run(); err != nil {
			fmt.Fprintln(os.Stderr, "\x1b[31mError installing PHP 8.2 and modules:", err, "\x1b[0m")
			return
		}
		fmt.Println("\x1b[32mPHP 8.2 and PHP 8.2-FPM installed successfully.\x1b[0m")
	} else {
		fmt.Println("\x1b[33mPHP 8.2 and PHP 8.2-FPM are already installed.\x1b[0m")
	}
}
