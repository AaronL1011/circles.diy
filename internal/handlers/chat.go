package handlers

import (
	"log"
	"net/http"

	"circles.diy/internal/templates"
)

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	// Get mock chat data
	data := templates.GetMockChatData()
	
	// Render the chat template
	err := templates.GetTemplates().Chat.ExecuteTemplate(w, "chat", data)
	if err != nil {
		log.Printf("Error rendering chat template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}