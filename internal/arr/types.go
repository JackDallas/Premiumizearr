package arr

import (
	"strings"
	"sync"
	"time"

	"github.com/jackdallas/premiumizearr/internal/utils"
	"github.com/jackdallas/premiumizearr/pkg/premiumizeme"
	"golift.io/starr/radarr"
	"golift.io/starr/sonarr"
)

func CompareFileNamesFuzzy(a, b string) bool {
	//Strip file extension
	a = utils.StripDownloadTypesExtention(a)
	b = utils.StripDownloadTypesExtention(b)
	//Strip media type extension
	a = utils.StripMediaTypesExtention(a)
	b = utils.StripMediaTypesExtention(b)
	//Strip Spaces
	a = strings.ReplaceAll(a, " ", "")
	b = strings.ReplaceAll(b, " ", "")
	//Strip periods
	a = strings.ReplaceAll(a, ".", "")
	b = strings.ReplaceAll(b, ".", "")
	//Strip dashes
	a = strings.ReplaceAll(a, "-", "")
	b = strings.ReplaceAll(b, "-", "")
	//Strip underscores
	a = strings.ReplaceAll(a, "_", "")
	b = strings.ReplaceAll(b, "_", "")
	//Convert to lowercase
	a = strings.ToLower(a)
	b = strings.ToLower(b)

	return a == b
}

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
