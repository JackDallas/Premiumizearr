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

```
wget https://github.com/JackDallas/Premiumizearr/releases/download/x.x.x/Premiumizearr_x.x.x_linux_amd64.tar.gz
tar xf Premiumizearr_x.x.x.x_linux_amd64.tar.gz
cd Premiumizearr_x.x.x.x_linux_amd64
sudo mkdir /opt/premiumizearrd/
sudo cp -r premiumizearrd static/ /opt/premiumizearrd/
sudo cp premiumizearrd /etc/systemd/system/
sudo systemctl-reload
sudo systemctl enable premiumizearrd.service
sudo systemctl start premiumizearrd.service
```

### deb file

```
wget https://github.com/JackDallas/Premiumizearr/releases/download/x.x.x/premiumizearr_x.x.x._linux_amd64.deb
sudo dpkg -i premiumizearr_x.x.x.x_linux_amd64.deb
```

### Docker

[Docker images are listed here](https://github.com/jackdallas/Premiumizearr/pkgs/container/premiumizearr)

`docker run -p 8182:8182 -v /host/data/path:/data -v /host/downloads/path:/downloads -v /host/blackhole/path:/blackhole ghcr.io/jackdallas/premiumizearr:latest`

> Note: The /data mount is where the `config.yaml` and log files are kept

## Setup

### Premiumizearrd

Running for the first time the server will start on http://0.0.0.0:8182

If you already use this binding you can edit them in the `config.yaml` 

> Note: Currently most changes in the config ui will not be used until a restart is complete

### Sonarr/Radarr

- Go to your Arr's `Download Client` settings page

- Add a new Torrent Blackhole client, set the `Torrent Folder` to the previously set `BlackholeDirectory` location, set the `Watch Folder` to the previously set `DownloadsDirectory` location

- Add a new Usenet Blackhole client, set the `Nzb Folder` to the previously set `BlackholeDirectory` location, set the `Watch Folder` to the previously set `DownloadsDirectory` location
