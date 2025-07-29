package gemini

import (
	"context"
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

func WaitForRateLimit(prompt string) (*genai.GenerateContentResponse, error) {
	for {
		resp, err := PromptQuietly(prompt)
		if err != nil {
			if IsRateLimitError(err) {
				sleep := parseRetryAfter(err) + 1

				logger.Warn("Warning prompting Gemini API: Rate limit reached - retry in %d seconds... Error: %+v\n", "sleep", sleep, "err", err)
				log.Printf("Warning prompting Gemini API: Rate limit reached - retry in %d seconds...\n", sleep)

				time.Sleep(time.Duration(sleep) * time.Second)

				continue
			}
			return nil, err
		}

		return resp, nil
	}
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
