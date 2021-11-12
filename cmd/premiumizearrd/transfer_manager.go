package main

import (
	"time"

	"github.com/jackdallas/premiumizearr/internal/utils"
	premiumizeme "github.com/jackdallas/premiumizearr/pkg/premiumizeme"
	"github.com/jackdallas/starr/sonarr"
	log "github.com/sirupsen/logrus"
)

type TransfersManager struct {
	premiumizearrd *premiumizearrd
	LastUpdated    int64
	Transfers      []premiumizeme.Transfer
	RunningTask    bool
	DownloadList   []string
}

func (manager *TransfersManager) Run(interval time.Duration) {
	for {
		manager.RunningTask = true
		manager.TaskUpdateTransfersList()
		manager.RunningTask = false
		manager.LastUpdated = time.Now().Unix()
		time.Sleep(interval)
	}
}

func (manager *TransfersManager) TaskUpdateTransfersList() {
	log.Debug("Running Task UpdateTransfersList")
	transfers, err := manager.premiumizearrd.premiumizearrClient.GetTransfers()
	if err != nil {
		log.Error(err)
		return
	}
	manager.updateTransfers(transfers)

	// TODO: pagination
	sonarrQueue, err := manager.premiumizearrd.SonarrClient.GetHistory(15000)
	if err != nil {
		log.Error(err)
		return
	}

	for _, transfer := range transfers {
		switch transfer.Status {
		case "error":
			log.Debugf("Processing transfer that has errored: ", transfer.Name)
			go manager.HandleErrorTransfer(transfer, sonarrQueue)
		case "finished":
			if utils.StringInSlice(transfer.Name, manager.DownloadList) == -1 {
				log.Debugf("Processing transfer that has finished: ", transfer.Name)
				go manager.HandleFinishedTransfer(transfer)
				manager.DownloadList = append(manager.DownloadList, transfer.Name)
			}
		case "running":
		case "queued":
			continue
		default:
			log.Tracef("Undefined Event: %s (This is not an error, just an unhandled condition)", transfer.Status)
		}
	}
}

func (manager *TransfersManager) updateTransfers(transfers []premiumizeme.Transfer) {
	manager.Transfers = transfers
}

func (manager *TransfersManager) HandleFinishedTransfer(transfer premiumizeme.Transfer) {
	log.Debug("Downloading: ", transfer.Name)
	log.Tracef("%+v", transfer)
	var link string
	var err error
	if len(transfer.FileID) > 0 {
		link, err = manager.premiumizearrd.premiumizearrClient.GenerateZippedFileLink(transfer.FileID)
	} else if len(transfer.FolderID) > 0 {
		link, err = manager.premiumizearrd.premiumizearrClient.GenerateZippedFolderLink(transfer.FolderID)
	} else {
		log.Errorf("No file or folder ID found on transfer!! Can't download %s", transfer.Name)
		return
	}
	if err != nil {
		log.Error(err)
		return
	}
	log.Trace("Downloading: ", link)
	err = utils.DownloadAndExtractZip(link, manager.premiumizearrd.Config.SonarrDownloadsDirectory)
	if err != nil {
		log.Error(err)
		return
	}
	err = manager.premiumizearrd.premiumizearrClient.DeleteTransfer(transfer.ID)
	if err != nil {
		log.Error(err)
		return
	}
}

func (manager *TransfersManager) HandleErrorTransfer(transfer premiumizeme.Transfer, sonarrQueue *sonarr.History) {
	cleanSourceTitle := utils.StripDownloadTypesExtention(transfer.Name)
	log.Debugf("Handling errored transfer %s", cleanSourceTitle)
	complete := false
	for _, queueItem := range sonarrQueue.Records {
		if queueItem.SourceTitle == cleanSourceTitle {
			if queueItem.EventType == "grabbed" {
				log.Info("Removing Sonarr history item: ", queueItem.SourceTitle)
				err := manager.premiumizearrd.SonarrClient.MarkHistoryItemAsFailed(queueItem.ID)
				if err != nil {
					log.Errorf("Failed to blacklist item in sonarr: %+v", err)
				}
				err = manager.premiumizearrd.premiumizearrClient.DeleteTransfer(transfer.ID)
				if err != nil {
					log.Errorf("Failed to delete transfer from premiumizearrze.me: %+v", err)
				}
				complete = true
				break
			} else {
				log.Debugf("Found matching item in sonarr history %s but not grabbed status, status was %s", queueItem.SourceTitle)
			}
		}
	}
	if !complete {
		log.Debugf("No matching item found in sonarr history for %s removing from transfers!", cleanSourceTitle)
		err := manager.premiumizearrd.premiumizearrClient.DeleteTransfer(transfer.ID)
		if err != nil {
			log.Errorf("Failed to delete transfer from premiumizearrze.me: %+v", err)
		}
	}
}
