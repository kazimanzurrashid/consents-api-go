package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kazimanzurrashid/consents-api-go/models"
	"github.com/kazimanzurrashid/consents-api-go/services"
)

type Event struct {
	srv services.Event
}

func NewEvent(srv services.Event) *Event {
	return &Event{srv}
}

func (h *Event) Create(w http.ResponseWriter, r *http.Request) {
	var req models.EventCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusUnprocessableEntity, "Malformed request body")
		return
	}

	if err := req.Validate(); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err := h.srv.Create(r.Context(), &req); err != nil {
		writeError(w, http.StatusUnprocessableEntity, "Invalid request")
		return
	}

	w.WriteHeader(http.StatusCreated)
}
