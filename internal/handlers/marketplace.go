package handlers

import (
	"log"
	"net/http"

	"circles.diy/internal/templates"
)

func MarketplaceHandler(w http.ResponseWriter, r *http.Request) {
	// Get mock marketplace data
	data := templates.GetMockMarketplaceData()

	// Render the marketplace template
	err := templates.GetTemplates().Marketplace.ExecuteTemplate(w, "marketplace", data)
	if err != nil {
		log.Printf("Error rendering marketplace template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
