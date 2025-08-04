package main

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/exec"

	"google.golang.org/genai"
)

func main() {
	apiKey := os.Getenv("GEMINI_APIKEY")
	if apiKey == "" {
		log.Fatal("Invalid API key. Please set the GEMINI_APIKEY environment variable.")
		sendNotification("[ERROR] GrammarFixer", "Invalid API key.")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: apiKey})
	if err != nil {
		log.Fatalf("Failed to create GenAI client: %v", err)
		sendNotification("[ERROR] GrammarFixer", "Unable to create GenAI client.")
	}

	sendNotification("[INFO] GrammarFixer", "Started fixing grammar")

	clipboardText, err := getClipboardText()
	if err != nil {
		log.Fatalf("Failed to read from clipboard: %v", err)
		sendNotification("[ERROR] GrammarFixer", "Unable to read clipboard.")
	}

	correctedText, err := fixGrammar(ctx, client, clipboardText)
	if err != nil {
		log.Fatalf("Failed to process text: %v", err)
		sendNotification("[ERROR] GrammarFixer", "Unable to process text.")
	}

	if err := writeToClipboard(correctedText); err != nil {
		log.Fatalf("Failed to write to clipboard: %v", err)
		sendNotification("[ERROR] GrammarFixer", "Unable to write text to clipboard.")
	}

	sendNotification("[INFO] GrammarFixer", "Finished fixing grammar")
}

func sendNotification(title string, text string) error {
	cmd := exec.Command("notify-send", "-t", "1500", title, text)
	err := cmd.Run()
	return err
}

func getClipboardText() (string, error) {
	var buf bytes.Buffer
	cmd := exec.Command("wl-paste")
	cmd.Stdout = &buf
	err := cmd.Run()
	return buf.String(), err
}

func writeToClipboard(text string) error {
	cmd := exec.Command("wl-copy")
	cmd.Stdin = bytes.NewBufferString(text)
	return cmd.Run()
}

func fixGrammar(ctx context.Context, client *genai.Client, input string) (string, error) {
	resp, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-lite",
		genai.Text(input),
		&genai.GenerateContentConfig{
			SystemInstruction: &genai.Content{
				Role: "system",
				Parts: []*genai.Part{
					{
						Text: "Fix the grammar and make it sound natural and human. Only send back the corrected version without explanations. Keep the original wording as much as possible. If the original message is rude, keep that tone. Make sure not to use em dashes. Make sure the original message's meaning stays the same. Use the same exact language as the original language was in.",
					},
				},
			},
			SafetySettings: []*genai.SafetySetting{
				{
					Category:  genai.HarmCategoryDangerousContent,
					Threshold: genai.HarmBlockThresholdBlockNone,
				},
				{
					Category:  genai.HarmCategoryHarassment,
					Threshold: genai.HarmBlockThresholdBlockNone,
				},
				{
					Category:  genai.HarmCategoryHateSpeech,
					Threshold: genai.HarmBlockThresholdBlockNone,
				},
				{
					Category:  genai.HarmCategorySexuallyExplicit,
					Threshold: genai.HarmBlockThresholdBlockNone,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Text(), nil
}
