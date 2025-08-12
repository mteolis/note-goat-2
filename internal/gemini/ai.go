package gemini

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"regexp"
	"strconv"
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

	model = client.GenerativeModel("gemini-2.5-flash-lite")
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

func PromptQuietly(prompt string) (*genai.GenerateContentResponse, error) {
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}
	if resp == nil {
		logger.Error("Error: response is nil.")
		log.Println("Error: response is nil.")
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

func WaitForErrors(prompt string) (*genai.GenerateContentResponse, error) {
	maxRetries := 17 // 2^N - 1 == 2^17 - 1 = 131071 seconds == 2184.52 minutes == 36.41 hours
	backoff := 1

	for i := range maxRetries {
		resp, err := PromptQuietly(prompt)
		if err != nil {
			if IsRateLimitError(err) {
				retry := parseRetryAfter(err) + 1

				logger.Warn("Warning prompting Gemini API: Rate limit reached - retry in %d seconds... Error: %+v\n", "retry", retry, "err", err)
				log.Printf("Warning prompting Gemini API: Rate limit reached - retry in %d seconds...\n", retry)

				time.Sleep(time.Duration(retry) * time.Second)

				continue
			} else {
				logger.Warn("Warning prompting Gemini API: retrying with exponential backoffs.", "attempt", i, "maxRetries", maxRetries, "backoff", backoff, "err", err)
				log.Printf("Warning prompting Gemini API: retrying with exponential backoffs. attempt: %d; maxRetries: %d, backoff: %d, err: %+v\n", i, maxRetries, backoff, err)

				time.Sleep(time.Duration(backoff) * time.Second)
				backoff *= 2

				continue
			}
		}

		return resp, nil
	}
	return nil, fmt.Errorf("failed to get response after %d exponential retries.", maxRetries)
}

func parseRetryAfter(err error) int {
	errStr := err.Error()

	re := regexp.MustCompile(`retry in (\d+)s`)
	match := re.FindStringSubmatch(errStr)

	if len(match) == 2 {
		secondsStr := match[1]
		seconds, err := strconv.Atoi(secondsStr)
		if err != nil {
			logger.Error("Error parsing retry after seconds: %+v\n", "err", err)
			log.Printf("Error parsing retry after seconds: %+v\n", err)
			return 0
		}
		return seconds
	}
	return 0
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
