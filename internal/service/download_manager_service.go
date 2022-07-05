package service

import (
	"os"
	"path"
	"strings"
	"time"

	"github.com/jackdallas/premiumizearr/internal/config"
	"github.com/jackdallas/premiumizearr/internal/utils"

	"github.com/jackdallas/premiumizearr/pkg/downloadmanager"
	"github.com/jackdallas/premiumizearr/pkg/premiumizeme"
	log "github.com/sirupsen/logrus"
)

type DownloadManagerService struct {
	downloadManager    *downloadmanager.DownloadManager
	premiumizemeClient *premiumizeme.Premiumizeme
	config             *config.Config

	downloadsFolderID string
}

func (DownloadManagerService) New() DownloadManagerService {
	return DownloadManagerService{
		downloadsFolderID: "",
		downloadManager:   &downloadmanager.DownloadManager{},
	}
}

func (manager *DownloadManagerService) Init(_premiumizemeClient *premiumizeme.Premiumizeme, _config *config.Config) {
	manager.premiumizemeClient = _premiumizemeClient
	manager.config = _config

	manager.downloadsFolderID = utils.GetDownloadsFolderIDFromPremiumizeme(manager.premiumizemeClient)
	manager.CleanUpUnzipDir()

	go manager.downloadManager.Run()
	go manager.TaskCheckPremiumizeDownloadsFolder()
}

func (manager *DownloadManagerService) CleanUpUnzipDir() {
	log.Info("Cleaning unzip directory")

	unzipBase, err := manager.config.GetUnzipBaseLocation()
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

func (manager *DownloadManagerService) ConfigUpdatedCallback(currentConfig config.Config, newConfig config.Config) {
	if currentConfig.UnzipDirectory != newConfig.UnzipDirectory {
		manager.CleanUpUnzipDir()
	}
}

func (manager *DownloadManagerService) TaskCheckPremiumizeDownloadsFolder() {
	for {
		time.Sleep(time.Second * 60)
		log.Debug("Running Task CheckPremiumizeDownloadsFolder")

		items, err := manager.premiumizemeClient.ListFolder(manager.downloadsFolderID)
		if err != nil {
			log.Errorf("Error listing downloads folder: %s", err.Error())
			return
		}

		for _, item := range items {
			manager.downloadFinishedTransfer(item, manager.config.DownloadsDirectory)
		}
	}
}

func (manager *TransferManagerService) updateTransfers(transfers []premiumizeme.Transfer) {
	manager.transfers = transfers
}

func (manager *DownloadManagerService) downloadFinishedTransfer(item premiumizeme.Item, downloadDirectory string) {
	log.Debug("Downloading: ", item.Name)
	log.Tracef("%+v", item)
	var link string
	var err error
	if item.Type == "file" {
		link, err = manager.premiumizemeClient.GenerateZippedFileLink(item.ID)
	} else if item.Type == "folder" {
		link, err = manager.premiumizemeClient.GenerateZippedFolderLink(item.ID)
	} else {
		log.Errorf("Item is not of type 'file' or 'folder' !! Can't download %s", item.Name)
		return
	}
	if err != nil {
		log.Error("Error generating download link: %s", err)
		return
	}
	log.Trace("Downloading from: ", link)

	tempDir, err := manager.config.GetNewUnzipLocation()
	if err != nil {
		log.Errorf("Could not create temp dir: %s", err)
		return
	}

	splitString := strings.Split(link, "/")
	savePath := path.Join(tempDir, splitString[len(splitString)-1])
	log.Trace("Downloading to: ", savePath)

	out, err := os.Create(savePath)
	if err != nil {
		log.Errorf("Could not create save path: %s", err)
		return
	}
	defer out.Close()

	transfer, err := manager.downloadManager.AddTransfer(link, savePath)
	if err != nil {
		log.Errorf("Could not add transfer: %s", err)
		return
	}

	go func() {
		<-transfer.Finished

		if transfer.GetStatus() == downloadmanager.STATUS_ERROR || transfer.GetStatus() == downloadmanager.STATUS_CANCELED {
			log.Errorf("Could not download file: %s", strings.Join(transfer.GetErrorStrings(), ", "))
			return
		}

		unzipped := true
		log.Tracef("Unzipping %s to %s", savePath, downloadDirectory)
		err = utils.Unzip(savePath, downloadDirectory)
		if err != nil {
			log.Errorf("Could not unzip file: %s", err)
			unzipped = false
		}

		log.Tracef("Removing zip %s from system", savePath)
		err = os.RemoveAll(savePath)
		if err != nil {
			log.Errorf("Could not remove zip: %s", err)
			return
		}

		if unzipped {
			err = manager.premiumizemeClient.DeleteFolder(item.ID)
			if err != nil {
				log.Error("Error deleting folder on premiumize.me: %s", err)
				return
			}
		}

	}()
}
