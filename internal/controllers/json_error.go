package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/randnull/banner-service/pkg/models"
)


func sendJSONError(w http.ResponseWriter, statusCode int, Message error) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    ErrorResponse := models.ErrorModel{Error: Message.Error()}
    json.NewEncoder(w).Encode(ErrorResponse)
}
