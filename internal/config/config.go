package config

import (
	"errors"
	"io/ioutil"
	"log"
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
	PremiumizemeAPIKey string `yaml:"PremiumizemeAPIKey"`

	Arrs []ArrConfig `yaml:"Arrs"`

	BlackholeDirectory string `yaml:"BlackholeDirectory"`
	DownloadsDirectory string `yaml:"DownloadsDirectory"`

	UnzipDirectory string `yaml:"UnzipDirectory"`

	BindIP   string `yaml:"bindIP"`
	BindPort string `yaml:"bindPort"`

	WebRoot string `yaml:"WebRoot"`
}

func loadConfigFromDisk() (Config, error) {
	var config Config
	file, err := ioutil.ReadFile("config.yaml")

	if err != nil {
		return config, ErrFailedToFindConfigFile
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return config, ErrInvalidConfigFile
	}

	data, err := yaml.Marshal(config)
	if err == nil {
		//Save config to disk to add missing fields
		ioutil.WriteFile("config.yaml", data, 0644)
	}

	return config, nil
}

func createDefaultConfig() error {
	config := Config{
		PremiumizemeAPIKey: "xxxxxxxxx",
		Arrs: []ArrConfig{
			{URL: "http://localhost:8989", APIKey: "xxxxxxxxx", Type: Sonarr},
			{URL: "http://localhost:7878", APIKey: "xxxxxxxxx", Type: Radarr},
		},
		BlackholeDirectory: "",
		DownloadsDirectory: "",
		UnzipDirectory:     "",
		BindIP:             "0.0.0.0",
		BindPort:           "8182",
		WebRoot:            "",
	}

	file, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("config.yaml", file, 0644)
	if err != nil {
		return err
	}

	return nil
}

func LoadOrCreateConfig(altConfigLocation string) (Config, error) {
	if altConfigLocation != "" {
		if _, err := ioutil.ReadFile(altConfigLocation); err != nil {
			log.Panicf("Failed to find config file at %s Error: %+v", altConfigLocation, err)
		}
	}

	config, err := loadConfigFromDisk()
	if err != nil {
		if err == ErrFailedToFindConfigFile {
			err = createDefaultConfig()
			if err != nil {
				return config, err
			}
			panic("Default config created, please fill it out")
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
