package handler

import (
	"encoding/json"
	"net/http"
	"short-link/cmd/internal/service"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	short, err := h.svc.Shorten(req.URL)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"short": short,
	})
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	short := r.URL.Path[1:]
	url, err := h.svc.Resolve(short)

	if err != nil || url == ""{
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}
