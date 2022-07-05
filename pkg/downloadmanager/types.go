package downloadmanager

import (
	"sync"
)

// type DownloadManager interface {
// 	GetTransfers() []Transfer

// 	GetTransfer(id int64) (*Transfer, error)
// 	AddTransfer(url string) (*Transfer, error)
// 	RemoveTransfer(id int64) error
// }

type DownloadManager struct {
	MaxSimultaneousDownloads int

	transfers     []Transfer
	transfersLock sync.Mutex

	idCounter int64

	CancelChannel chan bool
}
