package downloadmanager

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func (d *DownloadManager) Run() {
	for {
		select {
		case <-d.CancelChannel:
			return
		default:
			time.Sleep(time.Millisecond * 100)
			for i := 0; i < len(d.transfers); i++ {
				t := &d.transfers[i]
				switch t.GetStatus() {
				case STATUS_QUEUED:
					if d.GetActiveTransferCount() < d.MaxSimultaneousDownloads {
						if err := t.Download(); err != nil {
							log.Errorf("Error downloading: %s", err)
						}
					} else {
						log.Debugf("Too many active transfers, skipping %d", t.GetID())
					}
					return
				}
			}
		}
	}
}

func (d *DownloadManager) GetTransfers() []Transfer {
	d.transfersLock.Lock()
	defer d.transfersLock.Unlock()
	return d.transfers
}

func (d *DownloadManager) GetTransfer(id int64) (*Transfer, error) {
	d.transfersLock.Lock()
	defer d.transfersLock.Unlock()
	for i := 0; i < len(d.transfers); i++ {
		if d.transfers[i].GetID() == id {
			return &d.transfers[i], nil
		}
	}
	return nil, ErrorNoTransferWithID
}

func (d *DownloadManager) AddTransfer(url string, savePath string) (*Transfer, error) {
	d.transfersLock.Lock()
	defer d.transfersLock.Unlock()

	nextID := d.IdCounter.Add(1)

	d.transfers = append(d.transfers, NewTransfer(nextID, url, savePath))

	log.Debugf("Added transfer %d", nextID)
	return d.GetTransfer(nextID)
}

func (d *DownloadManager) GetActiveTransferCount() int {
	c := 0

	for i := 0; i < len(d.transfers); i++ {
		if d.transfers[i].GetStatus() == STATUS_DOWNLOADING {
			c++
		}
	}

	return c
}

func (d *DownloadManager) RemoveTransfer(id int64) error {
	d.transfersLock.Lock()
	defer d.transfersLock.Unlock()

	for i := range d.transfers {
		if d.transfers[i].GetID() == id {
			return d.transfers[i].Cancel()
		}
	}

	return ErrorNoTransferWithID
}
