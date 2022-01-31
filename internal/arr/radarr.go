package arr

import (
	"fmt"
	"math"
	"time"

	"github.com/jackdallas/premiumizearr/internal/utils"
	"github.com/jackdallas/premiumizearr/pkg/premiumizeme"
	log "github.com/sirupsen/logrus"
	"golift.io/starr"
	"golift.io/starr/radarr"
)

//////
//Radarr
//////

//Data Access

//GetHistory: Updates the history if it's been more than 15 seconds since last update
func (arr *RadarrArr) GetHistory() (radarr.History, error) {
	arr.LastUpdateMutex.Lock()
	defer arr.LastUpdateMutex.Unlock()
	arr.HistoryMutex.Lock()
	defer arr.HistoryMutex.Unlock()
	arr.ClientMutex.Lock()
	defer arr.ClientMutex.Unlock()
	arr.LastUpdateCountMutex.Lock()
	defer arr.LastUpdateCountMutex.Unlock()

	if time.Since(arr.LastUpdate) > 60*time.Second || arr.History == nil {
		//Get first page of records
		his, err := arr.Client.GetHistoryPage(&starr.Req{PageSize: 250, Page: 1})
		if err != nil {
			return radarr.History{}, fmt.Errorf("failed to get history from radarr: %+v", err)
		}

		if his.TotalRecords == arr.LastUpdateCount && his.TotalRecords > 0 {
			return *arr.History, nil
		}

		if his.TotalRecords > 250 {
			cachedPages := int(math.Ceil(float64(arr.LastUpdateCount) / 250))
			log.Tracef("Loaded %d cached pages of history\n", cachedPages)
			remotePages := int(math.Ceil(float64(his.TotalRecords) / float64(250)))
			log.Tracef("Found %d pages of history on the radarr server\n", cachedPages)
			for i := 2; i <= remotePages-cachedPages; i++ {
				log.Tracef("Radarr.GetHistory(): Getting History Page %d", i)
				h, err := arr.Client.GetHistoryPage(&starr.Req{PageSize: 250, Page: i})
				if err != nil {
					return radarr.History{}, fmt.Errorf("failed to get history from radarr: %+v", err)
				}
				his.Records = append(his.Records, h.Records...)
			}
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
	log.Tracef("Radarr.HistoryContains(): Checking history for %s", name)
	his, err := arr.GetHistory()
	if err != nil {
		log.Errorf("Radarr.HistoryContains(): Failed to get history: %+v", err)
		return -1, false
	}
	log.Trace("Radarr.HistoryContains(): Got History, now Locking History")
	arr.HistoryMutex.Lock()
	defer arr.HistoryMutex.Unlock()

	name = utils.StripDownloadTypesExtention(name)
	// name = strings.ReplaceAll(name, ".", " ")
	for _, item := range his.Records {
		if item.SourceTitle == name {
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
