package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"shien/internal/rpc"
)

// GameCommand handles gamification status display
type GameCommand struct{}

// NewGameCommand creates a new game command
func NewGameCommand() *GameCommand {
	return &GameCommand{}
}

// Name returns the command name
func (c *GameCommand) Name() string {
	return "game"
}

// Description returns the command description
func (c *GameCommand) Description() string {
	return "Show gamification status"
}

// Usage returns the command usage
func (c *GameCommand) Usage() string {
	return "game [--json] [--detail]"
}

// Execute runs the game command
func (c *GameCommand) Execute(client *rpc.Client, args []string) error {
	// Parse flags
	jsonOutput := false
	detailOutput := false
	
	for _, arg := range args {
		switch arg {
		case "--json", "-j":
			jsonOutput = true
		case "--detail", "-d":
			detailOutput = true
		}
	}
	
	if detailOutput {
		return c.showDetailedStatus(client, jsonOutput)
	}
	
	return c.showBasicStatus(client, jsonOutput)
}

func (c *GameCommand) showBasicStatus(client *rpc.Client, jsonOutput bool) error {
	status, err := client.GetGamificationStatus("")
	if err != nil {
		return fmt.Errorf("failed to get gamification status: %w", err)
	}
	
	if jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(status)
	}
	
	// Display basic status
	fmt.Println("📊 Gamification Status")
	fmt.Println("=" + strings.Repeat("=", 40))
	
	// Level and experience
	expPercent := float64(status.Experience) / float64(status.NextLevelExp-status.TotalExp+status.Experience) * 100
	expBar := c.makeProgressBar(int(expPercent), 20)
	fmt.Printf("Level: %d %s %d/%d XP (%.0f%%)\n", 
		status.Level, expBar, status.Experience, 
		status.NextLevelExp-status.TotalExp+status.Experience, expPercent)
	fmt.Printf("Total Experience: %d\n", status.TotalExp)
	
	fmt.Println()
	fmt.Println("⚡ Attributes")
	fmt.Println("-" + strings.Repeat("-", 40))
	
	// Display attributes with visual bars
	c.displayAttribute("Focus", status.Focus, 100)
	c.displayAttribute("Productivity", status.Productivity, 100)
	c.displayAttribute("Creativity", status.Creativity, 100)
	c.displayAttribute("Stamina", status.Stamina, 100)
	c.displayAttribute("Knowledge", status.Knowledge, 100)
	c.displayAttribute("Collaboration", status.Collaboration, 100)
	
	fmt.Println()
	fmt.Printf("Last Updated: %s\n", status.UpdatedAt.Format("2006-01-02 15:04:05"))
	
	return nil
}

func (c *GameCommand) showDetailedStatus(client *rpc.Client, jsonOutput bool) error {
	details, err := client.GetGamificationDetails("")
	if err != nil {
		return fmt.Errorf("failed to get gamification details: %w", err)
	}
	
	if jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(details)
	}
	
	// Show basic status first
	if details.Status != nil {
		status := details.Status
		fmt.Println("📊 Gamification Status (Detailed)")
		fmt.Println("=" + strings.Repeat("=", 40))
		
		// Level and experience
		expPercent := float64(status.Experience) / float64(status.NextLevelExp-status.TotalExp+status.Experience) * 100
		expBar := c.makeProgressBar(int(expPercent), 20)
		fmt.Printf("Level: %d %s %d/%d XP (%.0f%%)\n", 
			status.Level, expBar, status.Experience, 
			status.NextLevelExp-status.TotalExp+status.Experience, expPercent)
		fmt.Printf("Total Experience: %d\n", status.TotalExp)
		
		fmt.Println()
		fmt.Println("⚡ Attributes")
		fmt.Println("-" + strings.Repeat("-", 40))
		
		// Display attributes
		c.displayAttribute("Focus", status.Focus, 100)
		c.displayAttribute("Productivity", status.Productivity, 100)
		c.displayAttribute("Creativity", status.Creativity, 100)
		c.displayAttribute("Stamina", status.Stamina, 100)
		c.displayAttribute("Knowledge", status.Knowledge, 100)
		c.displayAttribute("Collaboration", status.Collaboration, 100)
	}
	
	// Show modifiers if any
	if len(details.Modifiers) > 0 {
		fmt.Println()
		fmt.Println("🎯 Active Modifiers")
		fmt.Println("-" + strings.Repeat("-", 40))
		
		for _, mod := range details.Modifiers {
			sign := "+"
			if mod.Value < 0 {
				sign = ""
			}
			expires := "permanent"
			if mod.ExpiresAt != nil {
				expires = mod.ExpiresAt.Format("15:04:05")
			}
			fmt.Printf("  %s%d %s - %s (expires: %s)\n", 
				sign, mod.Value, mod.Attribute, mod.Reason, expires)
		}
	}
	
	// Show recent app usage
	if len(details.RecentApps) > 0 {
		fmt.Println()
		fmt.Println("💻 Today's App Usage")
		fmt.Println("-" + strings.Repeat("-", 40))
		
		for app, minutes := range details.RecentApps {
			hours := minutes / 60
			mins := minutes % 60
			if hours > 0 {
				fmt.Printf("  %-20s %dh %dm\n", app, hours, mins)
			} else {
				fmt.Printf("  %-20s %dm\n", app, mins)
			}
		}
	}
	
	return nil
}

func (c *GameCommand) displayAttribute(name string, value int, maxDisplay int) {
	// Cap display at maxDisplay for visual consistency, but show actual value
	displayValue := value
	if displayValue > maxDisplay {
		displayValue = maxDisplay
	}
	
	bar := c.makeProgressBar(displayValue, 15)
	fmt.Printf("  %-14s %s %3d\n", name, bar, value)
}

func (c *GameCommand) makeProgressBar(percent int, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	
	filled := width * percent / 100
	empty := width - filled
	
	return "[" + strings.Repeat("█", filled) + strings.Repeat("░", empty) + "]"
}