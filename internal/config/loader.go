package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// GetConfigDir returns the configuration directory path
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// Check XDG_CONFIG_HOME first
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, "one"), nil
	}

	return filepath.Join(home, ".config", "one"), nil
}

// GetProjectsDir returns the projects configuration directory
func GetProjectsDir() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "projects"), nil
}

// LoadProjectConfig discovers and loads the project configuration for the current directory
func LoadProjectConfig() (*ProjectConfig, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	projectsDir, err := GetProjectsDir()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(projectsDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("no projects configured (directory not found: %s)", projectsDir)
	}

	// Read all config files in projects directory
	entries, err := os.ReadDir(projectsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read projects directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := filepath.Ext(entry.Name())
		if ext != ".yml" && ext != ".yaml" {
			continue
		}

		configPath := filepath.Join(projectsDir, entry.Name())
		config, err := parseProjectConfig(configPath)
		if err != nil {
			// Skip invalid configs
			continue
		}

		// Check if current directory matches any project path
		if matchesPath(currentDir, config.Project.Paths) {
			return config, nil
		}
	}

	return nil, fmt.Errorf("no project configuration found for current directory: %s", currentDir)
}

// parseProjectConfig parses a project configuration file
func parseProjectConfig(path string) (*ProjectConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config ProjectConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// matchesPath checks if the current directory matches any of the project paths
func matchesPath(currentDir string, projectPaths []string) bool {
	// Normalize current directory
	currentNorm := normalizePath(currentDir)

	for _, path := range projectPaths {
		pathNorm := normalizePath(path)

		// Check if current directory starts with the project path
		if strings.HasPrefix(currentNorm, pathNorm) {
			return true
		}
	}

	return false
}

// normalizePath normalizes a file system path
func normalizePath(path string) string {
	// Resolve symlinks and get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	// Evaluate symlinks
	evalPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		evalPath = absPath
	}

	// Clean the path
	cleanPath := filepath.Clean(evalPath)

	// Remove trailing slash
	cleanPath = strings.TrimSuffix(cleanPath, string(filepath.Separator))

	return cleanPath
}

// ListProjects returns all configured projects
func ListProjects() ([]*ProjectConfig, error) {
	projectsDir, err := GetProjectsDir()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(projectsDir); os.IsNotExist(err) {
		return []*ProjectConfig{}, nil
	}

	entries, err := os.ReadDir(projectsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read projects directory: %w", err)
	}

	var projects []*ProjectConfig
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := filepath.Ext(entry.Name())
		if ext != ".yml" && ext != ".yaml" {
			continue
		}

		configPath := filepath.Join(projectsDir, entry.Name())
		config, err := parseProjectConfig(configPath)
		if err == nil {
			projects = append(projects, config)
		}
	}

	return projects, nil
}

// SaveProjectConfig saves a project configuration
func SaveProjectConfig(config *ProjectConfig, filename string) error {
	projectsDir, err := GetProjectsDir()
	if err != nil {
		return err
	}

	// Create projects directory if it doesn't exist
	if err := os.MkdirAll(projectsDir, 0700); err != nil {
		return fmt.Errorf("failed to create projects directory: %w", err)
	}

	// Ensure filename has .yml extension
	if !strings.HasSuffix(filename, ".yml") && !strings.HasSuffix(filename, ".yaml") {
		filename = filename + ".yml"
	}

	configPath := filepath.Join(projectsDir, filename)

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
