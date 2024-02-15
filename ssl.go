package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func configureLetsEncryptSSL(domain string, email string) {
	fmt.Println("Configuring Let's Encrypt SSL for", domain, "...")
	var cmd *exec.Cmd
	if strings.Count(domain, ".") > 1 {
		cmd = exec.Command("certbot", "--nginx", "-d", domain, "--non-interactive", "--agree-tos", "--email", email)
	} else {
		cmd = exec.Command("certbot", "--nginx", "-d", domain, "-d", "www."+domain, "--non-interactive", "--agree-tos", "--email", email)
	}

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "\x1b[31mError configuring Let's Encrypt SSL for "+domain+":", err, "\x1b[0m")
	} else {
		fmt.Println("\x1b[32mSuccessfully configured Let's Encrypt SSL for", domain, "\x1b[0m")
	}
}
