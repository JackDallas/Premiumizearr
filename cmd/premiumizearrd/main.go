package main

import (
	"flag"
	"io"
	"os"
	"time"

	"github.com/jackdallas/premiumizearr/internal/config"
	"github.com/jackdallas/premiumizearr/pkg/premiumizeme"
	"github.com/jackdallas/starr"
	"github.com/jackdallas/starr/sonarr"
	log "github.com/sirupsen/logrus"
)

type premiumizearrd struct {
	Config              *config.Config
	premiumizearrClient *premiumizeme.Premiumizeme
	SonarrClient        *sonarr.Sonarr
	DirectoryWatcher    *DirectoryWatcherService
	TransferManager     *TransfersManager
}

func main() {
	//Flags
	var logLevel string
	var configFile string

	//Parse flags
	flag.StringVar(&logLevel, "log", "info", "Logging level: \n \tinfo,debug,trace")
	flag.StringVar(&configFile, "config", "", "Config file path")
	flag.Parse()

	lvl, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Errorf("Error flag not recognized, defaulting to Info!!", err)
		lvl = log.InfoLevel
	}
	log.SetLevel(lvl)

	logFile, err := os.Create("premiumizearr.log")
	if err != nil {
		log.Error(err)
	} else {
		log.SetOutput(io.MultiWriter(logFile, os.Stdout))
	}

	log.Info("Starting premiumizearr daemon")
	config, err := config.LoadOrCreateConfig(configFile)

	if err != nil {
		panic(err)
	}

	if config.PremiumizemeAPIKey == "" {
		panic("premiumizearr API Key is empty")
	}

	premiumizearrze_client := premiumizeme.NewPremiumizemeClient(config.PremiumizemeAPIKey)

	starr_config := starr.New(config.SonarrAPIKey, config.SonarrURL, 0)
	sonarr_client := sonarr.New(starr_config)

	var premarrd premiumizearrd

	transfer_manager := TransfersManager{
		premiumizearrd: &premarrd,
		LastUpdated:    time.Now().Unix(),
	}

	directory_watcher := DirectoryWatcherService{
		premiumizearrd: &premarrd,
	}

	premarrd = premiumizearrd{
		Config:              &config,
		premiumizearrClient: premiumizearrze_client,
		SonarrClient:        sonarr_client,
		DirectoryWatcher:    &directory_watcher,
		TransferManager:     &transfer_manager,
	}

	go directory_watcher.Watch()

	go StartWebServer(&premarrd)
	//Block until the program is terminated
	transfer_manager.Run(1 * time.Minute)
}
