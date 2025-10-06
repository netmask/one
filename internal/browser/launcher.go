package browser

import (
	"fmt"
	"os/exec"
	"runtime"
)

// OpenURL opens a URL in the specified browser with an optional profile
func OpenURL(browserType, profile, url string) error {
	switch browserType {
	case "chrome":
		return openChrome(profile, url)
	case "firefox":
		return openFirefox(profile, url)
	case "safari":
		return openSafari(url)
	default:
		return fmt.Errorf("unsupported browser: %s", browserType)
	}
}

func openChrome(profile, url string) error {
	var cmd *exec.Cmd

	args := []string{}
	if profile != "" {
		args = append(args, fmt.Sprintf("--profile-directory=%s", profile))
	}
	args = append(args, url)

	switch runtime.GOOS {
	case "darwin":
		cmdArgs := []string{"-a", "Google Chrome", "--args"}
		cmdArgs = append(cmdArgs, args...)
		cmd = exec.Command("open", cmdArgs...)
	case "linux":
		cmd = exec.Command("google-chrome", args...)
	case "windows":
		cmd = exec.Command("chrome.exe", args...)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}

func openFirefox(profile, url string) error {
	var cmd *exec.Cmd

	args := []string{}
	if profile != "" {
		args = append(args, "-P", profile)
	}
	args = append(args, url)

	switch runtime.GOOS {
	case "darwin":
		cmdArgs := []string{"-a", "Firefox", "--args"}
		cmdArgs = append(cmdArgs, args...)
		cmd = exec.Command("open", cmdArgs...)
	case "linux":
		cmd = exec.Command("firefox", args...)
	case "windows":
		cmd = exec.Command("firefox.exe", args...)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}

func openSafari(url string) error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("Safari is only available on macOS")
	}

	cmd := exec.Command("open", "-a", "Safari", url)
	return cmd.Start()
}
