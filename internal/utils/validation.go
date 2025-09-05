package utils

import (
	"html"
	"regexp"
	"strings"
)

func ValidateFeedback(feedback string) (string, bool) {
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