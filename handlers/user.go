package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/kazimanzurrashid/consents-api-go/models"
	"github.com/kazimanzurrashid/consents-api-go/services"
)

type User struct {
	srv services.User
}

func NewUser(srv services.User) *User {
	return &User{srv}
}

func (h *User) Create(w http.ResponseWriter, r *http.Request) {
	var req models.UserCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusUnprocessableEntity, "Malformed request body")
		return
	}

	if err := req.Validate(); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	user, err := h.srv.Create(r.Context(), &req)

	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, "Email already exists")
		return
	}

	writeSuccess(w, http.StatusCreated, user)
}

func (h *User) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := h.srv.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *User) Detail(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	user, err := h.srv.Detail(r.Context(), id)

	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if user == nil {
		writeError(w, http.StatusNotFound, "User not found")
		return
	}

	writeSuccess(w, http.StatusOK, user)
}
