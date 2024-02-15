package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func configureLetsEncryptSSL(domain string, email string) error {
	fmt.Println("Configuring Let's Encrypt SSL for", domain, "...")
	var cmd *exec.Cmd
	if strings.Count(domain, ".") > 1 {
		cmd = exec.Command("certbot", "--nginx", "--nginx-server-root", nginxAvailable, "-d", domain, "--non-interactive", "--agree-tos", "--email", email)
	} else {
		cmd = exec.Command("certbot", "--nginx", "--nginx-server-root", nginxAvailable, "-d", domain, "-d", "www."+domain, "--non-interactive", "--agree-tos", "--email", email)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error configuring Let's Encrypt SSL for %s: %w", domain, err)
	} else {
		fmt.Println("\x1b[32mSuccessfully configured Let's Encrypt SSL for", domain, "\x1b[0m")
	}
	return nil
}
