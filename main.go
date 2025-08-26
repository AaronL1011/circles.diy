package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type PageData struct {
	Success bool
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	data := PageData{Success: false}
	tmpl.Execute(w, data)
}

func feedbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	feedback := r.FormValue("feedback")
	if feedback != "" {
		// Log feedback to file in data directory
		feedbackPath := "/app/data/feedback.txt"

		// Ensure data directory exists
		if err := os.MkdirAll("/app/data", 0755); err != nil {
			log.Printf("Error creating data directory: %v", err)
		}

		file, err := os.OpenFile(feedbackPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Error opening feedback file: %v", err)
		} else {
			defer file.Close()
			submissionTime := time.Now()
			fmt.Fprintf(file, "=== Feedback received %s ===\n%s\n\n", submissionTime.Format("2006-01-02 15:04:05"), feedback)
		}
		log.Printf("Feedback received: %s", feedback)
	}

	// Show success message
	tmpl := template.Must(template.ParseFiles("index.html"))
	data := PageData{Success: true}
	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/feedback", feedbackHandler)

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
