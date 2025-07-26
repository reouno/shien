package ui

import (
	"fmt"
	"strings"
	"time"
)

type Display struct{}

func NewDisplay() *Display {
	return &Display{}
}

// ShowBanner displays a formatted banner message
func (d *Display) ShowBanner(title, message string) {
	width := 60
	border := strings.Repeat("=", width)
	
	fmt.Println(border)
	fmt.Printf("  %s\n", title)
	fmt.Println(border)
	fmt.Printf("  %s\n", message)
	fmt.Printf("  Time: %s\n", time.Now().Format("15:04:05"))
	fmt.Println(border)
	fmt.Println()
}

// ShowAlert displays an important message with visual emphasis
func (d *Display) ShowAlert(message string) {
	fmt.Printf("\n⚠️  ALERT: %s\n\n", message)
}

// ShowInfo displays an informational message
func (d *Display) ShowInfo(message string) {
	fmt.Printf("ℹ️  %s\n", message)
}

// ShowSuccess displays a success message
func (d *Display) ShowSuccess(message string) {
	fmt.Printf("✅ %s\n", message)
}

// ShowError displays an error message
func (d *Display) ShowError(message string) {
	fmt.Printf("❌ ERROR: %s\n", message)
}