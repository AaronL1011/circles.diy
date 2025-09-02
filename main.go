package main

import (
	"log"
	"net/http"
	"time"

	"circles.diy/internal/config"
	"circles.diy/internal/handlers"
	"circles.diy/internal/middleware"
	"circles.diy/internal/templates"
	"circles.diy/internal/utils"
)

func main() {
	// Load configuration
	cfg := config.NewConfig()

	// Build CSS on startup
	log.Println("Building CSS...")
	if err := utils.BuildCSS(); err != nil {
		log.Fatalf("Failed to build CSS: %v", err)
	}

	// Initialize templates
	log.Println("Initializing templates...")
	if err := templates.InitTemplates(); err != nil {
		log.Fatalf("Failed to initialize templates: %v", err)
	}

	// Start CSS file watcher in development mode
	if cfg.IsDev {
		log.Println("Starting CSS file watcher...")
		go utils.WatchCSSFiles()
	}

	// Setup routes
	mux := http.NewServeMux()

	// Manifesto & user research routes
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/feedback", handlers.FeedbackHandler)

	// App routes
	mux.HandleFunc("/dashboard", handlers.DashboardHandler)
	mux.HandleFunc("/dashboard/", handlers.DashboardHandler)
	mux.HandleFunc("/profile", handlers.ProfileHandler)
	mux.HandleFunc("/profile/", handlers.ProfileHandler)

	// Static asset routes
	mux.HandleFunc("/static/css/style.css", func(w http.ResponseWriter, r *http.Request) {
		handlers.ServeStaticFile(w, r, "static/css/style.css", "text/css; charset=utf-8")
	})
	mux.HandleFunc("/static/js/htmx.min.js", func(w http.ResponseWriter, r *http.Request) {
		handlers.ServeStaticFile(w, r, "static/js/htmx.min.js", "application/javascript; charset=utf-8")
	})

	// Apply middleware chain
	handler := middleware.Chain(mux, middleware.SecurityMiddleware, middleware.RateLimitMiddleware)

	// Configure server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	log.Printf("Routes available:")
	log.Printf("  / - Original landing page (index.html)")
	log.Printf("  /dashboard - New dashboard with templates + HTMX")
	log.Printf("  /profile - New profile page with templates + HTMX")
	log.Printf("  /concept-demo/ - Original concept demo")
	log.Fatal(server.ListenAndServe())
}
