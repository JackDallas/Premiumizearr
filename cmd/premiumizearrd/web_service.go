package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

// http Router
func StartWebServer(pme *premiumizearrd) {
	spa := spaHandler{staticPath: "static", indexPath: "index.html"}

	r := mux.NewRouter()

	r.HandleFunc("/api/transfers", pme.TransfersHandler)
	r.HandleFunc("/api/downloads", pme.DownloadsHandler)

	r.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("%s:%s", pme.Config.BindIP, pme.Config.BindPort),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	srv.ListenAndServe()
}

func (pme *premiumizearrd) TransfersHandler(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(pme.TransferManager.Transfers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func (pme *premiumizearrd) DownloadsHandler(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(pme.DirectoryWatcher.Queue.GetQueue())
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

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}
