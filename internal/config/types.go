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
	Name   string  `yaml:"Name" json:"Name"`
	URL    string  `yaml:"URL" json:"URL"`
	APIKey string  `yaml:"APIKey" json:"APIKey"`
	Type   ArrType `yaml:"Type" json:"Type"`
}

type Config struct {
	altConfigLocation string
	appCallback       AppCallback

	//PremiumizemeAPIKey string with yaml and json tag
	PremiumizemeAPIKey string `yaml:"PremiumizemeAPIKey" json:"PremiumizemeAPIKey"`

	Arrs []ArrConfig `yaml:"Arrs" json:"Arrs"`

	BlackholeDirectory           string `yaml:"BlackholeDirectory" json:"BlackholeDirectory"`
	PollBlackholeDirectory       bool   `yaml:"PollBlackholeDirectory" json:"PollBlackholeDirectory"`
	PollBlackholeIntervalMinutes int    `yaml:"PollBlackholeIntervalMinutes" json:"PollBlackholeIntervalMinutes"`

	DownloadsDirectory string `yaml:"DownloadsDirectory" json:"DownloadsDirectory"`

	UnzipDirectory string `yaml:"UnzipDirectory" json:"UnzipDirectory"`

	BindIP   string `yaml:"bindIP" json:"BindIP"`
	BindPort string `yaml:"bindPort" json:"BindPort"`

	WebRoot string `yaml:"WebRoot" json:"WebRoot"`

	SimultaneousDownloads int `yaml:"SimultaneousDownloads" json:"SimultaneousDownloads"`
}
