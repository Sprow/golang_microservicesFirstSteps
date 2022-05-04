package handler

//import (
//	"content_parser/internal/parser"
//	"encoding/json"
//	"github.com/go-chi/chi/v5"
//	"log"
//	"net/http"
//)
//
//type Handler struct {
//	parser *parser.ContentParser
//}
//
//func NewHandler(parser *parser.ContentParser) *Handler {
//	return &Handler{parser: parser}
//}
//
//func (h *Handler) Register(r *chi.Mux) {
//	r.Get("/", h.GetContentByOID)
//}
//
//func (h *Handler) GetContentByOID(w http.ResponseWriter, r *http.Request) {
//	d := json.NewDecoder(r.Body)
//	var content []byte
//	err := d.Decode(&content)
//	if err != nil {
//		log.Println(err)
//		http.Error(w, "Bad request", http.StatusBadRequest)
//		return
//	}
//	//ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
//	//defer cancel()
//	//err = h.parser.ParseAndSave(ctx, content)
//	if err != nil {
//		log.Println(err)
//	}
//}
