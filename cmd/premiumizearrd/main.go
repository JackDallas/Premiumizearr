package main

import (
	"flag"
	"time"

	"github.com/jackdallas/premiumizearr/internal/arr"
	"github.com/jackdallas/premiumizearr/internal/config"
	"github.com/jackdallas/premiumizearr/internal/service"
	"github.com/jackdallas/premiumizearr/internal/web_service"
	"github.com/jackdallas/premiumizearr/pkg/premiumizeme"
	"github.com/orandin/lumberjackrus"
	log "github.com/sirupsen/logrus"
	"golift.io/starr"
	"golift.io/starr/radarr"
	"golift.io/starr/sonarr"
)

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
	hook, err := lumberjackrus.NewHook(
		&lumberjackrus.LogFile{
			Filename:   "/var/log/premiumizearr.general.log",
			MaxSize:    100,
			MaxBackups: 1,
			MaxAge:     1,
			Compress:   false,
			LocalTime:  false,
		},
		log.InfoLevel,
		&log.TextFormatter{},
		&lumberjackrus.LogFileOpts{
			log.InfoLevel: &lumberjackrus.LogFile{
				Filename:   "/var/log/premiumizearr.info.log",
				MaxSize:    100,
				MaxBackups: 1,
				MaxAge:     1,
				Compress:   false,
				LocalTime:  false,
			},
			log.ErrorLevel: &lumberjackrus.LogFile{
				Filename:   "/var/log/premiumizearr.error.log",
				MaxSize:    100,   // optional
				MaxBackups: 1,     // optional
				MaxAge:     1,     // optional
				Compress:   false, // optional
				LocalTime:  false, // optional
			},
		},
	)

	if err != nil {
		panic(err)
	}

	log.AddHook(hook)

	log.Info("---------- Starting premiumizearr daemon ----------")
	log.Info("")

	config, err := config.LoadOrCreateConfig(configFile)

	if err != nil {
		panic(err)
	}

	if config.PremiumizemeAPIKey == "" {
		panic("premiumizearr API Key is empty")
	}

	// Initialisation

	premiumizearr_client := premiumizeme.NewPremiumizemeClient(config.PremiumizemeAPIKey)

	starr_config_sonarr := starr.New(config.SonarrAPIKey, config.SonarrURL, 0)
	starr_config_radarr := starr.New(config.RadarrAPIKey, config.RadarrURL, 0)

	sonarr_wrapper := arr.SonarrArr{
		Client:     sonarr.New(starr_config_sonarr),
		History:    nil,
		LastUpdate: time.Now(),
	}
	radarr_wrapper := arr.RadarrArr{
		Client:     radarr.New(starr_config_radarr),
		History:    nil,
		LastUpdate: time.Now(),
	}

	arrs := []arr.IArr{
		&sonarr_wrapper,
		&radarr_wrapper,
	}

	transfer_manager := service.NewTransferManagerService(premiumizearr_client, &arrs, &config)

	directory_watcher := service.NewDirectoryWatcherService(premiumizearr_client, &config)

	go directory_watcher.Watch()

	go web_service.StartWebServer(&transfer_manager, &directory_watcher, &config)
	//Block until the program is terminated
	transfer_manager.Run(15 * time.Second)
}
