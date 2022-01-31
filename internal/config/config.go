package config

import (
	"errors"
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	ErrInvalidConfigFile      = errors.New("invalid Config File")
	ErrFailedToFindConfigFile = errors.New("failed to find config file")
)

type Config struct {
	PremiumizemeAPIKey string `yaml:"PremiumizemeAPIKey"`

	SonarrURL    string `yaml:"SonarrURL"`
	SonarrAPIKey string `yaml:"SonarrAPIKey"`

	RadarrURL    string `yaml:"RadarrURL"`
	RadarrAPIKey string `yaml:"RadarrAPIKey"`

	BlackholeDirectory string `yaml:"BlackholeDirectory"`
	DownloadsDirectory string `yaml:"DownloadsDirectory"`

	BindIP   string `yaml:"bindIP"`
	BindPort string `yaml:"bindPort"`
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

	return config, nil
}

func createDefaultConfig() error {
	config := Config{
		PremiumizemeAPIKey: "",
		SonarrURL:          "http://localhost:8989",
		SonarrAPIKey:       "",
		RadarrURL:          "http://localhost:7878",
		RadarrAPIKey:       "",
		BlackholeDirectory: "",
		DownloadsDirectory: "",
		BindIP:             "0.0.0.0",
		BindPort:           "8182",
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
	//Clean up url
	if strings.HasSuffix(config.SonarrURL, ("/")) {
		config.SonarrURL = config.SonarrURL[:len(config.SonarrURL)-1]
	}

	return config, nil
}
