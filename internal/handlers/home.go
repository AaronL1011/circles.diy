package handlers

import (
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"circles.diy/internal/models"
	"circles.diy/internal/utils"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles("index.html"))
	data := models.PageData{
		Success:   false,
		CSRFToken: utils.GenerateCSRFToken(),
	}
	tmpl.Execute(w, data)
}

func FeedbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	csrfToken := r.FormValue("csrf_token")
	if csrfToken == "" {
		http.Error(w, "Missing CSRF token", http.StatusBadRequest)
		return
	}

	// Honeypot validation - if these fields are filled, it's likely a bot
	if r.FormValue("website") != "" || r.FormValue("email_address") != "" {
		log.Printf("Bot detected from %s: honeypot fields filled", r.RemoteAddr)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	feedback := r.FormValue("feedback")
	personas := r.Form["persona"]

	validatedFeedback, isValid := utils.ValidateFeedback(feedback)

	if !isValid {
		http.Error(w, "Invalid feedback content", http.StatusBadRequest)
		return
	}

	if validatedFeedback != "" {
		feedbackPath := "/app/data/feedback.txt"

		if err := os.MkdirAll("/app/data", 0750); err != nil {
			log.Printf("Error creating data directory: %v", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		file, err := os.OpenFile(feedbackPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
		if err != nil {
			log.Printf("Error opening feedback file: %v", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		submissionTime := time.Now()
		var personaTexts []string
		if len(personas) == 0 {
			personaTexts = append(personaTexts, "curious-observer")
		} else {
			for _, persona := range personas {
				personaTexts = append(personaTexts, html.EscapeString(persona))
			}
		}
		personaText := strings.Join(personaTexts, ", ")

		if _, err := fmt.Fprintf(file, "=== Feedback received %s ===\\nPersonas: %s\\n%s\\n\\n",
			submissionTime.Format("2006-01-02 15:04:05"), personaText, validatedFeedback); err != nil {
			log.Printf("Error writing feedback: %v", err)
		}

		log.Printf("Feedback received from %s (personas: %s, length: %d)", r.RemoteAddr, personaText, len(validatedFeedback))
	}

	tmpl := template.Must(template.ParseFiles("index.html"))
	data := models.PageData{
		Success:   true,
		CSRFToken: utils.GenerateCSRFToken(),
	}
	tmpl.Execute(w, data)
}