package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"
	
	"shien/internal/database/repository"
	"shien/internal/paths"
	"shien/internal/rpc"
	"shien/internal/version"
)

func main() {
	// Parse global flags first
	var dataDir string
	newArgs := []string{os.Args[0]} // Keep program name
	
	i := 1
	for i < len(os.Args) {
		if os.Args[i] == "--data-dir" && i+1 < len(os.Args) {
			dataDir = os.Args[i+1]
			i += 2 // Skip both --data-dir and its value
		} else {
			newArgs = append(newArgs, os.Args[i])
			i++
		}
	}
	os.Args = newArgs
	
	// Set custom data directory if provided
	if dataDir != "" {
		if err := paths.SetDataDir(dataDir); err != nil {
			log.Fatalf("Failed to set data directory: %v", err)
		}
	}
	
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
	
	command := os.Args[1]
	
	// Create RPC client
	client, err := rpc.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	
	switch command {
	case "status":
		handleStatus(client)
	case "activity":
		handleActivity(client, os.Args[2:])
	case "config":
		handleConfig(client, os.Args[2:])
	case "ping":
		handlePing(client)
	case "help", "-h", "--help":
		printUsage()
	case "version", "--version", "-v":
		fmt.Printf("shien version %s\n", version.GetFullVersion())
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: shien [--data-dir <path>] <command> [options]")
	fmt.Println()
	fmt.Println("Global Options:")
	fmt.Println("  --data-dir <path>   Use custom data directory")
	fmt.Println("  --version, -v       Show version information")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  status              Show daemon status")
	fmt.Println("  activity [options]  Show activity logs")
	fmt.Println("    -from <date>      Start date (YYYY-MM-DD)")
	fmt.Println("    -to <date>        End date (YYYY-MM-DD)")
	fmt.Println("    -today            Show today's activity")
	fmt.Println("  config              Show current configuration")
	fmt.Println("  ping               Check if daemon is running")
	fmt.Println("  version             Show version information")
	fmt.Println("  help               Show this help message")
}

func handleStatus(client *rpc.Client) {
	status, err := client.GetStatus()
	if err != nil {
		log.Fatalf("Failed to get status: %v", err)
	}
	
	fmt.Println("Shien Service Status")
	fmt.Println("==================")
	fmt.Printf("Running: %v\n", status.Running)
	fmt.Printf("Started: %s\n", status.StartedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Uptime:  %s\n", time.Since(status.StartedAt).Round(time.Second))
	fmt.Printf("Version: %s\n", status.Version)
}

func handleActivity(client *rpc.Client, args []string) {
	flags := flag.NewFlagSet("activity", flag.ExitOnError)
	from := flags.String("from", "", "Start date (YYYY-MM-DD)")
	to := flags.String("to", "", "End date (YYYY-MM-DD)")
	today := flags.Bool("today", false, "Show today's activity")
	
	flags.Parse(args)
	
	params := make(map[string]interface{})
	
	if *today {
		now := time.Now()
		params["from"] = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Format(time.RFC3339)
		params["to"] = now.Format(time.RFC3339)
	} else {
		if *from != "" {
			t, err := time.Parse("2006-01-02", *from)
			if err != nil {
				log.Fatalf("Invalid from date: %v", err)
			}
			params["from"] = t.Format(time.RFC3339)
		}
		
		if *to != "" {
			t, err := time.Parse("2006-01-02", *to)
			if err != nil {
				log.Fatalf("Invalid to date: %v", err)
			}
			// Set to end of day
			t = t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			params["to"] = t.Format(time.RFC3339)
		}
	}
	
	resp, err := client.Call(rpc.MethodGetActivityLogs, params)
	if err != nil {
		log.Fatalf("Failed to get activity logs: %v", err)
	}
	
	if !resp.Success {
		log.Fatalf("Error: %s", resp.Error)
	}
	
	// Convert response data to activity logs
	data, err := json.Marshal(resp.Data)
	if err != nil {
		log.Fatalf("Failed to parse response: %v", err)
	}
	
	var logs []repository.ActivityLog
	if err := json.Unmarshal(data, &logs); err != nil {
		log.Fatalf("Failed to parse activity logs: %v", err)
	}
	
	if len(logs) == 0 {
		fmt.Println("No activity logs found for the specified period")
		return
	}
	
	fmt.Println("Activity Logs")
	fmt.Println("=============")
	fmt.Printf("Total records: %d (≈ %d minutes)\n\n", len(logs), len(logs)*5)
	
	// Group by hour for display
	hourlyCount := make(map[string]int)
	for _, log := range logs {
		hour := log.RecordedAt.Format("2006-01-02 15:00")
		hourlyCount[hour]++
	}
	
	// Sort hours in chronological order
	var hours []string
	for hour := range hourlyCount {
		hours = append(hours, hour)
	}
	sort.Strings(hours)
	
	fmt.Println("Activity by hour:")
	for _, hour := range hours {
		count := hourlyCount[hour]
		bar := strings.Repeat("█", count)
		fmt.Printf("%s: %s (%d)\n", hour, bar, count*5)
	}
}

func handleConfig(client *rpc.Client, args []string) {
	resp, err := client.Call(rpc.MethodGetConfig, nil)
	if err != nil {
		log.Fatalf("Failed to get config: %v", err)
	}
	
	if !resp.Success {
		log.Fatalf("Error: %s", resp.Error)
	}
	
	// Pretty print config
	data, err := json.MarshalIndent(resp.Data, "", "  ")
	if err != nil {
		log.Fatalf("Failed to format config: %v", err)
	}
	
	fmt.Println("Current Configuration")
	fmt.Println("====================")
	fmt.Println(string(data))
}

func handlePing(client *rpc.Client) {
	if err := client.Ping(); err != nil {
		fmt.Println("❌ Daemon is not running")
		os.Exit(1)
	}
	
	fmt.Println("✅ Daemon is running")
}