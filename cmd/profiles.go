package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"one/internal/browser"
)

var profilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "List available browser profiles",
	Long:  `Display all detected browser profiles with their names and associated emails.`,
	RunE:  runProfiles,
}

func init() {
	rootCmd.AddCommand(profilesCmd)
}

func runProfiles(cmd *cobra.Command, args []string) error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
	dimStyle := lipgloss.NewStyle().Faint(true)

	fmt.Println()
	fmt.Println(titleStyle.Render("Available Browser Profiles"))
	fmt.Println()

	// Chrome profiles
	chromeProfiles, err := browser.ListChromeProfiles()
	if err == nil && len(chromeProfiles) > 0 {
		fmt.Println(headerStyle.Render("● Chrome"))
		fmt.Println()
		
		for _, profile := range chromeProfiles {
			fmt.Printf("  %s\n", lipgloss.NewStyle().Bold(true).Render(profile.Name))
			
			if profile.Email != "" {
				fmt.Printf("    Email: %s\n", profile.Email)
			}
			
			fmt.Printf("    %s\n", dimStyle.Render(fmt.Sprintf("Directory: %s", profile.Directory)))
			fmt.Println()
		}
	} else if err != nil {
		fmt.Printf("  %s\n\n", dimStyle.Render("Chrome not found or not installed"))
	}

	// Firefox profiles
	firefoxProfiles, err := browser.ListFirefoxProfiles()
	if err == nil && len(firefoxProfiles) > 0 {
		fmt.Println(headerStyle.Render("● Firefox"))
		fmt.Println()
		
		for _, profile := range firefoxProfiles {
			fmt.Printf("  %s\n", lipgloss.NewStyle().Bold(true).Render(profile.Name))
			fmt.Printf("    %s\n", dimStyle.Render(fmt.Sprintf("Directory: %s", profile.Directory)))
			fmt.Println()
		}
	} else if err != nil {
		fmt.Printf("  %s\n\n", dimStyle.Render("Firefox not found or not installed"))
	}

	// Safari note (macOS only)
	fmt.Println(headerStyle.Render("● Safari"))
	fmt.Println()
	fmt.Printf("  %s\n", dimStyle.Render("Safari doesn't support multiple profiles via command line"))
	fmt.Println()

	// Usage example
	fmt.Println(titleStyle.Render("Usage in Configuration"))
	fmt.Println()
	fmt.Println("  Add to your project config (~/.config/one/projects/project.yml):")
	fmt.Println()
	fmt.Println("  " + lipgloss.NewStyle().Faint(true).Render("browser:"))
	fmt.Println("  " + lipgloss.NewStyle().Faint(true).Render("  type: chrome"))
	
	if len(chromeProfiles) > 0 {
		// Show first non-default profile as example
		exampleProfile := chromeProfiles[0]
		if len(chromeProfiles) > 1 {
			exampleProfile = chromeProfiles[1]
		}
		fmt.Printf("  %s\n", lipgloss.NewStyle().Faint(true).Render(fmt.Sprintf("  profile: \"%s\"", exampleProfile.Directory)))
	} else {
		fmt.Println("  " + lipgloss.NewStyle().Faint(true).Render("  profile: \"Profile 1\""))
	}
	fmt.Println()

	fmt.Println(dimStyle.Render("Note: Use the Directory value (not the name) in your configuration"))
	fmt.Println()

	return nil
}
