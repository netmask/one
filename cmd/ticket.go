package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"one/internal/browser"
	"one/internal/config"
	"one/internal/template"
)

var ticketCmd = &cobra.Command{
	Use:   "ticket <TICKET-ID>",
	Short: "Open a ticket in the browser",
	Args:  cobra.ExactArgs(1),
	RunE:  runTicket,
}

func init() {
	rootCmd.AddCommand(ticketCmd)
}

func runTicket(cmd *cobra.Command, args []string) error {
	ticketID := args[0]

	// Load config
	cfg, err := config.LoadProjectConfig()
	if err != nil {
		return err
	}

	if cfg.Ticket == nil {
		return fmt.Errorf("no ticket system configured for this project")
	}

	// Generate ticket URL
	url := template.BuildTicketURL(cfg.Ticket.System, cfg.Ticket.BaseURL, ticketID)

	// Open in browser
	if err := browser.OpenURL(cfg.Browser.Type, cfg.Browser.Profile, url); err != nil {
		return fmt.Errorf("failed to open browser: %w", err)
	}

	fmt.Printf("✓ Opening ticket %s in browser\n", ticketID)
	fmt.Printf("  URL: %s\n", url)

	return nil
}
