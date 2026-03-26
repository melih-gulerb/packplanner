package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"packplanner/internal/domain/pack"
)

const (
	defaultPort      = "8680"
	defaultPackSizes = "250,500,1000,2000,5000"
)

type Config struct {
	Port             string
	DefaultPackSizes []int
	AllowedOrigins   []string
}

// Load reads runtime settings from the environment and falls back to sensible defaults.
func Load() (Config, error) {
	packSizes, err := parsePackSizes(getEnv("PACK_SIZES", defaultPackSizes))
	if err != nil {
		return Config{}, err
	}

	return Config{
		Port:             getEnv("PORT", defaultPort),
		DefaultPackSizes: packSizes,
		AllowedOrigins:   parseAllowedOrigins(getEnv("ALLOWED_ORIGINS", "*")),
	}, nil
}

func parsePackSizes(raw string) ([]int, error) {
	parts := strings.Split(raw, ",")
	sizes := make([]int, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		size, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("parse pack size %q: %w", part, err)
		}

		sizes = append(sizes, size)
	}

	// Reuse the domain normalization rules so config input matches API input behavior.
	return pack.NormalizePackSizes(sizes)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func parseAllowedOrigins(raw string) []string {
	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		origins = append(origins, part)
	}

	return origins
}
