package main

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"google.golang.org/genai"
)

func main() {
	apiKey := os.Getenv("GEMINI_APIKEY")
	if apiKey == "" {
		sendNotification("Invalid API key.", NotificationLevelError)
		log.Fatal("Invalid API key. Please set the GEMINI_APIKEY environment variable.")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: apiKey})
	if err != nil {
		sendNotification("Unable to create GenAI client.", NotificationLevelError)
		log.Fatalf("Failed to create GenAI client: %v", err)
	}

	sendNotification("Started", NotificationLevelInfo)

	clipboardText, err := getClipboardText()
	if err != nil {
		sendNotification("Unable to read clipboard.", NotificationLevelError)
		log.Fatalf("Failed to read from clipboard: %v", err)
	}

	correctedText, err := fixGrammar(ctx, client, clipboardText)
	if err != nil {
		sendNotification("Unable to process text.", NotificationLevelError)
		log.Fatalf("Failed to process text: %v", err)
	}

	if err := writeToClipboard(correctedText); err != nil {
		sendNotification("Unable to write text to clipboard.", NotificationLevelError)
		log.Fatalf("Failed to write to clipboard: %v", err)

	}

	summary := correctedText
	if len(correctedText) > 50 {
		summary = correctedText[:50] + "..."
	}
	sendNotification("Finished: "+summary, NotificationLevelInfo)
}

const (
	NotificationLevelError = 2
	NotificationLevelWarn  = 1
	NotificationLevelInfo  = 0
)

func sendNotification(text string, level int) error {
	var title string
	switch level {
	case NotificationLevelError:
		{
			title = "[ERROR]"
		}

	case NotificationLevelWarn:
		{
			title = "[WARNING]"
		}

	case NotificationLevelInfo:
		{
			title = "[INFO]"
		}
	}
	title += " GrammarFixer"

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
	model := "gemini-2.5-flash"
	if strings.Contains(input, "[QUALITY]") {
		model = "gemini-2.5-pro"
		input = strings.ReplaceAll(input, "[QUALITY]", "")
	}

	if strings.Contains(input, "[FAST]") {
		model = "gemini-2.5-flash-lite"
		input = strings.ReplaceAll(input, "[QUALITY]", "")
	}

	re := regexp.MustCompile(`(?i)\[INST=(.*?)\]`)
	matches := re.FindStringSubmatch(input)
	var customInstruction string
	if len(matches) > 1 {
		customInstruction = matches[1]
		input = strings.ReplaceAll(input, matches[0], "")
		customInstruction = " " + customInstruction
	}

	re = regexp.MustCompile(`(?i)\[STYLE=(.*?)\]`)
	matches = re.FindStringSubmatch(input)
	var customStyle string
	if len(matches) > 1 {
		customStyle = matches[1]
		input = strings.ReplaceAll(input, matches[0], "")
		customStyle = " Using this style description as a guide: ```" + customStyle + "```"
	}

	resp, err := client.Models.GenerateContent(
		ctx,
		model,
		genai.Text(input),
		&genai.GenerateContentConfig{
			SystemInstruction: &genai.Content{
				Role: "system",
				Parts: []*genai.Part{
					{
						Text: "Fix the grammar and make it sound natural and human. Only send back the corrected version without explanations. Keep the original wording as much as possible. If the original message is rude, keep that tone. Make sure not to use em dashes. Make sure the original message's meaning stays the same. Use the same exact language as the original language was in." + customStyle + customInstruction,
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
