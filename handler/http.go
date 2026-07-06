package handler

import (
	"encoding/json"
	"net/http"

	"fit24/domain"
)

type HTTPHandler struct {
	service domain.LeadService
}

func NewHTTPHandler(service domain.LeadService) *HTTPHandler {
	return &HTTPHandler{service: service}
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
		h.sendJSON(w, http.StatusMethodNotAllowed, response{Success: false, Message: "Метод не поддерживается"})
		return
	}

	var dto orderDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.sendJSON(w, http.StatusBadRequest, response{Success: false, Message: "Некорректный запрос"})
		return
	}

	_, err := h.service.SubmitOrder(r.Context(), dto.Plan, dto.Name, dto.Phone, dto.Email)
	if err != nil {
		h.sendJSON(w, http.StatusBadRequest, response{Success: false, Message: err.Error()})
		return
	}

	h.sendJSON(w, http.StatusOK, response{Success: true})
}

func (h *HTTPHandler) HandleContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendJSON(w, http.StatusMethodNotAllowed, response{Success: false, Message: "Метод не поддерживается"})
		return
	}

	var dto contactDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.sendJSON(w, http.StatusBadRequest, response{Success: false, Message: "Некорректный запрос"})
		return
	}

	_, err := h.service.SubmitContact(r.Context(), dto.Name, dto.Phone, dto.Message)
	if err != nil {
		h.sendJSON(w, http.StatusBadRequest, response{Success: false, Message: err.Error()})
		return
	}

	h.sendJSON(w, http.StatusOK, response{Success: true})
}

func (h *HTTPHandler) sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}