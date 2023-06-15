package service

import (
	"time"

	"github.com/jackdallas/premiumizearr/internal/arr"
	"github.com/jackdallas/premiumizearr/internal/config"
	log "github.com/sirupsen/logrus"
	"golift.io/starr"
	"golift.io/starr/radarr"
	"golift.io/starr/sonarr"
)

type ArrsManagerService struct {
	arrs   []arr.IArr
	config *config.Config
}

func (am ArrsManagerService) New() ArrsManagerService {
	am.arrs = []arr.IArr{}
	return am
}

func (am *ArrsManagerService) Init(_config *config.Config) {
	am.config = _config
}

func (am *ArrsManagerService) Start() {
	am.arrs = []arr.IArr{}
	log.Debugf("Starting ArrsManagerService")
	for _, arr_config := range am.config.Arrs {
		switch arr_config.Type {
		case config.Sonarr:
			c := starr.New(arr_config.APIKey, arr_config.URL, 0)
			wrapper := arr.SonarrArr{
				Name:       arr_config.Name,
				Client:     sonarr.New(c),
				History:    nil,
				LastUpdate: time.Now(),
				Config:     am.config,
			}
			am.arrs = append(am.arrs, &wrapper)
			log.Tracef("Added Sonarr arr: %s", arr_config.Name)
		case config.Radarr:
			c := starr.New(arr_config.APIKey, arr_config.URL, 0)
			wrapper := arr.RadarrArr{
				Name:       arr_config.Name,
				Client:     radarr.New(c),
				History:    nil,
				LastUpdate: time.Now(),
				Config:     am.config,
			}
			am.arrs = append(am.arrs, &wrapper)
			log.Tracef("Added Radarr arr: %s", arr_config.Name)
		default:
			log.Errorf("Unknown arr type: %s, not adding Arr %s", arr_config.Type, arr_config.Name)
		}
	}
	log.Debugf("Created %d Arrs", len(am.arrs))
}

func (am *ArrsManagerService) Stop() {
	//noop
}

func (am *ArrsManagerService) ConfigUpdatedCallback(currentConfig config.Config, newConfig config.Config) {
	if len(currentConfig.Arrs) != len(newConfig.Arrs) {
		am.Start()
		return
	}
	for i, arr_config := range newConfig.Arrs {
		if currentConfig.Arrs[i].Type != arr_config.Type ||
			currentConfig.Arrs[i].APIKey != arr_config.APIKey ||
			currentConfig.Arrs[i].URL != arr_config.URL {
			am.Start()
			return
		}
	}
}

func (am *ArrsManagerService) GetArrs() []arr.IArr {
	return am.arrs
}

func TestArrConnection(arr config.ArrConfig) error {
	c := starr.New(arr.APIKey, arr.URL, 0)

	switch arr.Type {
	case config.Sonarr:
		_, err := sonarr.New(c).GetSystemStatus()
		return err
	case config.Radarr:
		_, err := radarr.New(c).GetSystemStatus()
		return err
	default:
		return nil
	}
}
