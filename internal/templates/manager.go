package templates

import (
	"fmt"
	"html/template"
	"log"

	"circles.diy/internal/models"
)

type Templates struct {
	Dashboard       *template.Template
	ProfilePublic   *template.Template
	ProfileInternal *template.Template
	Circles         *template.Template
	Chat            *template.Template
}

var templates *Templates

func InitTemplates() error {
	templates = &Templates{}

	// Define custom template functions
	funcMap := template.FuncMap{
		"sub": func(a, b int) int {
			return a - b
		},
		"slice": func(s interface{}, start, end int) interface{} {
			switch v := s.(type) {
			case []models.MediaItem:
				if start < 0 || end > len(v) || start > end {
					return []models.MediaItem{}
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
	templates.Dashboard = dashboardTemplate

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
	templates.ProfilePublic = profileTemplate

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
	templates.ProfileInternal = profileInternalTemplate

	// Parse circles template
	circlesTemplate := template.New("circles").Funcs(funcMap)
	circlesTemplate, err = circlesTemplate.ParseGlob("templates/layouts/*.html")
	if err != nil {
		return fmt.Errorf("failed to parse layout templates for circles: %v", err)
	}

	circlesTemplate, err = circlesTemplate.ParseGlob("templates/components/*.html")
	if err != nil {
		return fmt.Errorf("failed to parse component templates for circles: %v", err)
	}

	circlesTemplate, err = circlesTemplate.ParseFiles("templates/pages/circles.html")
	if err != nil {
		return fmt.Errorf("failed to parse circles template: %v", err)
	}
	templates.Circles = circlesTemplate

	// Parse chat template
	chatTemplate := template.New("chat").Funcs(funcMap)
	chatTemplate, err = chatTemplate.ParseGlob("templates/layouts/*.html")
	if err != nil {
		return fmt.Errorf("failed to parse layout templates for chat: %v", err)
	}

	chatTemplate, err = chatTemplate.ParseGlob("templates/components/*.html")
	if err != nil {
		return fmt.Errorf("failed to parse component templates for chat: %v", err)
	}

	chatTemplate, err = chatTemplate.ParseFiles("templates/pages/chat.html")
	if err != nil {
		return fmt.Errorf("failed to parse chat template: %v", err)
	}
	templates.Chat = chatTemplate

	log.Println("Templates initialized successfully")
	return nil
}

func GetTemplates() *Templates {
	return templates
}