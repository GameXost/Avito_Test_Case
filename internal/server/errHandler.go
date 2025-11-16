package server

import (
	"encoding/json"
	"github.com/GameXost/Avito_Test_Case/models"
	"net/http"
)

func WriteError(w http.ResponseWriter, status int, code error, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := models.ErrorResponse{
		Error: models.ErrorMessage{
			Code:    code,
			Message: message,
		},
	}

	_ = json.NewEncoder(w).Encode(resp)
}
