package config

import (
	"errors"
	"os"
	"strings"
)

// Config holds the application configuration
type Config struct {
	AdminPassword   string
	JWTSecret       string
	KratosAdminURL  string
	KratosPublicURL string
	Port            string
	CORSOrigins     []string
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		return nil, errors.New("ADMIN_PASSWORD environment variable is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET environment variable is required")
	}

	kratosAdminURL := os.Getenv("KRATOS_ADMIN_URL")
	if kratosAdminURL == "" {
		kratosAdminURL = "http://localhost:4434"
	}

	kratosPublicURL := os.Getenv("KRATOS_PUBLIC_URL")
	if kratosPublicURL == "" {
		kratosPublicURL = "http://localhost:4433"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Parse CORS origins from comma-separated list
	corsOrigins := parseCORSOrigins(os.Getenv("CORS_ORIGINS"))

	return &Config{
		AdminPassword:   adminPassword,
		JWTSecret:       jwtSecret,
		KratosAdminURL:  kratosAdminURL,
		KratosPublicURL: kratosPublicURL,
		Port:            port,
		CORSOrigins:     corsOrigins,
	}, nil
}

// parseCORSOrigins parses a comma-separated list of origins
// Returns wildcard "*" if empty (allow all origins)
func parseCORSOrigins(origins string) []string {
	if origins == "" {
		// Default to wildcard (allow all origins)
		return []string{"*"}
	}

	// Split by comma and trim whitespace
	parts := strings.Split(origins, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

