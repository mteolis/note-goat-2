package gemini

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/googleapis/gax-go/v2/apierror"
	"google.golang.org/api/option"
)

var (
	ctx    context.Context
	model  *genai.GenerativeModel
	logger *slog.Logger
)

func InitModel(slogger *slog.Logger, geminiApiKey string) {
	logger = slogger
	ctx = context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiApiKey))
	if err != nil {
		logger.Error("Error creating client: %+v\n", "err", err)
		log.Fatalf("Error creating client: %+v", err)
	}

	model = client.GenerativeModel("gemini-2.0-flash")
}

func Prompt(prompt string) (*genai.GenerateContentResponse, error) {
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		logger.Error("Error generating content: %+v\n", "err", err)
		log.Printf("Error generating content: %+v\n", err)
		return nil, err
	}
	if resp == nil {
		logger.Error("Response is nil")
		log.Println("Response is nil")
		return nil, nil
	}

	return resp, nil
}

func IsRateLimitError(err error) bool {
	// if err HTTP code is 429 (too many requests) then rate limit is reached
	if apiErr, ok := err.(*apierror.APIError); ok {
		return apiErr.HTTPCode() == 429
	}
	return false
}

func WaitForRateLimit(prompt string) (*genai.GenerateContentResponse, error) {
	sleep := 5
	maxSleep := 60

	for {
		resp, err := Prompt(prompt)
		if err != nil {
			if IsRateLimitError(err) {
				logger.Warn("Error prompting Gemini API: Rate limit reached - waiting %d seconds...\n", "sleep", sleep)
				log.Printf("Error prompting Gemini API: Rate limit reached - waiting %d seconds...\n", sleep)

				time.Sleep(time.Duration(sleep) * time.Second)

				sleep *= 2

				if sleep > maxSleep {
					sleep = maxSleep
				}
				continue
			}
			return nil, err
		}

		return resp, nil
	}
}

func ExtractAnswer(resp *genai.GenerateContentResponse) string {
	if len(resp.Candidates) == 0 {
		return ""
	}
	if len(resp.Candidates[0].Content.Parts) == 0 {
		return ""
	}
	// check if the part is of type genai.Text and return it as a string
	if text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
		return string(text)
	}
	return ""
}
