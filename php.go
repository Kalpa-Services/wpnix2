package main

import (
	"fmt"
	"os/exec"
)

func checkAndInstallPHP() {
	if _, err := exec.LookPath("php8.2"); err != nil {
		fmt.Println("PHP 8.2 is not installed. Installing PHP 8.2 and PHP 8.2-FPM...")
		exec.Command("apt", "update").Run()
		exec.Command("apt", "install", "software-properties-common").Run()
		exec.Command("add-apt-repository", "ppa:ondrej/php", "-y").Run()
		exec.Command("apt", "update").Run()
		exec.Command("apt", "install", "php8.2", "php8.2-fpm", "php8.2-mysql", "php8.2-xml", "php8.2-mbstring", "php8.2-curl", "php8.2-zip", "-y").Run()
	} else {
		fmt.Println("PHP 8.2 and PHP 8.2-FPM are already installed.")
	}
}
