package main

import (
	"flag"
	"path"
	"time"

	"github.com/jackdallas/premiumizearr/internal/arr"
	"github.com/jackdallas/premiumizearr/internal/config"
	"github.com/jackdallas/premiumizearr/internal/service"
	"github.com/jackdallas/premiumizearr/internal/utils"
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
	var loggingDirectory string

	//Parse flags
	flag.StringVar(&logLevel, "log", utils.EnvOrDefault("PREMIUMIZEARR_LOG_LEVEL", "info"), "Logging level: \n \tinfo,debug,trace")
	flag.StringVar(&configFile, "config", utils.EnvOrDefault("PREMIUMIZEARR_CONFIG_DIR_PATH", "./"), "The directory the config.yml is located in")
	flag.StringVar(&loggingDirectory, "logging-dir", utils.EnvOrDefault("PREMIUMIZEARR_LOGGING_DIR_PATH", "./"), "The directory logs are to be written to")
	flag.Parse()

	lvl, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Errorf("Error flag not recognized, defaulting to Info!!", err)
		lvl = log.InfoLevel
	}
	log.SetLevel(lvl)
	hook, err := lumberjackrus.NewHook(
		&lumberjackrus.LogFile{
			Filename:   path.Join(loggingDirectory, "premiumizearr.general.log"),
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
				Filename:   path.Join(loggingDirectory, "premiumizearr.info.log"),
				MaxSize:    100,
				MaxBackups: 1,
				MaxAge:     1,
				Compress:   false,
				LocalTime:  false,
			},
			log.ErrorLevel: &lumberjackrus.LogFile{
				Filename:   path.Join(loggingDirectory, "premiumizearr.error.log"),
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
		log.Warn("Premiumizeme API key not set, application will not work until it's set")
	}

	// Initialisation

	premiumizearr_client := premiumizeme.NewPremiumizemeClient(config.PremiumizemeAPIKey)

	arrs := []arr.IArr{}

	for _, arr_config := range config.Arrs {
		switch arr_config.Type {
		case "Sonarr":
			config := starr.New(arr_config.APIKey, arr_config.URL, 0)
			wrapper := arr.SonarrArr{
				Name:       arr_config.Name,
				Client:     sonarr.New(config),
				History:    nil,
				LastUpdate: time.Now(),
			}
			arrs = append(arrs, &wrapper)
		case "Radarr":
			config := starr.New(arr_config.APIKey, arr_config.URL, 0)
			wrapper := arr.RadarrArr{
				Name:       arr_config.Name,
				Client:     radarr.New(config),
				History:    nil,
				LastUpdate: time.Now(),
			}
			arrs = append(arrs, &wrapper)
		default:
			log.Error("Unknown arr type: %s, not adding Arr %s", arr_config.Type, arr_config.Name)
		}
	}

	transfer_manager := service.NewTransferManagerService(premiumizearr_client, &arrs, &config)

	directory_watcher := service.NewDirectoryWatcherService(premiumizearr_client, &config)

	go directory_watcher.Watch()

	go web_service.StartWebServer(&transfer_manager, &directory_watcher, &config)
	//Block until the program is terminated
	transfer_manager.Run(15 * time.Second)
}
