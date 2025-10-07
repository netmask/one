package browser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Profile represents a browser profile with metadata
type Profile struct {
	Name      string // Display name (e.g., "Work Profile")
	Directory string // Internal directory name (e.g., "Profile 1")
	Email     string // Associated email if available
	Type      string // chrome, firefox, safari
}

// ListChromeProfiles returns all Chrome profiles with their metadata
func ListChromeProfiles() ([]Profile, error) {
	chromeDir, err := getChromeDir()
	if err != nil {
		return nil, err
	}

	var profiles []Profile

	// Check Default profile
	defaultPath := filepath.Join(chromeDir, "Default")
	if _, err := os.Stat(defaultPath); err == nil {
		if profile := readChromeProfile(defaultPath, "Default"); profile != nil {
			profiles = append(profiles, *profile)
		}
	}

	// Check numbered profiles
	entries, err := os.ReadDir(chromeDir)
	if err != nil {
		return profiles, nil // Return what we have
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Look for "Profile 1", "Profile 2", etc.
		if strings.HasPrefix(entry.Name(), "Profile ") {
			profilePath := filepath.Join(chromeDir, entry.Name())
			if profile := readChromeProfile(profilePath, entry.Name()); profile != nil {
				profiles = append(profiles, *profile)
			}
		}
	}

	return profiles, nil
}

func readChromeProfile(profilePath, directory string) *Profile {
	prefsPath := filepath.Join(profilePath, "Preferences")
	
	data, err := os.ReadFile(prefsPath)
	if err != nil {
		return nil
	}

	var prefs struct {
		Profile struct {
			Name string `json:"name"`
		} `json:"profile"`
		Account []struct {
			Email    string `json:"email"`
			FullName string `json:"full_name"`
		} `json:"account_info"`
		SigninAllowed bool `json:"signin_allowed"`
	}

	if err := json.Unmarshal(data, &prefs); err != nil {
		return nil
	}

	profile := &Profile{
		Name:      prefs.Profile.Name,
		Directory: directory,
		Type:      "chrome",
	}

	// If no name set, use directory name
	if profile.Name == "" {
		if directory == "Default" {
			profile.Name = "Default"
		} else {
			profile.Name = directory
		}
	}

	// Get email if available
	if len(prefs.Account) > 0 {
		profile.Email = prefs.Account[0].Email
		
		// If no custom name but has account, use full name or email
		if prefs.Profile.Name == "" && prefs.Account[0].FullName != "" {
			profile.Name = prefs.Account[0].FullName
		}
	}

	return profile
}

// ListFirefoxProfiles returns all Firefox profiles
func ListFirefoxProfiles() ([]Profile, error) {
	firefoxDir, err := getFirefoxDir()
	if err != nil {
		return nil, err
	}

	profilesIni := filepath.Join(firefoxDir, "profiles.ini")
	data, err := os.ReadFile(profilesIni)
	if err != nil {
		return nil, err
	}

	var profiles []Profile
	var currentName string
	var currentPath string

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "Name=") {
			currentName = strings.TrimPrefix(line, "Name=")
		}

		if strings.HasPrefix(line, "Path=") {
			currentPath = strings.TrimPrefix(line, "Path=")
		}

		// When we hit a new section or end, save the current profile
		if (strings.HasPrefix(line, "[") || line == "") && currentName != "" {
			profiles = append(profiles, Profile{
				Name:      currentName,
				Directory: currentPath,
				Type:      "firefox",
			})
			currentName = ""
			currentPath = ""
		}
	}

	// Don't forget the last profile
	if currentName != "" {
		profiles = append(profiles, Profile{
			Name:      currentName,
			Directory: currentPath,
			Type:      "firefox",
		})
	}

	return profiles, nil
}

// FormatProfile returns a human-readable string for a profile
func FormatProfile(p Profile) string {
	if p.Email != "" {
		return fmt.Sprintf("%s (%s)", p.Name, p.Email)
	}
	return p.Name
}

// FormatProfileWithDirectory includes the directory for disambiguation
func FormatProfileWithDirectory(p Profile) string {
	if p.Email != "" {
		return fmt.Sprintf("%s (%s) [%s]", p.Name, p.Email, p.Directory)
	}
	if p.Name != p.Directory {
		return fmt.Sprintf("%s [%s]", p.Name, p.Directory)
	}
	return p.Name
}

func getChromeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "Google", "Chrome"), nil
	case "linux":
		return filepath.Join(home, ".config", "google-chrome"), nil
	case "windows":
		return filepath.Join(os.Getenv("LOCALAPPDATA"), "Google", "Chrome", "User Data"), nil
	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

func getFirefoxDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "Firefox"), nil
	case "linux":
		return filepath.Join(home, ".mozilla", "firefox"), nil
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), "Mozilla", "Firefox"), nil
	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}
