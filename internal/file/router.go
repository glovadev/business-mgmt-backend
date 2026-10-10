package file

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, authMiddleware func(http.Handler) http.Handler) {
	mux.Handle("POST /api/v1/files/upload-url", authMiddleware(http.HandlerFunc(handler.GetUploadURL)))
	mux.Handle("GET /api/v1/files/download-url", authMiddleware(http.HandlerFunc(handler.GetDownloadURL)))
}
