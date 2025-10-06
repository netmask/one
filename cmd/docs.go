package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "View documentation",
	Long:  `Display beautifully formatted documentation including the full specification.`,
	RunE:  runDocs,
}

func init() {
	rootCmd.AddCommand(docsCmd)
	docsCmd.Flags().BoolP("spec", "s", false, "Show full technical specification")
	docsCmd.Flags().BoolP("examples", "e", false, "Show configuration examples")
}

func runDocs(cmd *cobra.Command, args []string) error {
	showSpec, _ := cmd.Flags().GetBool("spec")
	showExamplesFlag, _ := cmd.Flags().GetBool("examples")

	if showSpec {
		return showSpecification()
	}

	if showExamplesFlag {
		return showExamples()
	}

	// Default: show the help content
	return runHelp(cmd, args)
}

func showSpecification() error {
	// Try to find SPECIFICATION.md
	specPaths := []string{
		"SPECIFICATION.md",
		"../SPECIFICATION.md",
		"../../SPECIFICATION.md",
	}

	var content []byte
	var err error
	for _, path := range specPaths {
		content, err = os.ReadFile(path)
		if err == nil {
			break
		}
	}

	if err != nil {
		return fmt.Errorf("could not find SPECIFICATION.md")
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(120),
	)
	if err != nil {
		// Fallback to plain text
		fmt.Println(string(content))
		return nil
	}

	rendered, err := r.Render(string(content))
	if err != nil {
		fmt.Println(string(content))
		return nil
	}

	fmt.Print(rendered)
	return nil
}

func showExamples() error {
	exampleDirs := []string{
		"examples",
		"../examples",
		"../../examples",
	}

	var examplesDir string
	for _, dir := range exampleDirs {
		if _, err := os.Stat(dir); err == nil {
			examplesDir = dir
			break
		}
	}

	if examplesDir == "" {
		return fmt.Errorf("could not find examples directory")
	}

	// Read all example files
	files, err := os.ReadDir(examplesDir)
	if err != nil {
		return fmt.Errorf("failed to read examples directory: %w", err)
	}

	var markdown string
	markdown += "# Configuration Examples\n\n"
	markdown += "These examples show different configurations for various setups.\n\n"

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".yml" {
			continue
		}

		content, err := os.ReadFile(filepath.Join(examplesDir, file.Name()))
		if err != nil {
			continue
		}

		markdown += fmt.Sprintf("## %s\n\n", file.Name())
		markdown += "```yaml\n"
		markdown += string(content)
		markdown += "\n```\n\n"
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100),
	)
	if err != nil {
		fmt.Println(markdown)
		return nil
	}

	rendered, err := r.Render(markdown)
	if err != nil {
		fmt.Println(markdown)
		return nil
	}

	fmt.Print(rendered)
	return nil
}
