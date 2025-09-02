package handlers

import (
	"net/http"
	"os"
	"strings"
	"time"
)

func ServeStaticFile(w http.ResponseWriter, r *http.Request, filePath, contentType string) {
	// Security: prevent path traversal
	if strings.Contains(filePath, "..") {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	// Get file info for cache headers
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Set cache headers for production performance
	lastModified := fileInfo.ModTime().UTC().Format(http.TimeFormat)
	ifModifiedSince := r.Header.Get("If-Modified-Since")

	if ifModifiedSince != "" {
		if t, err := time.Parse(http.TimeFormat, ifModifiedSince); err == nil {
			if fileInfo.ModTime().UTC().Before(t.Add(1 * time.Second)) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}
	}

	// Set response headers
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Last-Modified", lastModified)
	w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour

	// For CSS files, set additional headers
	if strings.HasSuffix(filePath, ".css") {
		w.Header().Set("Vary", "Accept-Encoding")
	}

	// Serve the file
	http.ServeFile(w, r, filePath)
}