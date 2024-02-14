package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func installWordPress(domain, dbUser, dbPass, dbName, dbHost string) {
	webPath := filepath.Join(webDir, domain)
	if _, err := os.Stat(webPath); os.IsNotExist(err) {
		fmt.Println("Creating web directory for", domain, "...")
		if err := os.MkdirAll(webPath, os.ModePerm); err != nil {
			fmt.Fprintln(os.Stderr, "\x1b[31mError creating web directory:", err, "\x1b[0m")
			return
		}
	}
	fmt.Println("Downloading WordPress...")
	cmd := exec.Command("curl", "-O", "https://wordpress.org/latest.tar.gz")
	cmd.Dir = webPath
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "\x1b[31mError downloading WordPress:", err, "\x1b[0m")
		return
	}
	cmd = exec.Command("tar", "-zxvf", "latest.tar.gz")
	cmd.Dir = webPath
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "\x1b[31mError extracting WordPress:", err, "\x1b[0m")
		return
	}
	if err := os.Remove(filepath.Join(webPath, "latest.tar.gz")); err != nil {
		fmt.Fprintln(os.Stderr, "\x1b[31mError cleaning up zip file:", err, "\x1b[0m")
		return
	}
	wpConfigPath := filepath.Join(webPath, "wordpress", "wp-config-sample.php")
	input, err := os.ReadFile(wpConfigPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "\x1b[31mError reading wp-config-sample.php:", err, "\x1b[0m")
		return
	}
	output := bytes.Replace(input, []byte("database_name_here"), []byte(dbName), -1)
	output = bytes.Replace(output, []byte("username_here"), []byte(dbUser), -1)
	output = bytes.Replace(output, []byte("password_here"), []byte(dbPass), -1)
	output = bytes.Replace(output, []byte("localhost"), []byte(dbHost), -1)
	if err = os.WriteFile(filepath.Join(webPath, "wordpress", "wp-config.php"), output, 0666); err != nil {
		fmt.Fprintln(os.Stderr, "\x1b[31mError writing wp-config.php:", err, "\x1b[0m")
		return
	}
	fmt.Println("\x1b[32mWordPress installed and configured successfully.\x1b[0m")
}
