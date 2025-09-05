package handlers

import (
	"log"
	"net/http"
	"strings"

	"circles.diy/internal/templates"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Check if this is the internal profile view (/profile) or external (/profile/:handle)
	if path == "/profile" {
		// Internal profile view (owner's dashboard)
		data := templates.GetMockProfileInternalData()

		err := templates.GetTemplates().ProfileInternal.ExecuteTemplate(w, "profile-internal", data)
		if err != nil {
			log.Printf("Error rendering profile-internal template: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else if strings.HasPrefix(path, "/profile/") {
		// External profile view (/profile/:handle)
		// Extract handle from URL path
		handle := strings.TrimPrefix(path, "/profile/")
		if handle == "" {
			http.NotFound(w, r)
			return
		}

		// For demo purposes, we'll use the mock data regardless of handle
		// In a real app, you'd look up the user by handle
		data := templates.GetMockProfileData(handle, false) // isOwner = false for external view

		err := templates.GetTemplates().ProfilePublic.ExecuteTemplate(w, "profile-public", data)
		if err != nil {
			log.Printf("Error rendering profile-public template: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else {
		http.NotFound(w, r)
	}
}