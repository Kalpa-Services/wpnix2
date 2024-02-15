package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

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

	if err := exec.Command("systemctl", "stop", "apache2").Run(); err != nil {
		fmt.Fprintln(os.Stderr, "\x1b[31mError stopping Apache2:", err, "\x1b[0m")
		return
	}

	if err := exec.Command("systemctl", "disable", "apache2").Run(); err != nil {
		fmt.Fprintln(os.Stderr, "\x1b[31mError disabling Apache2:", err, "\x1b[0m")
		return
	}

	if err := exec.Command("systemctl", "restart", "nginx").Run(); err != nil {
		fmt.Fprintln(os.Stderr, "\x1b[31mError restarting Nginx:", err, "\x1b[0m")
		return
	}

	fmt.Println("\x1b[32mSuccessfully finalized setup and restarted Nginx for", domain, "\x1b[0m")
}
