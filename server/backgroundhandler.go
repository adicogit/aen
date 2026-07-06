package server

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"aen.it/poolmanager/config"
	"aen.it/poolmanager/log"
)

// listBackgrounds returns the list of image filenames available in the background folder
func (server *Server) listBackgrounds(w http.ResponseWriter, r *http.Request) {
	log.Log.Debug("Entering listBackgrounds")
	w.Header().Set("Content-Type", "application/json")

	folder := config.UIConfig.UIAspect.BackgroundImageFolder
	// folder is expected to be like /images/background
	dir := filepath.Join(WebuiHandler.staticPath, folder)

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Log.Error("Unable to read backgrounds dir", "dir", dir, "error", err)
		http.Error(w, "Unable to read backgrounds", http.StatusInternalServerError)
		return
	}

	var result []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		ext := strings.ToLower(filepath.Ext(name))
		switch ext {
		case ".png", ".jpg", ".jpeg", ".webp", ".jfif", ".gif", ".bmp":
			result = append(result, name)
		}
	}

	json.NewEncoder(w).Encode(result)
	log.Log.Debug("Exiting listBackgrounds", "count", len(result))
}
