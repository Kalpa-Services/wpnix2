package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func checkAndInstallCertbot() {
	if _, err := exec.LookPath("certbot"); err != nil {
		fmt.Println("\x1b[33mCertbot is not installed. Installing Certbot using Snap...\x1b[0m")
		if err := exec.Command("apt", "update").Run(); err != nil {
			fmt.Fprintln(os.Stderr, "\x1b[31mError updating package lists:", err, "\x1b[0m")
			return
		}
		if err := exec.Command("apt", "install", "snapd", "-y").Run(); err != nil {
			fmt.Fprintln(os.Stderr, "\x1b[31mError installing snapd:", err, "\x1b[0m")
			return
		}
		if err := exec.Command("snap", "install", "core").Run(); err != nil {
			fmt.Fprintln(os.Stderr, "\x1b[31mError installing snap core:", err, "\x1b[0m")
			return
		}
		if err := exec.Command("snap", "refresh", "core").Run(); err != nil {
			fmt.Fprintln(os.Stderr, "\x1b[31mError refreshing snap core:", err, "\x1b[0m")
			return
		}
		if err := exec.Command("snap", "install", "--classic", "certbot").Run(); err != nil {
			fmt.Fprintln(os.Stderr, "\x1b[31mError installing Certbot:", err, "\x1b[0m")
			return
		}
		if err := exec.Command("ln", "-s", "/snap/bin/certbot", "/usr/bin/certbot").Run(); err != nil {
			fmt.Fprintln(os.Stderr, "\x1b[31mError creating symlink for Certbot:", err, "\x1b[0m")
			return
		}
		fmt.Println("\x1b[32mCertbot installed successfully.\x1b[0m")
	} else {
		fmt.Println("\x1b[33mCertbot is already installed.\x1b[0m")
	}
}

func configureLetsEncryptSSL(domain string, email string) error {
	fmt.Println("Configuring Let's Encrypt SSL for", domain, "...")
	var cmd *exec.Cmd
	domainArgs := []string{"--nginx", "--non-interactive", "--agree-tos", "--email", email, "-d", domain}
	// Check if domain is a subdomain and include both www and non-www versions
	if strings.Count(domain, ".") == 1 {
		domainArgs = append(domainArgs, "-d", "www."+domain)
	}
	// Specify the nginx server block configuration file to ensure the correct one is used
	serverBlockPath := fmt.Sprintf("/etc/nginx/sites-available/%s", domain)
	domainArgs = append(domainArgs, "--nginx-server-root", serverBlockPath)

	cmd = exec.Command("certbot", domainArgs...)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error configuring Let's Encrypt SSL for %s: %w", domain, err)
	} else {
		fmt.Println("\x1b[32mSuccessfully configured Let's Encrypt SSL for", domain, "\x1b[0m")
	}
	return nil
}
