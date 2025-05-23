package goat

import (
	"log"
	"log/slog"
	"os"

	"github.com/mteolis/note-goat-2/internal/gemini"
	"github.com/mteolis/note-goat-2/internal/utils"
)

var (
	logger     *slog.Logger
	excelFile  string
	promptFile string
)

func InitGoat(slogger *slog.Logger, excelPath string, promptPath string, geminiApiKey string) {
	logger = slogger
	excelFile = excelPath
	promptFile = promptPath
	gemini.InitModel(logger, geminiApiKey)
}

func AddAISummary() {
	log.Println("Adding AI summary...")
	prompt := utils.ReadFileTrimmed(promptFile)
	response, err := gemini.Prompt(prompt)
	if err != nil {
		logger.Error("Error prompting Gemini AI: %+v\n", "err", err)
		log.Printf("Error prompting Gemini AI: %+v\n", err)
		return
	}

	answer := gemini.ExtractAnswer(response)

	saveStringToFile(answer)
}

func saveStringToFile(str string) {
	err := os.WriteFile("test/output.txt", []byte(str), 0644)
	if err != nil {
		logger.Error("Error writing to file: %+v", "err", err)
		log.Fatalf("Error writing to file: %+v", err)
	}
}
