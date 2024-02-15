package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	nginxAvailable = "/etc/nginx/sites-available"
	nginxEnabled   = "/etc/nginx/sites-enabled"
	webDir         = "/var/www"
	webUser        = "www-data:www-data"
)

func showHelp() {
	fmt.Println(`Usage: wpnix [-h] [-d DOMAIN] [-u DBUSER] [-p DBPASS] [-n DBNAME] [-H DBHOST]

This program installs WordPress and sets up an Nginx server block.

    -h          display this help and exit
    -d DOMAIN   specify the domain name
    -u DBUSER   database user
    -p DBPASS   database password
    -n DBNAME   database name
    -H DBHOST   database host`)
}

func finalizeSetupAndRestartNginx(domain string) {
	webPath := filepath.Join(webDir, domain)
	if err := exec.Command("chown", "-R", webUser, webPath).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "\x1b[31mError setting permissions:", err, "\x1b[0m")
		return
	}
	if err := exec.Command("chmod", "-R", "775", webPath).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "\x1b[31mError setting permissions:", err, "\x1b[0m")
		return
	}
	if err := exec.Command("ln", "-s", filepath.Join(nginxAvailable, domain), filepath.Join(nginxEnabled, domain)).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "\x1b[31mError creating symlink:", err, "\x1b[0m")
		return
	}

	if err := exec.Command("systemctl", "restart", "nginx").Run(); err != nil {
		fmt.Fprintln(os.Stderr, "\x1b[31mError restarting Nginx:", err, "\x1b[0m")
		return
	}

	fmt.Println("\x1b[32mSuccessfully finalized setup and restarted Nginx for", domain, "\x1b[0m")
}

func main() {
	var (
		domain, dbUser, dbPass, dbName, dbHost string
		help                                   bool
	)

	flag.StringVar(&domain, "d", "", "Domain name")
	flag.StringVar(&dbUser, "u", "", "Database user")
	flag.StringVar(&dbPass, "p", "", "Database password")
	flag.StringVar(&dbName, "n", "", "Database name")
	flag.StringVar(&dbHost, "H", "", "Database host")
	flag.BoolVar(&help, "h", false, "Show help")

	flag.Parse()

	if help {
		showHelp()
		return
	}

	if domain == "" || dbUser == "" || dbPass == "" || dbName == "" || dbHost == "" {
		fmt.Println("Error: All parameters are required.")
		showHelp()
		return
	}

	if os.Geteuid() != 0 {
		fmt.Println("This program must be run as root.")
		return
	}

	checkAndInstallNginx()
	checkAndInstallPerl()
	checkAndInstallPHP()
	createNginxConfig(domain)
	installWordPress(domain, dbUser, dbPass, dbName, dbHost)
	finalizeSetupAndRestartNginx(domain)
}
