package main

import (
	"path"
	"time"

	"github.com/jackdallas/premiumizearr/internal/config"
	"github.com/jackdallas/premiumizearr/internal/service"
	"github.com/jackdallas/premiumizearr/pkg/premiumizeme"
	"github.com/orandin/lumberjackrus"
	log "github.com/sirupsen/logrus"
)

type App struct {
	config             config.Config
	premiumizemeClient premiumizeme.Premiumizeme
	transferManager    service.TransferManagerService
	directoryWatcher   service.DirectoryWatcherService
	webServer          service.WebServerService
	arrsManager        service.ArrsManagerService
}

// Makes go vet error - prevents copies
func (app *App) Lock()   {}
func (app *App) UnLock() {}

func (app *App) Start(logLevel string, configFile string, loggingDirectory string) error {
	//Setup static login
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

	log.Trace("Running load or create config")
	log.Tracef("Reading config file location from flag or env: %s", configFile)
	app.config, err = config.LoadOrCreateConfig(configFile, app.ConfigUpdatedCallback)

	if err != nil {
		panic(err)
	}

	// Initialisation
	app.premiumizemeClient = premiumizeme.NewPremiumizemeClient(app.config.PremiumizemeAPIKey)

	app.transferManager = service.TransferManagerService{}.New()
	app.directoryWatcher = service.DirectoryWatcherService{}.New()
	app.webServer = service.WebServerService{}.New()
	app.arrsManager = service.ArrsManagerService{}.New()

	// Initialise Services
	app.arrsManager.Init(&app.config)
	app.directoryWatcher.Init(&app.premiumizemeClient, &app.config)

	// Must come after arrsManager
	app.transferManager.Init(&app.premiumizemeClient, &app.arrsManager, &app.config)
	// Must come after transfer, arrManager and directory
	app.webServer.Init(&app.transferManager, &app.directoryWatcher, &app.arrsManager, &app.config)

	app.arrsManager.Start()
	app.webServer.Start()
	app.directoryWatcher.Start()
	//Block until the program is terminated
	app.transferManager.Run(15 * time.Second)

	return nil
}

func (app *App) ConfigUpdatedCallback(currentConfig config.Config, newConfig config.Config) {
	app.transferManager.ConfigUpdatedCallback(currentConfig, newConfig)
	app.directoryWatcher.ConfigUpdatedCallback(currentConfig, newConfig)
	app.webServer.ConfigUpdatedCallback(currentConfig, newConfig)
	app.arrsManager.ConfigUpdatedCallback(currentConfig, newConfig)
}
