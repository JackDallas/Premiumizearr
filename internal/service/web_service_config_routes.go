package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackdallas/premiumizearr/internal/config"
)

type ConfigChangeResponse struct {
	Succeeded bool   `json:"succeeded"`
	Status    string `json:"status"`
}

func (s *WebServerService) ConfigHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		data, err := json.Marshal(s.config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(data)
	case http.MethodPost:
		var newConfig config.Config
		err := json.NewDecoder(r.Body).Decode(&newConfig)
		if err != nil {
			EncodeAndWriteConfigChangeResponse(w, &ConfigChangeResponse{
				Succeeded: false,
				Status:    fmt.Sprintf("Config failed to update %s", err.Error()),
			})
			return
		}
		s.config.UpdateConfig(newConfig)
		EncodeAndWriteConfigChangeResponse(w, &ConfigChangeResponse{
			Succeeded: true,
			Status:    "Config updated",
		})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func EncodeAndWriteConfigChangeResponse(w http.ResponseWriter, resp *ConfigChangeResponse) {
	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}
