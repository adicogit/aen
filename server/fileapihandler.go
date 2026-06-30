package server

import (
	"encoding/json"
	"net/http"

	"aen.it/poolmanager/config"
	"aen.it/poolmanager/log"
)

type uiConfig struct {
	Background_image string `json:"background_image,omitempty"`
	Theme            string `json:"theme,omitempty"`
	BilliardRoomName string `json:"billiard_room_name,omitempty"`
}

// return the Web UI configuration
func (server *Server) getUIConfig(w http.ResponseWriter, r *http.Request) {
	log.Log.Debug("Entering getBilliardManagerProperties")
	w.Header().Set("Content-Type", "application/json")
	result := uiConfig{
		Background_image: config.UIConfig.UIAspect.BackgroundImage,
		Theme:            config.UIConfig.UIAspect.Theme,
		BilliardRoomName: config.UIConfig.UIAspect.BilliardRoomName,
	}

	log.Log.Debug("Exiting getBilliardManagerProperties", "result", result)
	json.NewEncoder(w).Encode(result)
}

// return the Web UI configuration
func (server *Server) setUIConfig(w http.ResponseWriter, r *http.Request) {
	log.Log.Debug("Entering setUIConfig")
	w.Header().Set("Content-Type", "application/json")
	requestData := uiConfig{}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		log.Log.Error("Unable to decode received body", "body", r.Body, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	config.UIConfig.UIAspect.BackgroundImage = requestData.Background_image
	config.UIConfig.UIAspect.Theme = requestData.Theme
	config.UIConfig.UIAspect.BilliardRoomName = requestData.BilliardRoomName
	err = config.UIConfig.SaveConfig()
	if err != nil {
		log.Log.Error("Unable save UI configuration", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := struct {
		Result string `json:"result"`
	}{
		Result: "success",
	}
	log.Log.Debug("Exiting setUIConfig", "result", result)
	json.NewEncoder(w).Encode(result)
}
