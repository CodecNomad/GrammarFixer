package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"GrammarFixer/main/clipboard"
	"GrammarFixer/main/geminiwrapper"
	"GrammarFixer/main/notify"
	"GrammarFixer/main/projects"

	"google.golang.org/genai"
)

func main() {
	// Define command line flags
	var (
		listProjects = flag.Bool("list", false, "List saved projects")
		saveProject  = flag.String("save", "", "Save current clipboard correction as a project with given name")
		help         = flag.Bool("help", false, "Show help information")
	)
	flag.Parse()

	// Show help
	if *help {
		showHelp()
		return
	}

	// Handle project listing
	if *listProjects {
		if err := listProjectsCommand(); err != nil {
			log.Fatalf("Error listing projects: %v", err)
		}
		return
	}

	// Initialize project manager for potential saving
	pm, err := projects.NewProjectManager()
	if err != nil {
		log.Printf("Warning: Unable to initialize project manager: %v", err)
		pm = nil
	}

	// Normal clipboard processing mode
	apiKey := os.Getenv("GEMINI_APIKEY")
	if apiKey == "" {
		notify.SendNotification("Invalid API key.", notify.NotificationLevelError)
		log.Fatal("Invalid API key. Please set the GEMINI_APIKEY environment variable.")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: apiKey})
	if err != nil {
		notify.SendNotification("Unable to create GenAI client.", notify.NotificationLevelError)
		log.Fatalf("Failed to create GenAI client: %v", err)
	}

	notify.SendNotification("Started", notify.NotificationLevelInfo)

	clipboardText, err := clipboard.ReadClipboardText()
	if err != nil {
		notify.SendNotification("Unable to read clipboard.", notify.NotificationLevelError)
		log.Fatalf("Failed to read from clipboard: %v", err)
	}

	correctedText, err := geminiwrapper.FixGrammar(ctx, client, clipboardText)
	if err != nil {
		notify.SendNotification("Unable to process text.", notify.NotificationLevelError)
		log.Fatalf("Failed to process text: %v", err)
	}

	if err := clipboard.WriteToClipboard(correctedText); err != nil {
		notify.SendNotification("Unable to write text to clipboard.", notify.NotificationLevelError)
		log.Fatalf("Failed to write to clipboard: %v", err)
	}

	// Save as project if requested
	if *saveProject != "" && pm != nil {
		if _, err := pm.SaveProject(clipboardText, correctedText, *saveProject); err != nil {
			log.Printf("Warning: Unable to save project: %v", err)
		} else {
			notify.SendNotification("Project saved: "+*saveProject, notify.NotificationLevelInfo)
		}
	}

	summary := correctedText
	if len(correctedText) > 50 {
		summary = correctedText[:50] + "..."
	}
	notify.SendTimedNotification("Finished: "+summary, notify.NotificationLevelInfo, 4500)
}

func showHelp() {
	fmt.Println("GrammarFixer - AI-powered grammar correction tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  GrammarFixer [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -list              List all saved projects")
	fmt.Println("  -save NAME         Save current correction as a project with given name")
	fmt.Println("  -help              Show this help message")
	fmt.Println()
	fmt.Println("Default behavior (no flags):")
	fmt.Println("  Reads text from clipboard, fixes grammar using AI, and writes result back to clipboard")
	fmt.Println()
	fmt.Println("Environment variables:")
	fmt.Println("  GEMINI_APIKEY      Required. Your Gemini AI API key")
}

func listProjectsCommand() error {
	pm, err := projects.NewProjectManager()
	if err != nil {
		return err
	}

	projectList, err := pm.ListProjects()
	if err != nil {
		return err
	}

	if len(projectList) == 0 {
		fmt.Println("No projects found.")
		return nil
	}

	// Sort projects by creation date (newest first)
	sort.Slice(projectList, func(i, j int) bool {
		return projectList[i].CreatedAt.After(projectList[j].CreatedAt)
	})

	fmt.Printf("Found %d project(s):\n\n", len(projectList))

	for i, project := range projectList {
		fmt.Printf("%d. %s\n", i+1, project.Name)
		fmt.Printf("   ID: %s\n", project.ID)
		fmt.Printf("   Created: %s\n", project.CreatedAt.Format("2006-01-02 15:04:05"))

		// Show preview of original text
		original := project.OriginalText
		if len(original) > 100 {
			original = original[:100] + "..."
		}
		fmt.Printf("   Original: %s\n", original)

		// Show preview of corrected text
		corrected := project.CorrectedText
		if len(corrected) > 100 {
			corrected = corrected[:100] + "..."
		}
		fmt.Printf("   Corrected: %s\n", corrected)
		fmt.Println()
	}

	return nil
}
