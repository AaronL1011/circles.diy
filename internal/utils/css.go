package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func BuildCSS() error {
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

func WatchCSSFiles() {
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
			if err := BuildCSS(); err != nil {
				log.Printf("Error rebuilding CSS: %v", err)
			}
		}

		time.Sleep(1 * time.Second)
	}
}