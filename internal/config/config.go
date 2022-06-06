package config

import (
	"errors"
	"io/ioutil"

	"github.com/jackdallas/premiumizearr/internal/utils"
	log "github.com/sirupsen/logrus"

	"os"
	"path"

	"gopkg.in/yaml.v2"
)

// LoadOrCreateConfig - Loads the config from disk or creates a new one
func LoadOrCreateConfig(altConfigLocation string, _appCallback AppCallback) (Config, error) {
	config, err := loadConfigFromDisk(altConfigLocation)

	if err != nil {
		if err == ErrFailedToFindConfigFile {
			log.Warn("No config file found, created default config file")
			config = defaultConfig()
		}
		if err == ErrInvalidConfigFile || err == ErrFailedToSaveConfig {
			return config, err
		}
	}

	// Override unzip directory if running in docker
	if utils.IsRunningInDockerContainer() {
		log.Info("Running in docker, overriding unzip directory!")
		config.UnzipDirectory = "/unzip"
		// Override config data directories if blank
		if config.BlackholeDirectory == "" {
			log.Trace("Running in docker, overriding blank directory settings for blackhole directory")
			config.BlackholeDirectory = "/blackhole"
		}
		if config.DownloadsDirectory == "" {
			log.Trace("Running in docker, overriding blank directory settings for downloads directory")
			config.DownloadsDirectory = "/downloads"
		}
	}

	log.Tracef("Setting config location to %s", altConfigLocation)

	config.appCallback = _appCallback
	config.altConfigLocation = altConfigLocation

	config.Save()

	return config, nil
}

// Save - Saves the config to disk
func (c *Config) Save() error {
	log.Trace("Marshaling & saving config")
	data, err := yaml.Marshal(*c)
	if err != nil {
		log.Error(err)
		return err
	}

	savePath := "./config.yaml"
	if c.altConfigLocation != "" {
		savePath = path.Join(c.altConfigLocation, "config.yaml")
	}

	log.Tracef("Writing config to %s", savePath)
	err = ioutil.WriteFile(savePath, data, 0644)
	if err != nil {
		log.Errorf("Failed to save config file: %+v", err)
		return err
	}

	log.Trace("Config saved")
	return nil
}

func loadConfigFromDisk(altConfigLocation string) (Config, error) {
	var config Config

	log.Trace("Trying to load config from disk")

	configLocation := path.Join(altConfigLocation, "config.yaml")

	log.Tracef("Reading config from %s", configLocation)

	file, err := ioutil.ReadFile(configLocation)

	if err != nil {
		log.Trace("Failed to find config file")
		return config, ErrFailedToFindConfigFile
	}

	log.Trace("Attempting to unmarshal config")
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Trace("Failed to unmarshal config")
		return config, ErrInvalidConfigFile
	}

	log.Trace("Config loaded, running version update")
	config, updated := versionUpdateConfig(config)

	config.altConfigLocation = altConfigLocation

	if updated {
		log.Trace("Version updated saving")
		err = config.Save()

		if err == nil {
			log.Trace("Config saved")
			return config, nil
		} else {
			log.Errorf("Failed to save config to %s", configLocation)
			log.Error(err)
			return config, ErrFailedToSaveConfig
		}
	}

	log.Trace("Config loaded")
	return config, nil
}

func versionUpdateConfig(config Config) (Config, bool) {
	updated := false
	// 1.1.3
	if config.SimultaneousDownloads == 0 {
		config.SimultaneousDownloads = 5
		updated = true
	}

	return config, updated
}

func defaultConfig() Config {
	return Config{
		PremiumizemeAPIKey: "xxxxxxxxx",
		Arrs: []ArrConfig{
			{Name: "Sonarr", URL: "http://localhost:8989", APIKey: "xxxxxxxxx", Type: Sonarr},
			{Name: "Radarr", URL: "http://localhost:7878", APIKey: "xxxxxxxxx", Type: Radarr},
		},
		BlackholeDirectory:    "",
		DownloadsDirectory:    "",
		UnzipDirectory:        "",
		BindIP:                "0.0.0.0",
		BindPort:              "8182",
		WebRoot:               "",
		SimultaneousDownloads: 5,
	}
}

var (
	ErrUnzipDirectorySetToRoot    = errors.New("unzip directory set to root")
	ErrUnzipDirectoryNotWriteable = errors.New("unzip directory not writeable")
)

func (c *Config) GetUnzipBaseLocation() (string, error) {
	if c.UnzipDirectory == "" {
		log.Tracef("Unzip directory not set, using default: %s", os.TempDir())
		return path.Join(os.TempDir(), "premiumizearrd"), nil
	}

	if c.UnzipDirectory == "/" || c.UnzipDirectory == "\\" || c.UnzipDirectory == "C:\\" {
		log.Error("Unzip directory set to root, please set a directory")
		return "", ErrUnzipDirectorySetToRoot
	}

	if !utils.IsDirectoryWriteable(c.UnzipDirectory) {
		log.Errorf("Unzip directory not writeable: %s", c.UnzipDirectory)
		return c.UnzipDirectory, ErrUnzipDirectoryNotWriteable
	}

	log.Tracef("Unzip directory set to: %s", c.UnzipDirectory)
	return c.UnzipDirectory, nil
}

func (c *Config) GetNewUnzipLocation() (string, error) {
	// Create temp dir in os temp location or unzip-directory
	tempDir, err := c.GetUnzipBaseLocation()
	if err != nil {
		return "", err
	}

	log.Trace("Creating unzip directory")
	err = os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	log.Trace("Creating generated unzip directory")
	dir, err := ioutil.TempDir(tempDir, "unzip-")
	if err != nil {
		return "", err
	}
	return dir, nil
}
