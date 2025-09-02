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
	ID     string `json:"id"`
	Handle string `json:"handle"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Bio    string `json:"bio"`
	Banner string `json:"banner"`
}

type FeedItem struct {
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
	Feed             []FeedItem        `json:"feed"`
	FeedOffset       int               `json:"feed_offset"`
	Circles          []Circle          `json:"circles"`
	Discussions      []Discussion      `json:"discussions"`
	Events           []Event           `json:"events"`
	Ripples          []Ripple          `json:"ripples"`
	MarketplaceItems []MarketplaceItem `json:"marketplace_items"`
	Impact           []ImpactItem      `json:"impact"`
}

type ProfileStats struct {
	Posts       int `json:"posts"`
	Connections int `json:"connections"`
	Circles     int `json:"circles"`
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
	Stats   *PostStats  `json:"stats,omitempty"`
}

type PostStats struct {
	Replies int `json:"replies"`
	Shares  int `json:"shares"`
	Views   int `json:"views"`
}

type Extension struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

type Analytics struct {
	ProfileViews         int `json:"profile_views"`
	ProfileViewsChange   int `json:"profile_views_change"`
	PostEngagement       int `json:"post_engagement"`
	PostEngagementChange int `json:"post_engagement_change"`
	NewConnections       int `json:"new_connections"`
	NewConnectionsChange int `json:"new_connections_change"`
}

type DraftPost struct {
	ID      string      `json:"id"`
	Content string      `json:"content"`
	Image   *MediaItem  `json:"image,omitempty"`
	Video   *MediaItem  `json:"video,omitempty"`
	Gallery []MediaItem `json:"gallery,omitempty"`
}

type ProfileData struct {
	BaseData
	Profile      Profile     `json:"profile"`
	Posts        []Post      `json:"posts"`
	PostOffset   int         `json:"post_offset"`
	HasMorePosts bool        `json:"has_more_posts"`
	IsOwner      bool        `json:"is_owner"`
	Extensions   []Extension `json:"extensions"`
	Analytics    Analytics   `json:"analytics"`
	Drafts       []DraftPost `json:"drafts"`
	DraftCount   int         `json:"draft_count"`
}

// Template management
type Templates struct {
	dashboard       *template.Template
	profilePublic   *template.Template
	profileInternal *template.Template
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
		// if strings.HasPrefix(r.URL.Path, "/concept-demo") {
		// 	w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline'; img-src 'self' https://images.unsplash.com https://unsplash.com https://cdn.britannica.com https://media.tenor.com data:; media-src 'self' https://videos.pexels.com https://www.pexels.com data:; font-src 'self'")
		// } else {
		// 	w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'none'")
		// }

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
		"slice": func(s interface{}, start, end int) interface{} {
			switch v := s.(type) {
			case []MediaItem:
				if start < 0 || end > len(v) || start > end {
					return []MediaItem{}
				}
				return v[start:end]
			case []interface{}:
				if start < 0 || end > len(v) || start > end {
					return []interface{}{}
				}
				return v[start:end]
			default:
				return s
			}
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

	// Parse profile-internal template
	profileInternalTemplate := template.New("profile-internal").Funcs(funcMap)
	profileInternalTemplate, err = profileInternalTemplate.ParseGlob("templates/layouts/*.html")
	if err != nil {
		return fmt.Errorf("failed to parse layout templates for profile-internal: %v", err)
	}

	profileInternalTemplate, err = profileInternalTemplate.ParseGlob("templates/components/*.html")
	if err != nil {
		return fmt.Errorf("failed to parse component templates for profile-internal: %v", err)
	}

	profileInternalTemplate, err = profileInternalTemplate.ParseFiles("templates/pages/profile-internal.html")
	if err != nil {
		return fmt.Errorf("failed to parse profile-internal template: %v", err)
	}
	templates.profileInternal = profileInternalTemplate

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
					Avatar: "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=128&h=128&fit=crop&crop=face",
				},
				Content: "Just finished this oak coffee table! Happy to step out of my comfort-zone and share some joinery! This piece is available 💜💸",
				TimeAgo: "2m ago",
				Circle:  "Woodworking",
				Image: &MediaItem{
					URL: "https://images.unsplash.com/photo-1707749522150-e3b1b5f3e079?w=900&auto=format&fit=crop&q=60&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8NHx8b2FrJTIwdGFibGV8ZW58MHx8MHx8fDI%3D",
					Alt: "Oak coffee table project",
				},
				CanBuy: true,
			},
			{
				ID: "2",
				User: User{
					ID:     "heathtyler",
					Handle: "@heathtyler",
					Avatar: "https://unsplash.com/photos/_mscfMq3kKI/download?ixid=M3wxMjA3fDB8MXxzZWFyY2h8NXx8cHJvZmlsZSUyMHBpY3xlbnwwfHx8fDE3NTY2MzgyNTd8Mg&force=true&w=256&crop=face",
				},
				Content: "Warming up the barbeque and just got couple cases of the finest bread-water. Keen to see you all... remember 7PM dont be late!",
				TimeAgo: "15m ago",
				Circle:  "The Crop Circle",
				Image: &MediaItem{
					URL: "https://images.unsplash.com/photo-1664463758574-e640a7a998d4?q=80&w=2831&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
					Alt: "BBQ gathering setup",
				},
				CanBuy: false,
			},
			{
				ID: "3",
				User: User{
					ID:     "sara_pcb",
					Handle: "@sara_pcb",
					Avatar: "https://unsplash.com/photos/3TLl_97HNJo/download?ixid=M3wxMjA3fDB8MXxzZWFyY2h8MTV8fHByb2ZpbGUlMjBwaWN8ZW58MHx8fHwxNzU2NjM4MjU3fDI&force=true&w=128&crop=face",
				},
				Content: "New tutorial series starting: \"Arduino for Beginners\". First session covers basic circuits and programming fundamentals.",
				TimeAgo: "1h ago",
				Circle:  "DIY Electronics",
				Image: &MediaItem{
					URL: "https://images.unsplash.com/photo-1581091226825-a6a2a5aee158?w=400&h=250&fit=crop&crop=center",
					Alt: "Arduino tutorial setup",
				},
				CanBuy: false,
			},
			{
				ID: "4",
				User: User{
					ID:     "zucc",
					Handle: "@zucc",
					Avatar: "https://media.tenor.com/y1mYLo66EuoAAAAM/zucky.gifs",
				},
				Content: "circles.diy has changed the game forever. \n\nFuck, I wish I'd thought of that.",
				TimeAgo: "1h ago",
				Circle:  "Communication Software",
			},
		},
		FeedOffset: 4,
		Circles: []Circle{
			{ID: "1", Name: "Woodworking", Thumbnail: "https://images.unsplash.com/photo-1504148455328-c376907d081c?w=24&h=24&fit=crop&crop=center", MemberCount: "8 online"},
			{ID: "2", Name: "The Crop Circle", Thumbnail: "https://cdn.britannica.com/15/152615-050-3CB2AEEA/Crop-circle-England.jpg", MemberCount: "12 members"},
			{ID: "3", Name: "DIY Electronics", Thumbnail: "https://images.unsplash.com/photo-1518611012118-696072aa579a?w=24&h=24&fit=crop&crop=center", MemberCount: "24 members"},
		},
		Discussions: []Discussion{
			{ID: "1", Title: "Sustainable Materials: Where to Source?", Preview: "Looking for suppliers of ethically sourced hardwoods. What are your go-to sources for...", Circle: "Woodworking", TimeAgo: "23 replies • 45m ago"},
			{ID: "2", Title: "Pricing Creative Work: Community Wisdom", Preview: "How do you approach pricing custom commissions? Struggling to find the balance between...", Circle: "Local Artists", TimeAgo: "8 replies • 2h ago"},
		},
		Events: []Event{
			{ID: "1", Title: "Workshop: Intro to Drum & Bass Mixing", Time: "2:00 PM", Day: "31", Month: "Aug"},
			{ID: "2", Title: "Monthly Showcase", Time: "7:00 PM", Day: "02", Month: "Sep"},
		},
		Ripples: []Ripple{
			{ID: "1", User: "@maia", Content: "Quick progress update on the oak table project...", ExpiresIn: "1h left"},
			{ID: "2", User: "@craft_collective", Content: "Tip: pre-drill hardwood to avoid splitting", ExpiresIn: "18h left"},
		},
		MarketplaceItems: []MarketplaceItem{
			{ID: "1", Title: "Cordless Drill", Price: "$85", Image: "https://images.unsplash.com/photo-1504148455328-c376907d081c?w=120&h=80&fit=crop&crop=center", Location: "Sydney, NSW", TimeAgo: "3h ago"},
			{ID: "2", Title: "Ceramic Wheel & Tools", Price: "Trade", Image: "https://unsplash.com/photos/ZSgWcW70cTs/download?ixid=M3wxMjA3fDB8MXxzZWFyY2h8NXx8Y2VyYW1pYyUyMHdoZWVsfGVufDB8fHx8MTc1NjU1NjAxNXwy&force=true&w=640", Location: "Stanwell Park, NSW", TimeAgo: "1d ago"},
			{ID: "3", Title: "Oak Lumber Bundle", Price: "$120", Image: "https://unsplash.com/photos/urjasxHT9Ck/download?ixid=M3wxMjA3fDB8MXxzZWFyY2h8OHx8bHVtYmVyfGVufDB8fHx8MTc1NjU1NjcwNnwy&force=true&w=640", Location: "Granville, NSW", TimeAgo: "2d ago"},
		},
		Impact: []ImpactItem{
			{Label: "Contributions", Value: "47"},
			{Label: "Discussions", Value: "23"},
			{Label: "Circle Tithe", Value: "$5/month"},
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
				Circles:     5,
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
			{
				ID: "2",
				User: User{
					ID:     "maia",
					Handle: "@maia",
					Avatar: "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=32&h=32&fit=crop&crop=face",
				},
				Content: "Spending today selecting timber for the next commission. There's something meditative about running your hands along the grain, feeling for the perfect piece that wants to become a dining table. The wood tells its own story - weather marks, growth patterns, all the years it spent reaching toward light.",
				TimeAgo: "1d ago",
				Circle:  "Woodworking",
			},
			{
				ID: "2",
				User: User{
					ID:     "maia",
					Handle: "@maia",
					Avatar: "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=32&h=32&fit=crop&crop=face",
				},
				Content: "With over a decade in the business, it's nice to still have all my fingers.",
				TimeAgo: "3d ago",
				Circle:  "Woodworking",
			},
		},
		PostOffset:   1,
		HasMorePosts: true,
		IsOwner:      isOwner,
	}
}

func getMockProfileInternalData() ProfileData {
	return ProfileData{
		BaseData: BaseData{
			Title:     "My Profile",
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
				Circles:     5,
			},
			IsConnected: false,
		},
		Extensions: []Extension{
			{ID: "pos", Name: "Point of Sale", Description: "Sell directly from your profile", Enabled: true},
			{ID: "analytics", Name: "Analytics", Description: "Get insights about your posts", Enabled: true},
			{ID: "polls", Name: "Polls & Surveys", Description: "Create community polls", Enabled: false},
			{ID: "booking", Name: "Schedules and Booking", Description: "Manage a schedule for booking a service", Enabled: true},
			{ID: "skills", Name: "Skill Exchange", Description: "Offer and request skills", Enabled: false},
			{ID: "gallery", Name: "Portfolio Gallery", Description: "A customisable and currated showcase", Enabled: true},
		},
		Analytics: Analytics{
			ProfileViews:         1247,
			ProfileViewsChange:   12,
			PostEngagement:       89,
			PostEngagementChange: 5,
			NewConnections:       23,
			NewConnectionsChange: -3,
		},
		Drafts: []DraftPost{
			{
				ID:      "draft1",
				Content: "Workshop tour for anyone thinking of booking some time to create their next masterpiece. Also sharing some essential tools and layout considerations I've learned over 15 years of making.",
				Gallery: []MediaItem{
					{URL: "https://unsplash.com/photos/0CCVIuAjORE/download?ixid=M3wxMjA3fDB8MXxzZWFyY2h8N3x8d29ya3Nob3AlMjB3b29kfGVufDB8fHx8MTc1NjY5MzI3M3ww&force=true&w=640", Alt: "Workshop space"},
					{URL: "https://unsplash.com/photos/cSqDUEBQUAQ/download?ixid=M3wxMjA3fDB8MXxzZWFyY2h8Nnx8d29ya3Nob3AlMjB3b29kfGVufDB8fHx8MTc1NjY5MzI3M3ww&force=true&w=150&h=120&fit=crop&crop=center", Alt: "Tool rack"},
					{URL: "https://unsplash.com/photos/PC9EDk5aDtc/download?ixid=M3wxMjA3fDB8MXxzZWFyY2h8MTR8fHdvcmtzaG9wJTIwd29vZHxlbnwwfHx8fDE3NTY2OTMyNzN8MA&force=true&w=150&h=120&fit=crop&crop=center", Alt: "Wood storage"},
				},
			},
		},
		DraftCount: 1,
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
				Stats: &PostStats{
					Replies: 23,
					Shares:  47,
					Views:   156,
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
				Content: "Spending today selecting timber for the next commission. There's something meditative about running your hands along the grain, feeling for the perfect piece that wants to become a dining table. The wood tells its own story - weather marks, growth patterns, the years it spent reaching toward light.",
				TimeAgo: "1d ago",
				Circle:  "Woodworking",
				Stats: &PostStats{
					Replies: 12,
					Shares:  28,
					Views:   89,
				},
			},
			{
				ID: "3",
				User: User{
					ID:     "maia",
					Handle: "@maia",
					Avatar: "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=32&h=32&fit=crop&crop=face",
				},
				Content: "Traditional style - the backbone of solid furniture. Here's the technique I learned from my mentor, passed down through generations of craftspeople. No shortcuts, just sharp tools and patient hands.",
				TimeAgo: "3d ago",
				Circle:  "Woodworking",
				Video: &MediaItem{
					URL: "https://www.pexels.com/download/video/5972633/",
					Alt: "Video showing old school planing of an uneven timber edge",
				},
				Stats: &PostStats{
					Replies: 34,
					Shares:  67,
				},
			},
		},
		PostOffset:   3,
		HasMorePosts: false,
		IsOwner:      true,
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
	path := r.URL.Path

	// Check if this is the internal profile view (/profile) or external (/profile/:handle)
	if path == "/profile" {
		// Internal profile view (owner's dashboard)
		data := getMockProfileInternalData()

		err := templates.profileInternal.ExecuteTemplate(w, "profile-internal", data)
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
		data := getMockProfileData(handle, false) // isOwner = false for external view

		err := templates.profilePublic.ExecuteTemplate(w, "profile-public", data)
		if err != nil {
			log.Printf("Error rendering profile-public template: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else {
		http.NotFound(w, r)
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

	// Manifesto & user research routes
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/feedback", feedbackHandler)

	// App routes
	mux.HandleFunc("/dashboard", dashboardHandler)
	mux.HandleFunc("/profile", profileHandler)
	mux.HandleFunc("/profile/", profileHandler)

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
