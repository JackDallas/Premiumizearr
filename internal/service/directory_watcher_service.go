package service

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/jackdallas/premiumizearr/internal/config"
	"github.com/jackdallas/premiumizearr/internal/directory_watcher"
	"github.com/jackdallas/premiumizearr/internal/utils"
	"github.com/jackdallas/premiumizearr/pkg/premiumizeme"
	"github.com/jackdallas/premiumizearr/pkg/stringqueue"
	log "github.com/sirupsen/logrus"
)

type DirectoryWatcherService struct {
	premiumizemeClient *premiumizeme.Premiumizeme
	config             *config.Config
	Queue              *stringqueue.StringQueue
	status             string
	downloadsFolderID  string
}

const (
	ERROR_LIMIT_REACHED    = "Limit of transfers reached!"
	ERROR_ALREADY_UPLOADED = "You already added this job."
)

func NewDirectoryWatcherService(pm *premiumizeme.Premiumizeme, con *config.Config) DirectoryWatcherService {
	return DirectoryWatcherService{
		premiumizemeClient: pm,
		config:             con,
		status:             "",
	}
}

func (dw *DirectoryWatcherService) GetStatus() string {
	return dw.status
}

//TODO (Radarr): accept paths as a parameter, support multiple paths
//Watch: This is the entrypoint for the directory watcher
func (dw *DirectoryWatcherService) Watch() {
	log.Info("Starting directory watcher...")

	dw.downloadsFolderID = utils.GetDownloadsFolderIDFromPremiumizeme(dw.premiumizemeClient)

	log.Info("Clearing tmp directory...")
	tempDir := utils.GetTempBaseDir()
	err := os.RemoveAll(tempDir)
	if err != nil {
		log.Errorf("Error clearing tmp directory %s", tempDir)
	}
	os.Mkdir(tempDir, os.ModePerm)

	log.Info("Creating Queue...")
	dw.Queue = stringqueue.NewStringQueue()

	log.Info("Starting uploads processor...")
	go dw.processUploads()

	log.Info("Starting initial directory scans...")
	go dw.initialDirectoryScan(dw.config.BlackholeDirectory)

	// Build and start a DirectoryWatcher
	watcher := directory_watcher.NewDirectoryWatcher(dw.config.BlackholeDirectory,
		false,
		dw.checkFile,
		dw.addFileToQueue,
	)

	watcher.Watch()
}

func (dw *DirectoryWatcherService) initialDirectoryScan(p string) {
	log.Trace("Initial directory scan")
	files, err := ioutil.ReadDir(p)
	if err != nil {
		log.Errorf("Error with initial directory scan %+v", err)
	}

	for _, file := range files {
		go func(file os.FileInfo) {
			file_path := path.Join(p, file.Name())
			if dw.checkFile(file_path) {
				dw.addFileToQueue(file_path)
			}
		}(file)
	}
}

func (dw *DirectoryWatcherService) checkFile(path string) bool {
	log.Tracef("Checking file %s", path)

	fi, err := os.Stat(path)
	if err != nil {
		log.Errorf("Error checking file %s", path)
		return false
	}

	if fi.IsDir() {
		log.Errorf("Directory created in blackhole %s ignoring (Warning premiumizearrzed does not look in subfolders!)", path)
		return false
	}

	ext := filepath.Ext(path)
	if ext == ".nzb" || ext == ".magnet" {
		return true
	} else {
		return false
	}
}

func (dw *DirectoryWatcherService) addFileToQueue(path string) {
	dw.Queue.Add(path)
	log.Infof("File created in blackhole %s added to Queue. Queue length %d", path, dw.Queue.Len())
}

func (dw *DirectoryWatcherService) processUploads() {
	//TODO: Global running state
	for {
		if dw.Queue.Len() < 1 {
			log.Trace("No files in Queue, sleeping for 10 seconds")
			time.Sleep(time.Second * time.Duration(10))
		}

		isQueueFile, filePath := dw.Queue.PopTopOfQueue()
		if !isQueueFile {
			time.Sleep(time.Second * time.Duration(10))
			continue
		}

		sleepTimeSeconds := 2
		if filePath != "" {
			log.Debugf("Processing %s", filePath)
			err := dw.premiumizemeClient.CreateTransfer(filePath, dw.downloadsFolderID)
			if err != nil {
				switch err.Error() {
				case ERROR_LIMIT_REACHED:
					dw.status = "Limit of transfers reached!"
					log.Trace("Transfer limit reached waiting 10 seconds and retrying")
					sleepTimeSeconds = 10
				case ERROR_ALREADY_UPLOADED:
					log.Trace("File already uploaded, removing from Disk")
					os.Remove(filePath)
				default:
					log.Error("Error creating transfer: %s", err)
				}
			} else {
				dw.status = "Okay"
				os.Remove(filePath)
				if err != nil {
					log.Errorf("Error could not delete %s Error: %+v", filePath, err)
				}
				log.Infof("Removed %s from blackhole Queue. Queue Size: %d", filePath, dw.Queue.Len())
			}
			time.Sleep(time.Second * time.Duration(sleepTimeSeconds))
		} else {
			log.Errorf("Received %s from blackhole Queue. Appears to be an empty path.")
		}
	}
}
