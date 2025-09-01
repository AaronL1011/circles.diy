package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type PageData struct {
	Success   bool
	CSRFToken string
}

type RateLimiter struct {
	visitors map[string]*rate.Limiter
	mu       sync.RWMutex
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*rate.Limiter),
	}
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.visitors[ip]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		limiter = rate.NewLimiter(rate.Every(time.Minute), 10)
		rl.visitors[ip] = limiter
		rl.mu.Unlock()
	}

	return limiter
}

func (rl *RateLimiter) Allow(ip string) bool {
	return rl.getLimiter(ip).Allow()
}

var rateLimiter = NewRateLimiter()

func generateCSRFToken() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		log.Printf("Error generating CSRF token: %v", err)
		return ""
	}
	return hex.EncodeToString(bytes)
}

func validateFeedback(feedback string) (string, bool) {
	feedback = strings.TrimSpace(feedback)

	if len(feedback) == 0 {
		return "", false
	}

	if len(feedback) > 5000 {
		return "", false
	}

	maliciousPatterns := []string{
		`<script`, `javascript:`, `onload=`, `onerror=`,
		`eval\(`, `document\.`, `window\.`, `alert\(`,
	}

	feedbackLower := strings.ToLower(feedback)
	for _, pattern := range maliciousPatterns {
		matched, _ := regexp.MatchString(pattern, feedbackLower)
		if matched {
			return "", false
		}
	}

	feedback = html.EscapeString(feedback)

	return feedback, true
}

func securityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		
		// Different CSP for concept demo vs main site
		if strings.HasPrefix(r.URL.Path, "/concept-demo") {
			w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline'; img-src 'self' https://images.unsplash.com https://unsplash.com https://cdn.britannica.com https://media.tenor.com data:; font-src 'self'")
		} else {
			w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'none'")
		}
		
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

		next.ServeHTTP(w, r)
	})
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = strings.Split(forwarded, ",")[0]
		}

		if !rateLimiter.Allow(strings.TrimSpace(ip)) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles("index.html"))
	data := PageData{
		Success:   false,
		CSRFToken: generateCSRFToken(),
	}
	tmpl.Execute(w, data)
}

func conceptDemoHandler(w http.ResponseWriter, r *http.Request) {
	// Handle concept demo routes
	switch r.URL.Path {
	case "/concept-demo":
		// Redirect to ensure trailing slash
		http.Redirect(w, r, "/concept-demo/", http.StatusMovedPermanently)
		return
	case "/concept-demo/":
		// Serve the HTML file
		serveStaticFile(w, r, "concept-demo/index.html", "text/html; charset=utf-8")
		return
	case "/concept-demo/app.css":
		// Serve the CSS file
		serveStaticFile(w, r, "concept-demo/app.css", "text/css; charset=utf-8")
		return
	case "/concept-demo/profile-internal.html":
		// Serve the internal profile page (edit view)
		serveStaticFile(w, r, "concept-demo/profile-internal.html", "text/html; charset=utf-8")
		return
	case "/concept-demo/profile-external.html":
		// Serve the external profile page (public view)
		serveStaticFile(w, r, "concept-demo/profile-external.html", "text/html; charset=utf-8")
		return
	default:
		http.NotFound(w, r)
	}
}

func serveStaticFile(w http.ResponseWriter, r *http.Request, filePath, contentType string) {
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

func feedbackHandler(w http.ResponseWriter, r *http.Request) {
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
	
	validatedFeedback, isValid := validateFeedback(feedback)

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
		
		if _, err := fmt.Fprintf(file, "=== Feedback received %s ===\nPersonas: %s\n%s\n\n",
			submissionTime.Format("2006-01-02 15:04:05"), personaText, validatedFeedback); err != nil {
			log.Printf("Error writing feedback: %v", err)
		}

		log.Printf("Feedback received from %s (personas: %s, length: %d)", r.RemoteAddr, personaText, len(validatedFeedback))
	}

	tmpl := template.Must(template.ParseFiles("index.html"))
	data := PageData{
		Success:   true,
		CSRFToken: generateCSRFToken(),
	}
	tmpl.Execute(w, data)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/feedback", feedbackHandler)
	mux.HandleFunc("/concept-demo", conceptDemoHandler)
	mux.HandleFunc("/concept-demo/", conceptDemoHandler)
	mux.HandleFunc("/concept-demo/app.css", conceptDemoHandler)
	mux.HandleFunc("/concept-demo/profile-internal.html", conceptDemoHandler)
	mux.HandleFunc("/concept-demo/profile-external.html", conceptDemoHandler)

	handler := securityMiddleware(rateLimitMiddleware(mux))

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(server.ListenAndServe())
}
