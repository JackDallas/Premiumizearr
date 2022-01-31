package web_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackdallas/premiumizearr/internal/config"
	"github.com/jackdallas/premiumizearr/internal/service"
	"github.com/jackdallas/premiumizearr/pkg/premiumizeme"
	log "github.com/sirupsen/logrus"
)

type IndexTemplates struct {
	RootPath string
}

var indexBytes []byte

const webRoot = "premiumizearr"

type server struct {
	transferManager         *service.TransferManagerService
	directoryWatcherService *service.DirectoryWatcherService
}

// http Router
func StartWebServer(transferManager *service.TransferManagerService, directoryWatcher *service.DirectoryWatcherService, config *config.Config) {
	tmpl, err := template.ParseFiles("./static/index.html")
	if err != nil {
		log.Fatal(err)
	}

	var ibytes bytes.Buffer
	err = tmpl.Execute(&ibytes, &IndexTemplates{webRoot})
	if err != nil {
		log.Fatal(err)
	}
	indexBytes = ibytes.Bytes()

	s := server{
		transferManager:         transferManager,
		directoryWatcherService: directoryWatcher,
	}
	spa := spaHandler{
		staticPath: "static",
		indexPath:  "index.html",
	}

	r := mux.NewRouter()

	log.Infof("Creating route: %s", webRoot+"/api/transfers")
	r.HandleFunc("/"+webRoot+"/api/transfers", s.TransfersHandler)

	log.Infof("Creating route: %s", webRoot+"/api/downloads")
	r.HandleFunc("/"+webRoot+"/api/downloads", s.DownloadsHandler)

	log.Infof("Creating route: %s", webRoot+"/api/blackhole")
	r.HandleFunc("/"+webRoot+"/api/blackhole", s.BlackholeHandler)

	r.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("%s:%s", config.BindIP, config.BindPort),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	srv.ListenAndServe()
}

type TransfersResponse struct {
	Transfers []premiumizeme.Transfer `json:"data"`
	Status    string                  `json:"status"`
}

func (s *server) TransfersHandler(w http.ResponseWriter, r *http.Request) {
	var resp TransfersResponse
	resp.Transfers = *s.transferManager.GetTransfers()
	resp.Status = s.transferManager.GetStatus()
	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

type BlackholeFile struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type BlackholeResponse struct {
	BlackholeFiles []BlackholeFile `json:"data"`
	Status         string          `json:"status"`
}

type Download struct {
	Added    int64  `json:"added"`
	Name     string `json:"name"`
	Progress string `json:"progress"`
	Speed    string `json:"speed"`
}
type DownloadsResponse struct {
	Downloads []Download `json:"data"`
	Status    string     `json:"status"`
}

func (s *server) DownloadsHandler(w http.ResponseWriter, r *http.Request) {
	var resp DownloadsResponse

	for _, v := range s.transferManager.GetDownloads() {
		resp.Downloads = append(resp.Downloads, Download{
			Added:    v.Added.Unix(),
			Name:     v.Name,
			Progress: v.ProgressDownloader.GetProgress(),
			Speed:    v.ProgressDownloader.GetSpeed(),
		})
	}
	resp.Status = ""

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func (s *server) BlackholeHandler(w http.ResponseWriter, r *http.Request) {
	var resp BlackholeResponse
	for i, n := range s.directoryWatcherService.Queue.GetQueue() {
		name := path.Base(n)
		resp.BlackholeFiles = append(resp.BlackholeFiles, BlackholeFile{
			ID:   i,
			Name: name,
		})
	}
	resp.Status = s.directoryWatcherService.GetStatus()

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

// Shamlessly stolen from mux examples https://github.com/gorilla/mux#examples
type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path = strings.Replace(path, webRoot, "", 1)

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) || strings.HasSuffix(path, h.staticPath) {
		// file does not exist, serve index.html
		// http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		// file does not exist, serve index.html template
		w.Write(indexBytes)
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r.URL.Path = strings.Replace(path, h.staticPath, "", -1)
	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}
