package handlers

import (
	"log"
	"net/http"

	"circles.diy/internal/templates"
)

func CirclesHandler(w http.ResponseWriter, r *http.Request) {
	data := templates.GetMockCirclesPageData()

	err := templates.GetTemplates().Circles.ExecuteTemplate(w, "circles", data)
	if err != nil {
		log.Printf("Error rendering circles template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}