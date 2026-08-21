package httpapi

import "net/http"

func Routes(h *Handler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/records", h)
	mux.Handle("/records/", h)
	return mux
}
func Health(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) }
func WithHealth(h http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/health", http.HandlerFunc(Health))
	mux.Handle("/", h)
	return mux
}
