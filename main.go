package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
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
		// Log feedback to file
		file, err := os.OpenFile("feedback.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Error opening feedback file: %v", err)
		} else {
			defer file.Close()
			fmt.Fprintf(file, "=== Feedback received ===\n%s\n\n", feedback)
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
