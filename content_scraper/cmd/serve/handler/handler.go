package handler

import (
	"ContentScraper/internal/manager"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

type Handler struct {
	contentScraperManager *manager.ContentScraperManager
}

func NewHandler(contentScraperManager *manager.ContentScraperManager) *Handler {
	return &Handler{
		contentScraperManager: contentScraperManager,
	}
}

func (h *Handler) Register(r *chi.Mux) {
	r.Post("/", h.SentContentToParserByOID)
}

func (h *Handler) SentContentToParserByOID(w http.ResponseWriter, r *http.Request) {
	oid, err := primitive.ObjectIDFromHex(r.FormValue("oid"))
	if err != nil {
		log.Println(err)
	}
	content, err := h.contentScraperManager.DB.GetContentByOID(r.Context(), oid)
	if err != nil {
		log.Println(err)
	}

	err = json.NewEncoder(w).Encode(content)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
