package handlers

import (
	"log"
	"net/http"

	"circles.diy/internal/templates"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	data := templates.GetMockDashboardData()

	err := templates.GetTemplates().Dashboard.ExecuteTemplate(w, "dashboard", data)
	if err != nil {
		log.Printf("Error rendering dashboard template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}