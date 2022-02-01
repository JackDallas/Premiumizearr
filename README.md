# Premiumizearr

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

## Setup

### Premiumizearrd

Edit the config file at `/opt/premiumizearrd/config.yml`

`RadarrURL,SonarrURL` Url the Arr can be accessed on

`RadarrAPIKey,SonarrAPIKey` API key for the Arr

`PremiumizemeAPIKey` API key for your [premiumize.me](https://www.premiumize.me) account

`BlackholeDirectory` Path to Directory the Arr's will put magnet/torrent/nzb files in

`DownloadsDirectory` Path for Premiumizearr to download media files to, that the Arr's watch for new media

`bindIP` IP the web server binds to

`bindPort` Port the web server binds to

### Sonarr/Radarr

- Go to your Arr's `Download Client` settings page

- Add a new Torrent Blackhole client, set the `Torrent Folder` to the previously set `BlackholeDirectory` location, set the `Watch Folder` to the previously set `DownloadsDirectory` location

- Add a new Usenet Blackhole client, set the `Nzb Folder` to the previously set `BlackholeDirectory` location, set the `Watch Folder` to the previously set `DownloadsDirectory` location
