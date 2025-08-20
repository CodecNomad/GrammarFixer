package main

import (
	"context"
	"log"
	"os"

	"GrammarFixer/main/clipboard"
	"GrammarFixer/main/geminiwrapper"
	"GrammarFixer/main/notify"

	"google.golang.org/genai"
)

func main() {
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

	summary := correctedText
	if len(correctedText) > 50 {
		summary = correctedText[:50] + "..."
	}
	notify.SendTimedNotification("Finished: "+summary, notify.NotificationLevelInfo, 4500)
}
