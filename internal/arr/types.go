package arr

import (
	"sync"
	"time"

	"github.com/jackdallas/premiumizearr/pkg/premiumizeme"
	"golift.io/starr/radarr"
	"golift.io/starr/sonarr"
)

//TODO: Thread Safe version
type IArr interface {
	HistoryContains(string) (int64, bool)
	MarkHistoryItemAsFailed(int64) error
	HandleErrorTransfer(*premiumizeme.Transfer, int64, *premiumizeme.Premiumizeme) error
	GetArrName() string
}

type SonarrArr struct {
	Name                 string
	ClientMutex          sync.Mutex
	Client               *sonarr.Sonarr
	HistoryMutex         sync.Mutex
	History              *sonarr.History
	LastUpdateMutex      sync.Mutex
	LastUpdate           time.Time
	LastUpdateCount      int
	LastUpdateCountMutex sync.Mutex
}

type RadarrArr struct {
	Name                 string
	ClientMutex          sync.Mutex
	Client               *radarr.Radarr
	HistoryMutex         sync.Mutex
	History              *radarr.History
	LastUpdateMutex      sync.Mutex
	LastUpdate           time.Time
	LastUpdateCount      int
	LastUpdateCountMutex sync.Mutex
}
