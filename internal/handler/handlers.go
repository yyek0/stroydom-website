package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/yyek0/stroydom-website/internal/models"
)

type DummyDB struct {
}

func (d *DummyDB) Ping() error {
	return errors.New("abaob")
}

type Handlers struct {
	dummyDB *DummyDB
}

func NewHandler(d *DummyDB) *Handlers {
	return &Handlers{
		dummyDB: d,
	}
}

func (h *Handlers) HandleCheckHealth(w http.ResponseWriter, r *http.Request) {
	if err := h.dummyDB.Ping(); err != nil {
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
		errdto := models.ErrorDTO{
			Msg:  err.Error(),
			Time: time.Now(),
		}

		http.Error(w, errdto.ToString(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := struct {
		Status  string      `json:"status"`
		Message string      `json:"message"`
		Data    models.Lead `json:"data"`
	}{
		Status:  "ok",
		Message: "Заявка успешно создана",
		Data:    lead,
	}

	json.NewEncoder(w).Encode(response)

}
