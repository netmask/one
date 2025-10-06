package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"one/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage project configurations",
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured projects",
	RunE:  runConfigList,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current project configuration",
	RunE:  runConfigShow,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configShowCmd)
}

// renderMarkdown renders markdown with Glamour
func renderMarkdown(content string) (string, error) {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100),
	)
	if err != nil {
		return "", err
	}
	
	rendered, err := r.Render(content)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(rendered), nil
}

func runConfigList(cmd *cobra.Command, args []string) error {
	projects, err := config.ListProjects()
	if err != nil {
		return err
	}

	if len(projects) == 0 {
		markdown := `# No Projects Configured

Run ` + "`one init`" + ` to create your first project configuration.

## Quick Start

` + "```bash" + `
cd /path/to/your/project
one init
` + "```" + `
`
		rendered, err := renderMarkdown(markdown)
		if err != nil {
			fmt.Println("No projects configured.")
			fmt.Println("\nRun 'one init' to create your first project configuration.")
			return nil
		}
		fmt.Println(rendered)
		return nil
	}

	// Build markdown list
	var markdown strings.Builder
	markdown.WriteString("# Configured Projects\n\n")
	
	for _, project := range projects {
		markdown.WriteString(fmt.Sprintf("## %s\n\n", project.Project.Name))
		markdown.WriteString(fmt.Sprintf("**Provider**: %s  \n", project.Git.Provider))
		markdown.WriteString("**Paths**:\n")
		for _, path := range project.Project.Paths {
			markdown.WriteString(fmt.Sprintf("- `%s`\n", path))
		}
		markdown.WriteString("\n")
	}

	rendered, err := renderMarkdown(markdown.String())
	if err != nil {
		// Fallback to plain text
		fmt.Println("Configured Projects:\n")
		for _, project := range projects {
			fmt.Printf("  ● %s\n", project.Project.Name)
			fmt.Printf("    Provider: %s\n", project.Git.Provider)
			fmt.Println("    Paths:")
			for _, path := range project.Project.Paths {
				fmt.Printf("      - %s\n", path)
			}
			fmt.Println()
		}
		return nil
	}

	fmt.Println(rendered)
	return nil
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadProjectConfig()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	// Create markdown formatted output with syntax highlighting
	markdown := fmt.Sprintf("# Current Project Configuration\n\n**Project**: %s\n\n## Configuration\n\n```yaml\n%s```\n", 
		cfg.Project.Name, string(data))

	// Try to render with Glow
	rendered, err := renderMarkdown(markdown)
	if err != nil {
		// Fallback to plain text
		fmt.Printf("Current Project Configuration:\n\n")
		fmt.Println(string(data))
		return nil
	}

	fmt.Println(rendered)
	return nil
}
