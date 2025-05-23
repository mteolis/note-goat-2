package main

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/mteolis/note-goat-2/internal/constants"
	"github.com/mteolis/note-goat-2/internal/goat"
	"github.com/mteolis/note-goat-2/internal/utils"
)

var (
	logger *slog.Logger
)

func main() {
	start := time.Now()
	filename := "logs/" + constants.AppName + "_" + constants.Version + "_" + start.Format("20060102_150405") + ".log"

	err := os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		log.Fatalf("Error creating log directory: %+v", err)
	}

	logFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	handler := slog.NewTextHandler(logFile, nil)
	logger = slog.New(handler)

	logger.Info("%s %s executing...\n", constants.AppName, constants.Version)
	log.Printf("%s %s executing...\n", constants.AppName, constants.Version)

	utils.InitUtils(logger)
	excelPath := utils.SqweekInputExcel()
	promptPath := utils.SqweekInputPrompt()
	geminiApiKey := utils.GetGeminiApiKey()

	goat.InitGoat(logger, excelPath, promptPath, geminiApiKey)
	goat.AddAISummary()

	utils.CalcExecutionTime(start)
	utils.WaitForQuit()
}
