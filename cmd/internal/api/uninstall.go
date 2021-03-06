package api

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd"
)

//UninstallRequest for incoming uninstallation requests
type UninstallRequest struct {
	Game  *string  `json:"game"`
	Names []string `json:"names"`
}

//UninstallResponse for uninstallation requests
type UninstallResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message,omitempty"`
	Stats   *cmd.Stats `json:"stats,omitempty"`
}

//Uninstall a mod via api request
func Uninstall(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	setHeaders(w)

	decoder := json.NewDecoder(r.Body)
	stats, err := uninstall(decoder)

	encoder := json.NewEncoder(w)
	if err == nil {
		encoder.Encode(&UninstallResponse{
			Success: true,
			Stats:   stats,
		})
	} else {
		encoder.Encode(&UninstallResponse{
			Success: false,
			Message: err.Error(),
			Stats:   stats,
		})
	}
}

func uninstall(decoder *json.Decoder) (*cmd.Stats, error) {
	var req UninstallRequest
	if err := decoder.Decode(&req); err != nil {
		return nil, fmt.Errorf("cmd/internal/api: Could not parse request body: %s", err.Error())
	}

	if req.Game != nil {
		if err := flag.Set("game", *req.Game); err != nil {
			return nil, fmt.Errorf("cmd/internal/api: Could set game flag: %s", err.Error())
		}
	}

	return cmd.Uninstall(req.Names)
}
