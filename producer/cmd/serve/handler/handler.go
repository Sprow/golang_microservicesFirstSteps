package handler

import (
	"Producer/internal/producer"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"net/url"
)

type Handler struct {
	producer *producer.Producer
}

func NewHandler(producer *producer.Producer) *Handler {
	return &Handler{
		producer: producer,
	}
}

func (h *Handler) Register(r *chi.Mux) {
	r.Post("/", h.sandUrl)

}

type reqURL struct {
	URL string `json:"url"`
}

func (h *Handler) sandUrl(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	var newURL reqURL
	err := d.Decode(&newURL)
	if err != nil {
		log.Println(err)
	}
	_, err = url.ParseRequestURI(newURL.URL)
	if err != nil {
		http.Error(w, "Bad url", http.StatusInternalServerError)
		log.Println(err)
	} else {
		_ = h.producer.PublishURL(newURL.URL)
	}
}
