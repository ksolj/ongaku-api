#!/bin/bash
set -eu

# ==================================================================================== #
# VARIABLES
# ==================================================================================== #

TIMEZONE=Europe/Moscow

export LC_ALL=en_US.UTF-8 

# ==================================================================================== #
# SCRIPT LOGIC
# ==================================================================================== #

add-apt-repository --yes universe

apt update
apt --yes -o Dpkg::Options::="--force-confnew" upgrade

timedatectl set-timezone ${TIMEZONE}
apt --yes install locales-all

ufw allow 22
ufw allow 80/tcp
ufw allow 443/tcp
ufw --force enable

apt --yes install fail2ban

curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz
mv migrate.linux-amd64 /usr/local/bin/migrate

apt --yes install postgresql

sudo -i -u postgres psql -c "CREATE DATABASE ongaku"
sudo -i -u postgres psql -d ongaku -c "CREATE EXTENSION IF NOT EXISTS citext"
sudo -i -u postgres psql -d ongaku -c "CREATE ROLE ongaku WITH LOGIN PASSWORD '${DB_PASSWORD}'"

echo "ONGAKU_DB_DSN='postgres://ongaku:${DB_PASSWORD}@localhost/ongaku'" >> /etc/environment

apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
apt update
apt --yes install caddy

echo "Script complete! Rebooting..."
reboot