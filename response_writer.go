package gen

import "net/http"

type responseWriter struct {
	http.ResponseWriter
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w}
}

func (w *responseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}
