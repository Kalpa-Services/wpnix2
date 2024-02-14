package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Default Configurations
var (
	nginxAvailable = "/etc/nginx/sites-available"
	nginxEnabled   = "/etc/nginx/sites-enabled"
	webDir         = "/var/www"
	webUser        = "www-data:www-data"
)

// Function to display help
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

// Function to check and install Nginx
func checkAndInstallNginx() {
	if _, err := exec.LookPath("nginx"); err != nil {
		fmt.Println("Nginx is not installed. Installing Nginx...")
		exec.Command("apt", "update").Run()
		exec.Command("apt", "install", "nginx", "-y").Run()
	} else {
		fmt.Println("Nginx is already installed.")
	}
}

// Function to check and install Perl
func checkAndInstallPerl() {
	if _, err := exec.LookPath("perl"); err != nil {
		fmt.Println("Perl is not installed. Installing Perl...")
		exec.Command("apt", "install", "perl", "-y").Run()
	} else {
		fmt.Println("Perl is already installed.")
	}
}

// Function to check and install PHP 8.2 and PHP 8.2-FPM
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

// Function to create nginx config
func createNginxConfig(domain string) {
	config := fmt.Sprintf(`server {
    listen 80;
    listen [::]:80;
    server_tokens off;
    server_name %s www.%s;
    root %s/%s/wordpress;
    index index.php;

    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;
    add_header Content-Security-Policy "default-src 'self' http: https: data: blob: 'unsafe-inline'" always;

    location = /favicon.ico {
        log_not_found off;
        access_log off;
    }

    location = /robots.txt {
        allow all;
        log_not_found off;
        access_log off;
    }

    location / {
        try_files $uri $uri/ /index.php?$args;
    }

    location ~ \.php$ {
        include snippets/fastcgi-php.conf;
        fastcgi_intercept_errors on;
        fastcgi_pass unix:/var/run/php/php8.2-fpm.sock;
    }

    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|otf)$ {
        expires 365d;
        access_log off;
        add_header Cache-Control "public";
    }

    location ~ /\.ht {
        deny all;
    }
}`, domain, domain, webDir, domain)

	file, err := os.Create(filepath.Join(nginxAvailable, domain))
	if err != nil {
		fmt.Println("\x1b[31mError creating Nginx config file:", err, "\x1b[0m")
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(config)
	writer.Flush()
}

// Function to install WordPress
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

// Function to check and install Certbot
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

// Function to configure Let's Encrypt SSL for the domain
func configureLetsEncryptSSL(domain string) {
	fmt.Println("Configuring Let's Encrypt SSL for", domain, "...")
	cmd := exec.Command("certbot", "--nginx", "-d", domain, "-d", "www."+domain)
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "\x1b[31mError configuring Let's Encrypt SSL for "+domain+":", err, "\x1b[0m")
	} else {
		fmt.Println("\x1b[32mSuccessfully configured Let's Encrypt SSL for", domain, "\x1b[0m")
	}
}

// Function to finalize setup and restart Nginx
func finalizeSetupAndRestartNginx(domain string) {
	// Set permissions and create symlink
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

	// Restart Nginx to apply changes
	if err := exec.Command("systemctl", "restart", "nginx").Run(); err != nil {
		fmt.Fprintln(os.Stderr, "\x1b[31mError restarting Nginx:", err, "\x1b[0m")
		return
	}

	fmt.Println("\x1b[32mSuccessfully finalized setup and restarted Nginx for", domain, "\x1b[0m")
}

// Main function
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
	checkAndInstallCertbot()
	configureLetsEncryptSSL(domain)
	finalizeSetupAndRestartNginx(domain)
}
