package template

import (
	"strings"
	"time"
)

// Context holds template variables
type Context map[string]string

// Render replaces template variables with values from the context
func Render(template string, context Context) string {
	result := template

	for key, value := range context {
		placeholder := "{" + key + "}"
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result
}

// BuildTicketURL generates a ticket URL based on the ticket system
func BuildTicketURL(system, baseURL, ticketID string) string {
	baseURL = strings.TrimSuffix(baseURL, "/")

	switch system {
	case "jira":
		return baseURL + "/browse/" + ticketID
	case "linear":
		return baseURL + "/issue/" + ticketID
	case "github":
		return baseURL + "/issues/" + ticketID
	default:
		return ""
	}
}

// GetCurrentDate returns the current date in ISO 8601 format
func GetCurrentDate() string {
	return time.Now().Format("2006-01-02")
}
