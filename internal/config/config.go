package config

import (
	"errors"
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"os"
	"path"

	"gopkg.in/yaml.v2"
)

var (
	ErrInvalidConfigFile      = errors.New("invalid Config File")
	ErrFailedToFindConfigFile = errors.New("failed to find config file")
)

//ArrType enum for Sonarr/Radarr
type ArrType string

const (
	Sonarr ArrType = "Sonarr"
	Radarr ArrType = "Radarr"
)

type ArrConfig struct {
	Name   string  `yaml:"Name"`
	URL    string  `yaml:"URL"`
	APIKey string  `yaml:"APIKey"`
	Type   ArrType `yaml:"Type"`
}

type Config struct {
	altConfigLocation string

	PremiumizemeAPIKey string `yaml:"PremiumizemeAPIKey"`

	Arrs []ArrConfig `yaml:"Arrs"`

	BlackholeDirectory string `yaml:"BlackholeDirectory"`
	DownloadsDirectory string `yaml:"DownloadsDirectory"`

	UnzipDirectory string `yaml:"UnzipDirectory"`

	BindIP   string `yaml:"bindIP"`
	BindPort string `yaml:"bindPort"`

	WebRoot string `yaml:"WebRoot"`

	SimultaneousDownloads int `yaml:"SimultaneousDownloads"`
}

func loadConfigFromDisk(altConfigLocation string) (Config, error) {
	var config Config
	file, err := ioutil.ReadFile(path.Join(altConfigLocation, "config.yaml"))

	if err != nil {
		return config, ErrFailedToFindConfigFile
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return config, ErrInvalidConfigFile
	}

	config = versionUpdateConfig(config)

	config.Save()

	config.altConfigLocation = altConfigLocation
	return config, nil
}

func (c *Config) Save() bool {
	log.Trace("Marshaling & saving config")
	data, err := yaml.Marshal(*c)
	if err != nil {
		log.Error(err)
		return false
	}

	log.Tracef("Writing config to %s", path.Join(c.altConfigLocation, "config.yaml"))
	err = ioutil.WriteFile(path.Join(c.altConfigLocation, "config.yaml"), data, 0644)
	if err != nil {
		log.Errorf("Failed to save config file: %+v", err)
		return false
	}

	log.Trace("Config saved")
	return true
}

func versionUpdateConfig(config Config) Config {
	// 1.1.3
	if config.SimultaneousDownloads == 0 {
		config.SimultaneousDownloads = 5
	}

	return config
}

func defaultConfig(altConfigLocation string) Config {
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

func LoadOrCreateConfig(altConfigLocation string) (Config, error) {

	config, err := loadConfigFromDisk(altConfigLocation)
	if err != nil {
		if err == ErrFailedToFindConfigFile {
			config = defaultConfig(altConfigLocation)
			log.Warn("No config file found, created default config file")
			config.Save()
		}
		if err == ErrInvalidConfigFile {
			return config, ErrInvalidConfigFile
		}
	}

	return config, nil
}

func (c *Config) GetTempBaseDir() string {
	if c.UnzipDirectory != "" {
		return path.Dir(c.UnzipDirectory)
	}
	return path.Join(os.TempDir(), "premiumizearrd")
}

func (c *Config) GetTempDir() (string, error) {
	// Create temp dir in os temp location
	tempDir := c.GetTempBaseDir()

	err := os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		return "", err
	}
	dir, err := ioutil.TempDir(tempDir, "unzip-")
	if err != nil {
		return "", err
	}
	return dir, nil
}
