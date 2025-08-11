package utils

import (
	"os/exec"
	"runtime"
	"strings"
)

// GetForegroundApp returns the name of the currently active foreground application
func GetForegroundApp() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return getForegroundAppMacOS()
	case "windows":
		return getForegroundAppWindows()
	case "linux":
		return getForegroundAppLinux()
	default:
		return "Unknown", nil
	}
}

// getForegroundAppMacOS gets the foreground app on macOS
func getForegroundAppMacOS() (string, error) {
	script := `tell application "System Events" to get name of first application process whose frontmost is true`
	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	appName := strings.TrimSpace(string(output))

	// Normalize common app names for consistency
	appName = normalizeAppName(appName)

	return appName, nil
}

// getForegroundAppWindows gets the foreground app on Windows
func getForegroundAppWindows() (string, error) {
	// Windows implementation would use Win32 API
	// For now, return a placeholder
	return "Unknown", nil
}

// getForegroundAppLinux gets the foreground app on Linux
func getForegroundAppLinux() (string, error) {
	// Linux implementation would use X11 or Wayland
	// For now, return a placeholder
	return "Unknown", nil
}

// normalizeAppName standardizes application names for consistent tracking
func normalizeAppName(appName string) string {
	// Map common variations to standard names
	nameMap := map[string]string{
		"Visual Studio Code": "Code Editor",
		"Code":               "Code Editor",
		"Cursor":             "Code Editor",
		"IntelliJ IDEA":      "Code Editor",
		"Xcode":              "Code Editor",
		"sublime_text":       "Code Editor",
		"Terminal":           "Terminal",
		"iTerm2":             "Terminal",
		"iTerm":              "Terminal",
		"kitty":              "Terminal",
		"Google Chrome":      "Browser",
		"Safari":             "Browser",
		"Firefox":            "Browser",
		"Microsoft Edge":     "Browser",
		"Arc":                "Browser",
		"Dia":                "Browser",
		"Slack":              "Slack",
		"Microsoft Teams":    "Video Conference",
		"Zoom":               "Video Conference",
		"zoom.us":            "Video Conference",
		"Mail":               "Email",
		"Outlook":            "Email",
		"Thunderbird":        "Email",
		"Figma":              "Design Tool",
		"Sketch":             "Design Tool",
		"Adobe Photoshop":    "Design Tool",
		"Adobe Illustrator":  "Design Tool",
		"Notion":             "Documentation",
		"Obsidian":           "Documentation",
		"Notes":              "Documentation",
		"Microsoft Word":     "Documentation",
		"Pages":              "Documentation",
		"ChatGPT":            "AI Assistant",
		"Claude":             "AI Assistant",
		"Discord":            "Communication",
		"Telegram":           "Communication",
		"WhatsApp":           "Communication",
		"Messages":           "Communication",
		"Chatwork":           "Communication",
	}

	if normalized, exists := nameMap[appName]; exists {
		return normalized
	}

	// Return original name if no mapping exists
	return appName
}
