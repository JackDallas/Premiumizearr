package config

import "errors"

var (
	ErrInvalidConfigFile      = errors.New("invalid Config File")
	ErrFailedToFindConfigFile = errors.New("failed to find config file")
	ErrFailedToSaveConfig     = errors.New("failed to save config")
)

//ArrType enum for Sonarr/Radarr
type ArrType string

//AppCallback - Callback for the app to use
type AppCallback func(oldConfig Config, newConfig Config)

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
	appCallback       AppCallback

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
