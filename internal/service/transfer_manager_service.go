package service

import (
	"fmt"
	"time"

	"github.com/jackdallas/premiumizearr/internal/config"
	"github.com/jackdallas/premiumizearr/internal/utils"
	"github.com/jackdallas/premiumizearr/pkg/premiumizeme"
	log "github.com/sirupsen/logrus"
)

type TransferManagerService struct {
	premiumizemeClient *premiumizeme.Premiumizeme
	arrsManager        *ArrsManagerService
	config             *config.Config
	lastUpdated        int64
	transfers          []premiumizeme.Transfer
	runningTask        bool
	status             string
}

// Handle
func (t TransferManagerService) New() TransferManagerService {
	t.premiumizemeClient = nil
	t.arrsManager = nil
	t.config = nil
	t.lastUpdated = time.Now().Unix()
	t.transfers = make([]premiumizeme.Transfer, 0)
	t.runningTask = false
	t.status = ""
	return t
}

func (t *TransferManagerService) Init(pme *premiumizeme.Premiumizeme, arrsManager *ArrsManagerService, config *config.Config) {
	t.premiumizemeClient = pme
	t.arrsManager = arrsManager
	t.config = config
	t.CleanUpUnzipDir()
}

func (t *TransferManagerService) CleanUpUnzipDir() {
	log.Info("Cleaning unzip directory")

	unzipBase, err := t.config.GetUnzipBaseLocation()
	if err != nil {
		log.Errorf("Error getting unzip base location: %s", err.Error())
		return
	}

	err = utils.RemoveContents(unzipBase)
	if err != nil {
		log.Errorf("Error cleaning unzip directory: %s", err.Error())
		return
	}

}

func (manager *TransferManagerService) ConfigUpdatedCallback(currentConfig config.Config, newConfig config.Config) {
	//NOOP
}

func (manager *TransferManagerService) Run(interval time.Duration) {
	for {
		manager.runningTask = true
		manager.TaskUpdateTransfersList()
		manager.runningTask = false
		manager.lastUpdated = time.Now().Unix()
		time.Sleep(interval)
	}
}

func (manager *TransferManagerService) GetTransfers() *[]premiumizeme.Transfer {
	return &manager.transfers
}
func (manager *TransferManagerService) GetStatus() string {
	return manager.status
}

func (manager *TransferManagerService) updateTransfers(transfers []premiumizeme.Transfer) {
	manager.transfers = transfers
}

func (manager *TransferManagerService) TaskUpdateTransfersList() {
	log.Debug("Running Task UpdateTransfersList")
	transfers, err := manager.premiumizemeClient.GetTransfers()
	if err != nil {
		log.Errorf("Error getting transfers: %s", err.Error())
		return
	}
	manager.updateTransfers(transfers)

	log.Tracef("Checking %d transfers against %d Arr clients", len(transfers), len(manager.arrsManager.GetArrs()))
	earlyReturn := false

	if len(transfers) == 0 {
		manager.status = "No transfers"
		earlyReturn = true
	} else {
		manager.status = fmt.Sprintf("Got %d transfers", len(transfers))
	}

	if len(manager.arrsManager.GetArrs()) == 0 {
		manager.status = fmt.Sprintf("%s, no ARRs available", manager.status)
		earlyReturn = true
	}
	//else {
	// 	//TODO: Test
	// 	// if manager.status[len(manager.status)-19:] == ", no ARRs available" {
	// 	// 	manager.status = manager.status[:len(manager.status)-19]
	// 	// }
	// 	fmt.Print(manager.status)
	// }

	if earlyReturn {
		return
	}

	for _, transfer := range transfers {
		found := false
		for _, arr := range manager.arrsManager.GetArrs() {
			if found {
				break
			}
			if transfer.Status == "error" {
				log.Tracef("Checking errored transfer %s against %s history", transfer.Name, arr.GetArrName())
				arrID, contains := arr.HistoryContains(transfer.Name)
				if !contains {
					log.Tracef("%s history doesn't contain %s", arr.GetArrName(), transfer.Name)
					continue
				}
				log.Tracef("Found %s in %s history", transfer.Name, arr.GetArrName())
				found = true
				log.Debugf("Processing transfer that has errored: %s", transfer.Name)
				go arr.HandleErrorTransfer(&transfer, arrID, manager.premiumizemeClient)

			}
		}
	}
}
