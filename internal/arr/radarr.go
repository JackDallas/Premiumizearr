package arr

import (
	"fmt"
	"time"

	"github.com/jackdallas/premiumizearr/pkg/premiumizeme"
	log "github.com/sirupsen/logrus"
	"golift.io/starr/radarr"
)

//////
//Radarr
//////

//Data Access

// GetHistory: Updates the history if it's been more than 15 seconds since last update
func (arr *RadarrArr) GetHistory() (radarr.History, error) {
	arr.LastUpdateMutex.Lock()
	defer arr.LastUpdateMutex.Unlock()
	arr.HistoryMutex.Lock()
	defer arr.HistoryMutex.Unlock()
	arr.ClientMutex.Lock()
	defer arr.ClientMutex.Unlock()
	arr.LastUpdateCountMutex.Lock()
	defer arr.LastUpdateCountMutex.Unlock()

	if time.Since(arr.LastUpdate) > 30*time.Second || arr.History == nil {
		his, err := arr.Client.GetHistory(0, 1000)
		if err != nil {
			return radarr.History{}, err
		}

		arr.History = his
		arr.LastUpdate = time.Now()
		arr.LastUpdateCount = his.TotalRecords
	}

	log.Tracef("Radarr.GetHistory(): Returning from GetHistory")
	return *arr.History, nil
}

func (arr *RadarrArr) MarkHistoryItemAsFailed(id int64) error {
	arr.ClientMutex.Lock()
	defer arr.ClientMutex.Unlock()
	return arr.Client.Fail(id)

}

func (arr *RadarrArr) GetArrName() string {
	return "Radarr"
}

//Functions

func (arr *RadarrArr) HistoryContains(name string) (int64, bool) {
	log.Tracef("Radarr [%s]: Checking history for %s", arr.Name, name)
	his, err := arr.GetHistory()
	if err != nil {
		log.Errorf("Radarr [%s]: Failed to get history: %+v", arr.Name, err)
		return -1, false
	}
	log.Tracef("Radarr [%s]: Got History, now Locking History", arr.Name)
	arr.HistoryMutex.Lock()
	defer arr.HistoryMutex.Unlock()

	for _, item := range his.Records {
		if CompareFileNamesFuzzy(item.SourceTitle, name) {
			return item.ID, true
		}
	}

	return -1, false
}

func (arr *RadarrArr) HandleErrorTransfer(transfer *premiumizeme.Transfer, arrID int64, pm *premiumizeme.Premiumizeme) error {
	his, err := arr.GetHistory()
	if err != nil {
		return fmt.Errorf("failed to get history from radarr: %+v", err)
	}

	arr.HistoryMutex.Lock()
	defer arr.HistoryMutex.Unlock()

	complete := false

	for _, queueItem := range his.Records {
		if queueItem.ID == arrID {
			if queueItem.EventType == "grabbed" {
				err := arr.MarkHistoryItemAsFailed(queueItem.ID)
				if err != nil {
					return fmt.Errorf("failed to blacklist item in radarr: %+v", err)
				}
				err = pm.DeleteTransfer(transfer.ID)
				if err != nil {
					return fmt.Errorf("failed to delete transfer from premiumize.me: %+v", err)
				}
				complete = true
				break
			}
		}
	}

	if !complete {
		err := pm.DeleteTransfer(transfer.ID)
		if err != nil {
			return fmt.Errorf("failed to delete transfer from premiumize.me: %+v", err)
		}
	}

	return nil
}
