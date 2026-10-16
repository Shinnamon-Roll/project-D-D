// Package rules holds the campaign menu and the character roster (roster.go).
// JSON names follow docs/design/protocol.md.
package rules

type Campaign struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"` // "available" | "coming_soon"
}

// Campaigns is the hard-coded campaign menu. No adventure content is bundled.
func Campaigns() []Campaign {
	return []Campaign{
		{ID: "starter-set", Name: "Starter Set: Dragons of Stormwreck Isle", Status: "available"},
		{ID: "salem", Name: "Salem", Status: "coming_soon"},
	}
}
