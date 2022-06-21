# Premiumizearr

[![Build](https://github.com/JackDallas/Premiumizearr/actions/workflows/build.yml/badge.svg)](https://github.com/JackDallas/Premiumizearr/actions/workflows/build.yml)

## Features

- Monitor blackhole directory to push `.magnet` and `.nzb` to Premiumize.me
- Monitor and download Premiumize.me transfers (web ui on default port 8182)
- Mark transfers as failed in Radarr & Sonarr

## Support the project by using my invite code

[Invite Code](https://www.premiumize.me/ref/446038083)

## Install

[Grab the latest release artifact links here](https://github.com/JackDallas/Premiumizearr/releases/)

### Binary

#### System Install

```cli
wget https://github.com/JackDallas/Premiumizearr/releases/download/x.x.x/Premiumizearr_x.x.x_linux_amd64.tar.gz
tar xf Premiumizearr_x.x.x.x_linux_amd64.tar.gz
cd Premiumizearr_x.x.x.x_linux_amd64
sudo mkdir /opt/premiumizearrd/
sudo cp -r premiumizearrd static/ /opt/premiumizearrd/
sudo cp premiumizearrd.service /etc/systemd/system/
sudo systemctl-reload
sudo systemctl enable premiumizearrd.service
sudo systemctl start premiumizearrd.service
```

#### User Install

```cli
wget https://github.com/JackDallas/Premiumizearr/releases/download/x.x.x/Premiumizearr_x.x.x_linux_amd64.tar.gz
tar xf Premiumizearr_x.x.x.x_linux_amd64.tar.gz
cd Premiumizearr_x.x.x.x_linux_amd64
mkdir -p ~/.local/bin/
cp -r premiumizearrd static/ ~/.local/bin/
echo -e "export PATH=~/.local/bin/:$PATH" >> ~/.bashrc 
source ~/.bashrc
```

You're now able to run the daemon from anywhere just by typing `premiumizearrd`

### deb file

```cmd
wget https://github.com/JackDallas/Premiumizearr/releases/download/x.x.x/premiumizearr_x.x.x._linux_amd64.deb
sudo dpkg -i premiumizearr_x.x.x.x_linux_amd64.deb
```

### Docker

[Docker images are listed here](https://github.com/jackdallas/Premiumizearr/pkgs/container/premiumizearr)

```cmd
docker run \
    -v /home/dallas/test/data:/data \
    -v /home/dallas/test/blackhole:/blackhole \
    -v /home/dallas/test/downloads:/downloads \
    -p 8182:8182 \
    ghcr.io/jackdallas/premiumizearr:latest
```

If you wish to increase logging (which you'll be asked to do if you submit an issue) you can add `-e PREMIUMIZEARR_LOG_LEVEL=trace` to the command

> Note: The /data mount is where the `config.yaml` and log files are kept

## Setup

### Premiumizearrd

Running for the first time the server will start on `http://0.0.0.0:8182`

If you already use this binding for something else you can edit them in the `config.yaml`

> WARNING: This app exposes api keys in the ui and does not have authentication, it is strongly recommended you put it behind a reverse proxy with auth and set the host to `127.0.0.1` to hide the app from the web.

### Sonarr/Radarr

- Go to your Arr's `Download Client` settings page

- Add a new Torrent Blackhole client, set the `Torrent Folder` to the previously set `BlackholeDirectory` location, set the `Watch Folder` to the previously set `DownloadsDirectory` location

- Add a new Usenet Blackhole client, set the `Nzb Folder` to the previously set `BlackholeDirectory` location, set the `Watch Folder` to the previously set `DownloadsDirectory` location

### Reverse Proxy

Premiumizearr does not have authentication built in so it's strongly recommended you use a reverse proxy

#### Nginx

```nginx
location /premiumizearr/ {
    proxy_pass http://127.0.0.1:8182/;
    proxy_set_header Host $proxy_host;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Host $host;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_redirect off;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection $http_connection;
}
```
