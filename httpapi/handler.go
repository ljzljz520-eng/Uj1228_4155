package httpapi

import (
	"encoding/json"
	"gestureflame/domain"
	"gestureflame/service"
	"net/http"
	"strings"
)

type Handler struct{ app *service.Service }

func New(app *service.Service) *Handler { return &Handler{app: app} }
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.handleGet(w, r)
		return
	}
	if r.Method == http.MethodPost {
		h.handlePost(w, r)
		return
	}
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}
func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "records" {
		rs, err := h.app.Search(r.URL.Query().Get("q"))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		writeJSON(w, rs)
		return
	}
	id := strings.TrimPrefix(path, "records/")
	rec, events, err := h.app.Recall(id)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	writeJSON(w, map[string]any{"record": rec, "events": events})
}
func (h *Handler) handlePost(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	var payload map[string]string
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if path == "records" {
		d := domain.ImportDraft{BatchCode: payload["batch_code"], Title: payload["title"], Result: payload["result"], Actor: payload["actor"], AttachmentName: "handoff", AttachmentType: "text/plain", AttachmentDigest: "inline"}
		rec, err := h.app.Create(d)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		writeJSON(w, rec)
		return
	}
	parts := strings.Split(path, "/")
	if len(parts) != 3 {
		http.Error(w, "unknown route", 404)
		return
	}
	id, action := parts[1], parts[2]
	var rec domain.Record
	var err error
	switch action {
	case "review":
		rec, err = h.app.Review(id, payload["actor"])
	case "confirm":
		rec, err = h.app.Confirm(id, payload["actor"])
	case "archive":
		rec, err = h.app.Archive(id, payload["actor"])
	case "publish":
		rec, err = h.app.Publish(id, payload["actor"])
	default:
		http.Error(w, "unknown action", 404)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	writeJSON(w, rec)
}
func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}
