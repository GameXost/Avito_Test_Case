package errHandle

import (
	"encoding/json"
	"github.com/GameXost/Avito_Test_Case/models"
	"net/http"
)

func WriteError(w http.ResponseWriter, status int, errCode error, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := models.ErrorResponse{
		Error: models.ErrorMessage{
			Code:    errCode.Error(),
			Message: message,
		},
	}

	_ = json.NewEncoder(w).Encode(resp)
}
