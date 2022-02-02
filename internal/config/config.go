package config

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
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
		data, err := yaml.Marshal(config)
		if err == nil {
			//Save config to disk to add missing fields
			ioutil.WriteFile("config.yaml", data, 0644)
		}
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
	//Clean up url
	if strings.HasSuffix(config.SonarrURL, ("/")) {
		config.SonarrURL = config.SonarrURL[:len(config.SonarrURL)-1]
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
