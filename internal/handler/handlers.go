package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/yyek0/stroydom-website/internal/database"
	"github.com/yyek0/stroydom-website/internal/models"
)

type Handlers struct {
	Database database.LeadStorage
}

func NewHandler(db database.LeadStorage) *Handlers {
	return &Handlers{
		Database: db,
	}
}

func (h *Handlers) HandleCheckHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := h.Database.Ping(ctx); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) HandleCreateLead(w http.ResponseWriter, r *http.Request) {
	var lead models.Lead
	if err := json.NewDecoder(r.Body).Decode(&lead); err != nil {
		// will be able to be logged with errdto
		errdto := models.ErrorDTO{
			Msg:  err.Error(),
			Time: time.Now(),
		}

		http.Error(w, errdto.ToString(), http.StatusBadRequest)
		return
	}

	if err := lead.Validate(); err != nil {
		// will be able to be logged with errdto

		errdto := models.ErrorDTO{
			Msg:  err.Error(),
			Time: time.Now(),
		}

		http.Error(w, errdto.ToString(), http.StatusBadRequest)
		return
	}
	lead.CreatedAt = time.Now()
	if err := h.Database.Create(r.Context(), lead); err != nil {
		// will be able to be logged with errdto

		errdto := models.ErrorDTO{
			Msg:  err.Error(),
			Time: time.Now(),
		}

		http.Error(w, errdto.ToString(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := struct {
		Status  string      `json:"status"`
		Message string      `json:"message"`
		Data    models.Lead `json:"data"`
	}{
		Status:  "Created",
		Message: "Заявка успешно создана",
		Data:    lead,
	}

	json.NewEncoder(w).Encode(response)

}

func (h *Handlers) HandleGetAllLeads(w http.ResponseWriter, r *http.Request) {
	leads, err := h.Database.GetAll(r.Context())
	if err != nil {
		// will be able to be logged with errdto
		errdto := models.ErrorDTO{
			Msg:  err.Error(),
			Time: time.Now(),
		}

		http.Error(w, errdto.ToString(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(leads)
}

func (h *Handlers) HandleGetLead(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		// will be able to be logged with errdto
		err := errors.New("id not found")
		errdto := models.ErrorDTO{
			Msg:  err.Error(),
			Time: time.Now(),
		}

		http.Error(w, errdto.ToString(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		// will be able to be logged with errdto
		errdto := models.ErrorDTO{
			Msg:  err.Error(),
			Time: time.Now(),
		}

		http.Error(w, errdto.ToString(), http.StatusBadRequest)
		return
	}

	lead, err := h.Database.Get(r.Context(), id)

	if err != nil {
		errdto := models.ErrorDTO{
			Msg:  err.Error(),
			Time: time.Now(),
		}

		http.Error(w, errdto.ToString(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(lead)
}

func (h *Handlers) HandleDeleteLead(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		// will be able to be logged with errdto
		err := errors.New("id not found")
		errdto := models.ErrorDTO{
			Msg:  err.Error(),
			Time: time.Now(),
		}

		http.Error(w, errdto.ToString(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		// will be able to be logged with errdto
		errdto := models.ErrorDTO{
			Msg:  err.Error(),
			Time: time.Now(),
		}

		http.Error(w, errdto.ToString(), http.StatusBadRequest)
		return
	}

	if err := h.Database.Delete(r.Context(), id); err != nil {
		// will be able to be logged with errdto
		errdto := models.ErrorDTO{
			Msg:  err.Error(),
			Time: time.Now(),
		}

		http.Error(w, errdto.ToString(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
