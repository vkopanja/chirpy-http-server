package handler

import "net/http"

func NewHealth() *Health {
	return &Health{}
}

type Health struct{}

// ServeHTTP
// @Summary Health check
// @Description Health check
// @Tags Health handler
// @Produce text/plain
// @Success 200 {string} string "OK"
// @Router /api/healthz [get]
func (h *Health) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		panic(err)
	}
}
