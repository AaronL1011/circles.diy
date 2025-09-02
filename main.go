package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

// Template data structures
type BaseData struct {
	Title     string
	ActiveNav string
	Theme     ThemeSettings
	CSRFToken string
}

type ThemeSettings struct {
	Mode   string `json:"mode"`   // light, dark, system
	Radius string `json:"radius"` // 0, 6, 12, 32
}

type User struct {
	ID       string `json:"id"`
	Handle   string `json:"handle"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
	Banner   string `json:"banner"`
}

type FeedItem struct {
	ID       string    `json:"id"`
	User     User      `json:"user"`
	Content  string    `json:"content"`
	TimeAgo  string    `json:"time_ago"`
	Circle   string    `json:"circle"`
	Image    *MediaItem `json:"image,omitempty"`
	Video    *MediaItem `json:"video,omitempty"`
	Gallery  []MediaItem `json:"gallery,omitempty"`
	CanBuy   bool      `json:"can_buy"`
}

type MediaItem struct {
	URL string `json:"url"`
	Alt string `json:"alt"`
}

type Circle struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Thumbnail   string `json:"thumbnail"`
	MemberCount string `json:"member_count"`
	Active      bool   `json:"active"`
}

type Discussion struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Preview string `json:"preview"`
	Circle  string `json:"circle"`
	TimeAgo string `json:"time_ago"`
}

type Event struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Time  string `json:"time"`
	Day   string `json:"day"`
	Month string `json:"month"`
}

type Ripple struct {
	ID        string `json:"id"`
	User      string `json:"user"`
	Content   string `json:"content"`
	ExpiresIn string `json:"expires_in"`
}

type MarketplaceItem struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Price    string `json:"price"`
	Image    string `json:"image"`
	Location string `json:"location"`
	TimeAgo  string `json:"time_ago"`
}

type ImpactItem struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type DashboardData struct {
	BaseData
	Feed              []FeedItem        `json:"feed"`
	FeedOffset        int              `json:"feed_offset"`
	Circles           []Circle         `json:"circles"`
	Discussions       []Discussion     `json:"discussions"`
	Events            []Event          `json:"events"`
	Ripples           []Ripple         `json:"ripples"`
	MarketplaceItems  []MarketplaceItem `json:"marketplace_items"`
	Impact            []ImpactItem     `json:"impact"`
}

type ProfileStats struct {
	Posts       int `json:"posts"`
	Connections int `json:"connections"`
}

type Profile struct {
	ID          string       `json:"id"`
	Handle      string       `json:"handle"`
	Name        string       `json:"name"`
	Avatar      string       `json:"avatar"`
	Banner      string       `json:"banner"`
	Bio         string       `json:"bio"`
	Stats       ProfileStats `json:"stats"`
	IsConnected bool         `json:"is_connected"`
}

type Post struct {
	ID      string      `json:"id"`
	User    User        `json:"user"`
	Content string      `json:"content"`
	TimeAgo string      `json:"time_ago"`
	Circle  string      `json:"circle"`
	Image   *MediaItem  `json:"image,omitempty"`
	Video   *MediaItem  `json:"video,omitempty"`
	Gallery []MediaItem `json:"gallery,omitempty"`
	CanBuy  bool        `json:"can_buy"`
}

type ProfileData struct {
	BaseData
	Profile     Profile `json:"profile"`
	Posts       []Post  `json:"posts"`
	PostOffset  int     `json:"post_offset"`
	HasMorePosts bool   `json:"has_more_posts"`
	IsOwner     bool    `json:"is_owner"`
}

// Template management
type Templates struct {
	dashboard     *template.Template
	profilePublic *template.Template
}

var templates *Templates

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
		limiter = rate.NewLimiter(rate.Every(time.Second), 10)
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
			w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline'; img-src 'self' https://images.unsplash.com https://unsplash.com https://cdn.britannica.com https://media.tenor.com data:; media-src 'self' https://videos.pexels.com https://www.pexels.com data:; font-src 'self'")
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

// CSS building system
func buildCSS() error {
	cssFiles := []string{
		"static/css/variables.css",
		"static/css/base.css",
		"static/css/layout.css",
		"static/css/components.css",
		"static/css/themes.css",
	}
	
	// Create output file
	outputFile := "static/css/style.css"
	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output CSS file: %v", err)
	}
	defer output.Close()
	
	// Write header comment
	if _, err := output.WriteString("/* Generated CSS - circles.diy */\n\n"); err != nil {
		return fmt.Errorf("failed to write header: %v", err)
	}
	
	// Concatenate all CSS files
	for _, cssFile := range cssFiles {
		if _, err := os.Stat(cssFile); os.IsNotExist(err) {
			log.Printf("CSS file not found: %s, skipping", cssFile)
			continue
		}
		
		file, err := os.Open(cssFile)
		if err != nil {
			return fmt.Errorf("failed to open CSS file %s: %v", cssFile, err)
		}
		
		// Write file marker comment
		if _, err := output.WriteString(fmt.Sprintf("/* === %s === */\n", cssFile)); err != nil {
			file.Close()
			return fmt.Errorf("failed to write file marker: %v", err)
		}
		
		// Copy file content
		if _, err := io.Copy(output, file); err != nil {
			file.Close()
			return fmt.Errorf("failed to copy CSS file %s: %v", cssFile, err)
		}
		
		// Add separator
		if _, err := output.WriteString("\n\n"); err != nil {
			file.Close()
			return fmt.Errorf("failed to write separator: %v", err)
		}
		
		file.Close()
		log.Printf("Added %s to compiled CSS", cssFile)
	}
	
	log.Printf("CSS compiled successfully to %s", outputFile)
	return nil
}

// CSS file watcher for development
func watchCSSFiles() {
	cssDir := "static/css"
	lastModTime := time.Time{}
	
	for {
		hasChanges := false
		
		// Check if any CSS files have changed
		err := filepath.Walk(cssDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			
			if filepath.Ext(path) == ".css" && path != "static/css/style.css" {
				if info.ModTime().After(lastModTime) {
					hasChanges = true
					lastModTime = info.ModTime()
				}
			}
			return nil
		})
		
		if err != nil {
			log.Printf("Error watching CSS files: %v", err)
		} else if hasChanges {
			log.Println("CSS files changed, rebuilding...")
			if err := buildCSS(); err != nil {
				log.Printf("Error rebuilding CSS: %v", err)
			}
		}
		
		time.Sleep(1 * time.Second)
	}
}

// Template initialization
func initTemplates() error {
	templates = &Templates{}
	
	// Define custom template functions
	funcMap := template.FuncMap{
		"sub": func(a, b int) int {
			return a - b
		},
		"slice": func(s []interface{}, start, end int) []interface{} {
			if start < 0 || end > len(s) || start > end {
				return []interface{}{}
			}
			return s[start:end]
		},
		"add": func(a, b int) int {
			return a + b
		},
	}
	
	// Parse dashboard template
	dashboardTemplate := template.New("dashboard").Funcs(funcMap)
	dashboardTemplate, err := dashboardTemplate.ParseGlob("templates/layouts/*.html")
	if err != nil {
		return fmt.Errorf("failed to parse layout templates: %v", err)
	}
	
	dashboardTemplate, err = dashboardTemplate.ParseGlob("templates/components/*.html")
	if err != nil {
		return fmt.Errorf("failed to parse component templates: %v", err)
	}
	
	dashboardTemplate, err = dashboardTemplate.ParseFiles("templates/pages/dashboard.html")
	if err != nil {
		return fmt.Errorf("failed to parse dashboard template: %v", err)
	}
	templates.dashboard = dashboardTemplate
	
	// Parse profile template
	profileTemplate := template.New("profile").Funcs(funcMap)
	profileTemplate, err = profileTemplate.ParseGlob("templates/layouts/*.html")
	if err != nil {
		return fmt.Errorf("failed to parse layout templates for profile: %v", err)
	}
	
	profileTemplate, err = profileTemplate.ParseGlob("templates/components/*.html")
	if err != nil {
		return fmt.Errorf("failed to parse component templates for profile: %v", err)
	}
	
	profileTemplate, err = profileTemplate.ParseFiles("templates/pages/profile-public.html")
	if err != nil {
		return fmt.Errorf("failed to parse profile template: %v", err)
	}
	templates.profilePublic = profileTemplate
	
	log.Println("Templates initialized successfully")
	return nil
}

// Mock data generators for development
func getMockDashboardData() DashboardData {
	return DashboardData{
		BaseData: BaseData{
			Title:     "Dashboard",
			ActiveNav: "dashboard",
			Theme: ThemeSettings{
				Mode:   "system",
				Radius: "0",
			},
			CSRFToken: generateCSRFToken(),
		},
		Feed: []FeedItem{
			{
				ID: "1",
				User: User{
					ID:     "maia",
					Handle: "@maia",
					Avatar: "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=32&h=32&fit=crop&crop=face",
				},
				Content: "Just finished this oak coffee table! Happy to step out of my comfort-zone and share some joinery! This piece is available 💜💸",
				TimeAgo: "2h ago",
				Circle:  "Woodworking",
				Image: &MediaItem{
					URL: "https://images.unsplash.com/photo-1707749522150-e3b1b5f3e079?w=600&h=400&fit=crop&crop=center",
					Alt: "Oak coffee table project",
				},
				CanBuy: true,
			},
			{
				ID: "2",
				User: User{
					ID:     "maia",
					Handle: "@maia",
					Avatar: "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=32&h=32&fit=crop&crop=face",
				},
				Content: "Spending today selecting timber for the next commission. There's something meditative about running your hands along the grain, feeling for the perfect piece that wants to become a dining table.",
				TimeAgo: "1d ago",
				Circle:  "Woodworking",
			},
		},
		FeedOffset: 2,
		Circles: []Circle{
			{ID: "1", Name: "Woodworking", Thumbnail: "https://images.unsplash.com/photo-1416879595882-3373a0480b5b?w=24&h=24&fit=crop&crop=center", MemberCount: "234", Active: true},
			{ID: "2", Name: "Local Artists", Thumbnail: "https://images.unsplash.com/photo-1541961017774-22349e4a1262?w=24&h=24&fit=crop&crop=center", MemberCount: "89"},
		},
		Discussions: []Discussion{
			{ID: "1", Title: "Best wood for beginners?", Preview: "I'm just starting out and wondering what wood types are most forgiving...", Circle: "Woodworking", TimeAgo: "3h ago"},
			{ID: "2", Title: "Summer art market planning", Preview: "Who's planning to set up at the summer markets? Let's coordinate...", Circle: "Local Artists", TimeAgo: "5h ago"},
		},
		Events: []Event{
			{ID: "1", Title: "Workshop Planning", Time: "2:00 PM", Day: "15", Month: "FEB"},
			{ID: "2", Title: "Art Gallery Opening", Time: "7:00 PM", Day: "18", Month: "FEB"},
		},
		Ripples: []Ripple{
			{ID: "1", User: "@alex", Content: "Looking for someone to help move furniture this weekend!", ExpiresIn: "2 days"},
			{ID: "2", User: "@sam", Content: "Free wood scraps available - message me!", ExpiresIn: "5 hours"},
		},
		MarketplaceItems: []MarketplaceItem{
			{ID: "1", Title: "Handcrafted Oak Bookshelf", Price: "$240", Image: "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=150&h=100&fit=crop&crop=center", Location: "Downtown", TimeAgo: "2h ago"},
		},
		Impact: []ImpactItem{
			{Label: "Skills Shared", Value: "12"},
			{Label: "Hours Contributed", Value: "34"},
			{Label: "Connections Made", Value: "8"},
		},
	}
}

func getMockProfileData(handle string, isOwner bool) ProfileData {
	return ProfileData{
		BaseData: BaseData{
			Title:     fmt.Sprintf("%s - Profile", handle),
			ActiveNav: "profile",
			Theme: ThemeSettings{
				Mode:   "system",
				Radius: "0",
			},
			CSRFToken: generateCSRFToken(),
		},
		Profile: Profile{
			ID:     "maia",
			Handle: "@maia",
			Name:   "Maia Makes",
			Avatar: "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=128&h=128&fit=crop&crop=face",
			Banner: "https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=1200&h=300&fit=crop&crop=center",
			Bio:    "Woodworker & furniture maker crafting heirloom pieces from sustainably sourced timber. Teaching traditional joinery techniques and sharing the journey from tree to table.",
			Stats: ProfileStats{
				Posts:       3,
				Connections: 342,
			},
			IsConnected: !isOwner,
		},
		Posts: []Post{
			{
				ID: "1",
				User: User{
					ID:     "maia",
					Handle: "@maia",
					Avatar: "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=32&h=32&fit=crop&crop=face",
				},
				Content: "Just finished this oak coffee table! Happy to step out of my comfort-zone and share some joinery! This piece is available 💜💸",
				TimeAgo: "2h ago",
				Circle:  "Woodworking",
				Image: &MediaItem{
					URL: "https://images.unsplash.com/photo-1707749522150-e3b1b5f3e079?w=600&h=400&fit=crop&crop=center",
					Alt: "Oak coffee table project",
				},
				CanBuy: true,
			},
		},
		PostOffset:   1,
		HasMorePosts: false,
		IsOwner:      isOwner,
	}
}

// New route handlers
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	data := getMockDashboardData()
	
	err := templates.dashboard.ExecuteTemplate(w, "dashboard", data)
	if err != nil {
		log.Printf("Error rendering dashboard template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	// Extract handle from URL path if present, otherwise default to current user
	handle := "@maia" // In a real app, this would come from session/auth
	isOwner := true   // In a real app, this would be determined by comparing with logged-in user
	
	data := getMockProfileData(handle, isOwner)
	
	err := templates.profilePublic.ExecuteTemplate(w, "profile-public", data)
	if err != nil {
		log.Printf("Error rendering profile template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Build CSS on startup
	log.Println("Building CSS...")
	if err := buildCSS(); err != nil {
		log.Fatalf("Failed to build CSS: %v", err)
	}
	
	// Initialize templates
	log.Println("Initializing templates...")
	if err := initTemplates(); err != nil {
		log.Fatalf("Failed to initialize templates: %v", err)
	}
	
	// Start CSS file watcher in development mode
	isDev := os.Getenv("ENV") != "production"
	if isDev {
		log.Println("Starting CSS file watcher...")
		go watchCSSFiles()
	}

	mux := http.NewServeMux()
	
	// IMPORTANT: Keep existing root route unchanged - serves index.html
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/feedback", feedbackHandler)
	
	// New template-based routes
	mux.HandleFunc("/dashboard", dashboardHandler)
	mux.HandleFunc("/profile", profileHandler)
	
	// Existing concept-demo routes (preserved)
	mux.HandleFunc("/concept-demo", conceptDemoHandler)
	mux.HandleFunc("/concept-demo/", conceptDemoHandler)
	mux.HandleFunc("/concept-demo/app.css", conceptDemoHandler)
	mux.HandleFunc("/concept-demo/profile-internal.html", conceptDemoHandler)
	mux.HandleFunc("/concept-demo/profile-external.html", conceptDemoHandler)
	
	// Static asset routes
	mux.HandleFunc("/static/css/style.css", func(w http.ResponseWriter, r *http.Request) {
		serveStaticFile(w, r, "static/css/style.css", "text/css; charset=utf-8")
	})
	mux.HandleFunc("/static/js/htmx.min.js", func(w http.ResponseWriter, r *http.Request) {
		serveStaticFile(w, r, "static/js/htmx.min.js", "application/javascript; charset=utf-8")
	})

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
	log.Printf("Routes available:")
	log.Printf("  / - Original landing page (index.html)")
	log.Printf("  /dashboard - New dashboard with templates + HTMX")
	log.Printf("  /profile - New profile page with templates + HTMX")
	log.Printf("  /concept-demo/ - Original concept demo")
	log.Fatal(server.ListenAndServe())
}
