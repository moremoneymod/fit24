package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"fit24/domain"
)

type HTTPHandler struct {
	service   domain.LeadService
	dbTimeout time.Duration
}

func NewHTTPHandler(service domain.LeadService, dbTimeout time.Duration) *HTTPHandler {
	return &HTTPHandler{
		service:   service,
		dbTimeout: dbTimeout,
	}
}

type orderDTO struct {
	Plan  string `json:"plan"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

type contactDTO struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

type response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func (h *HTTPHandler) HandleOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendJSON(w, http.StatusMethodNotAllowed, response{Success: false, Message: "Метод не разрешен"})
		return
	}

	var dto orderDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.sendJSON(w, http.StatusBadRequest, response{Success: false, Message: "Некорректный запрос"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.dbTimeout)
	defer cancel()

	_, err := h.service.SubmitOrder(ctx, dto.Plan, dto.Name, dto.Phone, dto.Email)
	if err != nil {
		h.sendJSON(w, http.StatusBadRequest, response{Success: false, Message: err.Error()})
		return
	}

	h.sendJSON(w, http.StatusOK, response{Success: true})
}

func (h *HTTPHandler) HandleContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendJSON(w, http.StatusMethodNotAllowed, response{Success: false, Message: "Метод не разрешен"})
		return
	}

	var dto contactDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.sendJSON(w, http.StatusBadRequest, response{Success: false, Message: "Некорректный запрос"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.dbTimeout)
	defer cancel()

	_, err := h.service.SubmitContact(ctx, dto.Name, dto.Phone, dto.Message)
	if err != nil {
		h.sendJSON(w, http.StatusBadRequest, response{Success: false, Message: err.Error()})
		return
	}

	h.sendJSON(w, http.StatusOK, response{Success: true})
}

func (h *HTTPHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	h.sendJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

func (h *HTTPHandler) sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
