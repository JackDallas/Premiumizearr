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
	watchDirectory     *directory_watcher.WatchDirectory
}

const (
	ERROR_LIMIT_REACHED    = "Limit of transfers reached!"
	ERROR_ALREADY_UPLOADED = "You already added this job."
)

func (DirectoryWatcherService) New() DirectoryWatcherService {
	return DirectoryWatcherService{
		premiumizemeClient: nil,
		config:             nil,
		Queue:              nil,
		status:             "",
		downloadsFolderID:  "",
	}
}

func (dw *DirectoryWatcherService) Init(premiumizemeClient *premiumizeme.Premiumizeme, config *config.Config) {
	dw.premiumizemeClient = premiumizemeClient
	dw.config = config
}

func (dw *DirectoryWatcherService) ConfigUpdatedCallback(currentConfig config.Config, newConfig config.Config) {
	if currentConfig.BlackholeDirectory != newConfig.BlackholeDirectory {
		log.Info("Blackhole directory changed, restarting directory watcher...")
		log.Info("Running initial directory scan...")
		go dw.initialDirectoryScan(dw.config.BlackholeDirectory)
		dw.watchDirectory.UpdatePath(newConfig.BlackholeDirectory)
	}
}

func (dw *DirectoryWatcherService) GetStatus() string {
	return dw.status
}

//TODO (Radarr): accept paths as a parameter, support multiple paths
//Watch: This is the entrypoint for the directory watcher
func (dw *DirectoryWatcherService) Start() {
	log.Info("Starting directory watcher...")

	dw.downloadsFolderID = utils.GetDownloadsFolderIDFromPremiumizeme(dw.premiumizemeClient)

	log.Info("Clearing tmp directory...")
	tempDir := dw.config.GetTempBaseDir()

	if tempDir == "/" || tempDir == "\\" || tempDir == "C:\\" {
		panic("Unzip directory is set to system root don't do that it will wipe your system!")
	}

	err := os.RemoveAll(tempDir)
	if err != nil {
		log.Errorf("Error clearing tmp directory %s", tempDir)
	}
	os.Mkdir(tempDir, os.ModePerm)

	log.Info("Creating Queue...")
	dw.Queue = stringqueue.NewStringQueue()

	log.Info("Starting uploads processor...")
	go dw.processUploads()

	log.Info("Running initial directory scan...")
	go dw.initialDirectoryScan(dw.config.BlackholeDirectory)

	// Build and start a DirectoryWatcher
	dw.watchDirectory = directory_watcher.NewDirectoryWatcher(dw.config.BlackholeDirectory,
		false,
		dw.checkFile,
		dw.addFileToQueue,
	)

	dw.watchDirectory.Watch()
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
