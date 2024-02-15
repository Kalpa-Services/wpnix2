package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func checkAndInstallNginx() {
	if _, err := exec.LookPath("nginx"); err != nil {
		fmt.Println("Nginx is not installed. Installing Nginx...")
		if err := exec.Command("apt", "update").Run(); err != nil {
			fmt.Fprintln(os.Stderr, "\x1b[31mError updating package lists:", err, "\x1b[0m")
			return
		}
		if err := exec.Command("apt", "install", "nginx", "-y").Run(); err != nil {
			fmt.Fprintln(os.Stderr, "\x1b[31mError installing Nginx:", err, "\x1b[0m")
			return
		}
		fmt.Println("\x1b[32mNginx installed successfully.\x1b[0m")
	} else {
		fmt.Println("\x1b[33mNginx is already installed.\x1b[0m")
	}
}

func createNginxConfig(domain string) {
	config := fmt.Sprintf(`server {
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

	createSymlinkIfNotExists(filepath.Join(nginxAvailable, domain), filepath.Join(nginxEnabled, domain))
	stopAndDisableApache2()
	validateAndReloadNginx()

	fmt.Println("\x1b[32mSuccessfully finalized setup and restarted Nginx for", domain, "\x1b[0m")
}

func stopAndDisableApache2() {
	_, err := exec.LookPath("apache2")
	if err != nil {
		fmt.Println("Apache2 is not installed, skipping stop and disable steps.")
		return
	}
	if err := exec.Command("systemctl", "stop", "apache2").Run(); err != nil {
		fmt.Fprintln(os.Stderr, "\x1b[31mError stopping Apache2:", err, "\x1b[0m")
		return
	}

	if err := exec.Command("systemctl", "disable", "apache2").Run(); err != nil {
		fmt.Fprintln(os.Stderr, "\x1b[31mError disabling Apache2:", err, "\x1b[0m")
		return
	}
}

func validateAndReloadNginx() error {
	if err := exec.Command("nginx", "-t").Run(); err != nil {
		return fmt.Errorf("nginx configuration test failed: %w", err)
	}
	if err := exec.Command("systemctl", "reload", "nginx").Run(); err != nil {
		return fmt.Errorf("failed to reload Nginx: %w", err)
	}
	fmt.Println("Nginx configuration reloaded successfully.")
	return nil
}

func createSymlinkIfNotExists(source, target string) error {
	if _, err := os.Lstat(target); err == nil {
		fmt.Println("Symlink already exists:", target)
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check symlink: %w", err)
	}
	if err := os.Symlink(source, target); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}
	fmt.Println("Symlink created:", target)
	return nil
}
