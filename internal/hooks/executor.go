package hooks

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	
	"one/internal/config"
)

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	warnStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
)

// ExecuteHooks runs a list of hooks
func ExecuteHooks(hooks []config.Hook, stage string) error {
	if len(hooks) == 0 {
		return nil
	}

	fmt.Printf("\n%s Running %s hooks...\n", infoStyle.Render("⚡"), stage)
	fmt.Println()

	for i, hook := range hooks {
		if err := executeHook(hook, i+1, len(hooks)); err != nil {
			if hook.FailOnError {
				return fmt.Errorf("%s hook '%s' failed: %w", stage, hook.Name, err)
			}
			// Just warn if fail_on_error is false
			fmt.Printf("%s Hook '%s' failed but continuing (fail_on_error: false)\n", 
				warnStyle.Render("⚠️"), hook.Name)
			fmt.Println()
		}
	}

	fmt.Printf("%s All %s hooks completed\n\n", successStyle.Render("✓"), stage)
	return nil
}

func executeHook(hook config.Hook, current, total int) error {
	// Display hook info
	fmt.Printf("  [%d/%d] %s\n", current, total, lipgloss.NewStyle().Bold(true).Render(hook.Name))
	
	if hook.Description != "" {
		fmt.Printf("        %s\n", infoStyle.Render(hook.Description))
	}
	
	fmt.Printf("        $ %s\n", lipgloss.NewStyle().Faint(true).Render(hook.Command))
	fmt.Println()

	// Execute the command
	start := time.Now()
	output, err := runCommand(hook.Command)
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("%s\n", errorStyle.Render("  ✗ Failed"))
		if output != "" {
			fmt.Printf("\n%s\n", indentOutput(output))
		}
		fmt.Printf("\n  Duration: %s\n\n", duration.Round(time.Millisecond))
		return err
	}

	fmt.Printf("%s (took %s)\n", successStyle.Render("  ✓ Success"), duration.Round(time.Millisecond))
	
	// Show output if verbose or if hook wants it shown
	if output != "" && shouldShowOutput(output) {
		fmt.Printf("\n%s\n", indentOutput(output))
	}
	
	fmt.Println()
	return nil
}

func runCommand(command string) (string, error) {
	// Use shell to execute command so we can support pipes, etc.
	var cmd *exec.Cmd
	
	if isWindows() {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	// Set working directory to current directory
	cmd.Dir, _ = os.Getwd()

	// Capture both stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	
	// Combine output
	output := stdout.String()
	if stderr.Len() > 0 {
		if output != "" {
			output += "\n"
		}
		output += stderr.String()
	}

	return strings.TrimSpace(output), err
}

func indentOutput(output string) string {
	lines := strings.Split(output, "\n")
	var result strings.Builder
	
	for _, line := range lines {
		if line != "" {
			result.WriteString("    ")
			result.WriteString(lipgloss.NewStyle().Faint(true).Render(line))
			result.WriteString("\n")
		}
	}
	
	return result.String()
}

func shouldShowOutput(output string) bool {
	// Show output if it's short (< 10 lines) or seems important
	lines := strings.Split(output, "\n")
	if len(lines) <= 10 {
		return true
	}
	
	// Check for error indicators
	lowerOutput := strings.ToLower(output)
	keywords := []string{"error", "fail", "warning", "todo", "fixme"}
	for _, keyword := range keywords {
		if strings.Contains(lowerOutput, keyword) {
			return true
		}
	}
	
	return false
}

func isWindows() bool {
	return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}

// ValidateHooks validates hook configuration
func ValidateHooks(hooks []config.Hook) error {
	for i, hook := range hooks {
		if hook.Name == "" {
			return fmt.Errorf("hook %d: name is required", i+1)
		}
		if hook.Command == "" {
			return fmt.Errorf("hook '%s': command is required", hook.Name)
		}
	}
	return nil
}
