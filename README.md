## Introduction

This program automates the process of installing WordPress and setting up an Nginx server block on an Ubuntu server. It checks for the presence of Nginx and installs it if necessary. The script then configures a new Nginx server block for your domain and installs WordPress.

It is a rewrite of the original bash script [WPNIX](https://github.com/Kalpa-Services/wpnix). You can use either the original bash script or this Go version to install WordPress on your server but only this Go version is maintained and updated.

Prerequisites
-------------

-   An Ubuntu server (20.04 or later recommended) with a minumum of 2GB RAM.
-   Root privileges on the server.
-   Basic knowledge of terminal and command-line operations.
-   A domain name pointing to your server's IP address.
-   A MySQL database and user with all privileges.
-   A MySQL password for the user.
-   A MySQL host (usually `localhost`).

Installation
------------

run `curl -s https://packagecloud.io/install/repositories/carlHandy/wpnix/script.deb.sh?any=true | sudo bash` to install the repository and then `sudo apt-get install wpnix` to install the wpnix package.

Usage
-----

The script accepts the following arguments:

-   `-d DOMAIN`: The domain name for the WordPress site.
-   `-u DBUSER`: The database username.
-   `-p DBPASS`: The database password.
-   `-n DBNAME`: The database name.
-   `-H DBHOST`: The database host (usually `localhost`).
-   `-e EMAIL`: The email address for LetsEncrypt SSL.

Example usage:

`sudo wpnix -d example.com -u wordpressuser -p password -n wordpressdb -H localhost -e test@example.com`

If you're using a managed database service for example Digital Ocean that does not use the default `3306` port for MySQL, append your port to the DB Host. For example:

`sudo wpnix -d example.com -u wordpressuser -p password -n wordpressdb -H managedb:25062 -e test@example.com`

Features
--------

-   Checks for the presence of Nginx, Perl, PHP 8.2, and PHP 8.2-FPM, and installs them if not found.
-   Sets up an Nginx server block for the specified domain.
- Sets up letsencrypt ssl for the specified domain.
-   Installs the latest version of WordPress.
-   Configures WordPress with the provided database details.

Important Notes
---------------

-   Ensure that you have all the required database details before running the script.
-   The script must be run as `root` or with `sudo` to perform necessary system modifications.
-   It's recommended to use this script on a fresh Ubuntu installation to prevent any conflicts with existing configurations.

Troubleshooting
---------------

If you encounter any issues:

-   Check the syntax of the command and ensure all required arguments are provided.
-   Verify that your user has root privileges.
-   Ensure your server's package manager is functioning properly.

For more help, you can check the script's output for error messages.
