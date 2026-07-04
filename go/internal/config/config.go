// Package config resolves Bing Webmaster and IndexNow credentials from multiple sources.
// Priority order: CLI flags > environment variables > .env file.
package config

import (
	"bufio"
	"log/slog"
	"os"
	"strings"
)

const (
	apiKeyEnvVar   = "BING_WEBMASTER_API_KEY"
	indexNowEnvVar = "BING_INDEXNOW_KEY"
	dotEnvFile     = ".env"
)

// Config holds resolved configuration values.
type Config struct {
	// APIKey is the Bing Webmaster API key.
	APIKey string
	// IndexNowKey is the optional default IndexNow key.
	IndexNowKey string
}

// Resolve returns configuration values loaded from the highest-priority source for each key.
func Resolve(apiKeyFlag string, indexNowKeyFlag string) Config {
	dotenv := loadFromDotEnv()

	return Config{
		APIKey:      resolveValue("bing webmaster api key", apiKeyFlag, apiKeyEnvVar, dotenv),
		IndexNowKey: resolveValue("indexnow key", indexNowKeyFlag, indexNowEnvVar, dotenv),
	}
}

func resolveValue(label string, flagValue string, envVar string, dotenv map[string]string) string {
	if flagValue != "" {
		slog.Debug(label + " loaded from CLI flag")
		return flagValue
	}

	if value := os.Getenv(envVar); value != "" {
		slog.Debug(label+" loaded from environment variable", "envVar", envVar)
		return value
	}

	if value := dotenv[envVar]; value != "" {
		slog.Debug(label+" loaded from .env file", "envVar", envVar)
		return value
	}

	return ""
}

// loadFromDotEnv reads supported keys from a .env file in the current directory.
func loadFromDotEnv() map[string]string {
	f, err := os.Open(dotEnvFile)
	if err != nil {
		return map[string]string{}
	}
	defer func() { _ = f.Close() }()

	values := map[string]string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		if after, ok := strings.CutPrefix(line, apiKeyEnvVar+"="); ok {
			values[apiKeyEnvVar] = strings.Trim(after, `"'`)
			continue
		}
		if after, ok := strings.CutPrefix(line, indexNowEnvVar+"="); ok {
			values[indexNowEnvVar] = strings.Trim(after, `"'`)
		}
	}

	return values
}
