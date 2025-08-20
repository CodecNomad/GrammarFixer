// Package geminiwrapper provides a function to fix grammar using the Gemini AI model.
package geminiwrapper

import (
	"context"
	"regexp"
	"strings"

	"google.golang.org/genai"
)

func FixGrammar(ctx context.Context, client *genai.Client, input string) (string, error) {
	model := "gemini-2.5-flash"
	if strings.Contains(input, "[QUALITY]") {
		model = "gemini-2.5-pro"
		input = strings.ReplaceAll(input, "[QUALITY]", "")
	}

	if strings.Contains(input, "[FAST]") {
		model = "gemini-2.5-flash-lite"
		input = strings.ReplaceAll(input, "[FAST]", "")
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
