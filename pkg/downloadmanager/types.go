package downloadmanager

import (
	"sync"
	"sync/atomic"
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

	IdCounter atomic.Int64

	CancelChannel chan bool
}
