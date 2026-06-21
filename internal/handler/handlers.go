package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/yyek0/stroydom-website/internal/database"
	"github.com/yyek0/stroydom-website/internal/models"
	"go.uber.org/zap"
)

type Handlers struct {
	Database database.LeadStorage
	Logger   *zap.Logger
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	ReqID   string `json:"req_id"`
}

func (h *Handlers) respondWithError(w http.ResponseWriter, userMsg string, statusCode int, realErr error) {
	// 1. Генерим уникальный ID ошибки
	reqID := fmt.Sprintf("REQ-%d", time.Now().UnixNano())

	h.Logger.Error("Сбой при обработке запроса",
		zap.String("req_id", reqID),
		zap.String("user_msg", userMsg),
		zap.Error(realErr),
	)

	response := ErrorResponse{
		Status:  "error",
		Message: userMsg,
		ReqID:   reqID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func NewHandler(db database.LeadStorage, lg *zap.Logger) *Handlers {
	return &Handlers{
		Database: db,
		Logger:   lg,
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
		h.respondWithError(w, "Не удалось сохранить заявку", http.StatusBadRequest, err)
		return
	}

	if err := lead.Validate(); err != nil {
		h.respondWithError(w, "Неверный формат данных", http.StatusBadRequest, err)
		return
	}

	lead.CreatedAt = time.Now()

	id, err := h.Database.Create(r.Context(), lead)
	if err != nil {
		h.respondWithError(w, "Не удалось сохранить заявку", http.StatusInternalServerError, err)
		return
	}

	h.Logger.Info("Заявка успешно создана",
		zap.Int("lead_id", id),
		zap.String("name", lead.Name),
	)

	lead.ID = id

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Заявка успешно создана",
		"data":    lead,
	})

}

func (h *Handlers) HandleGetAllLeads(w http.ResponseWriter, r *http.Request) {
	leads, err := h.Database.GetAll(r.Context())
	if err != nil {
		h.respondWithError(w, "Не удалось получить все заявки", http.StatusInternalServerError, err)
		return
	}

	h.Logger.Info("Успешно получены все заявки", zap.Int("count", len(leads)))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(leads)
}

func (h *Handlers) HandleGetLead(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.respondWithError(w, "Не передан ID заявки", http.StatusBadRequest, errors.New("empty id query parameter"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.respondWithError(w, "Неверный формат ID", http.StatusBadRequest, err)
		return
	}

	lead, err := h.Database.Get(r.Context(), id)

	if err != nil {
		h.respondWithError(w, "Заявка не найдена", http.StatusNotFound, err)
		return
	}

	h.Logger.Info("Заявка успешно получена", zap.Int("lead_id", id))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(lead)
}

func (h *Handlers) HandleDeleteLead(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.respondWithError(w, "Не передан ID заявки", http.StatusBadRequest, errors.New("empty id query parameter"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.respondWithError(w, "Неверный формат ID", http.StatusBadRequest, err)
		return
	}

	if err := h.Database.Delete(r.Context(), id); err != nil {
		h.respondWithError(w, "Не удалось удалить заявку", http.StatusInternalServerError, err)
		return
	}

	h.Logger.Info("Заявка успешно удалена", zap.Int("lead_id", id))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
