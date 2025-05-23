package utils

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/mteolis/note-goat-2/internal/constants"
	"github.com/sqweek/dialog"
)

var (
	logger *slog.Logger
)

func InitUtils(slogger *slog.Logger) {
	logger = slogger
}

func SqweekInputExcel() string {
	log.Printf("Waiting for excel file input selection...\n")

	var excelExtensions = []string{"xlsx", "xls", "xlsm", "xlsb", "csv", "xltx", "xltm", "ods"}
	filePath, err := dialog.File().
		Filter("Excel Files", excelExtensions...).
		Title(fmt.Sprintf("Select an Excel file to %s", constants.AppName)).
		Load()
	if err != nil {
		logger.Error("Error selecting file: %+v\n", "err", err)
		log.Fatalf("Error selecting file: %+v\n", err)
	}
	return filePath
}

func SqweekInputPrompt() string {
	log.Printf("Waiting for prompt text file input selection...\n")

	filePath, err := dialog.File().
		Filter("Text Files", "txt").
		Title(fmt.Sprintf("Select a Text file as %s prompt", constants.AppName)).
		Load()
	if err != nil {
		logger.Error("Error selecting file: %+v\n", "err", err)
		log.Fatalf("Error selecting file: %+v\n", err)
	}
	return filePath
}

func GetGeminiApiKey() string {
	key := os.Getenv(constants.GeminiApiKeyVar)

	if key == "" {
		log.Printf("Waiting for %s text file input selection...\n", constants.GeminiApiKeyVar)

		filePath, err := dialog.File().
			Filter("Text Files", "txt").
			Title(fmt.Sprintf("Select the text file that contains the %s", constants.GeminiApiKeyVar)).
			Load()
		if err != nil {
			logger.Error("Error selecting file: %+v\n", "err", err)
			log.Fatalf("Error selecting file: %+v\n", err)
		}
		key = ReadFileTrimmed(filePath)
	}
	return key
}

func ReadFileTrimmed(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		logger.Error("Error reading file: %+v\n", "err", err)
		log.Fatalf("Error reading file: %+v\n", err)
	}
	return strings.TrimSpace(string(data))
}

func CalcExecutionTime(start time.Time) {
	duration := time.Since(start)

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	logger.Info("%s %s executed successfully in %dh:%dm:%ds\n",
		"App Name", constants.AppName, "Version", constants.Version, "hours", hours, "minutes", minutes, "seconds", seconds)
	log.Printf("%s %s executed successfully in %dh:%dm:%ds\n", constants.AppName, constants.Version, hours, minutes, seconds)
}

func WaitForQuit() {
	if err := keyboard.Open(); err != nil {
		log.Fatalf("Error opening keyboard: %+v", err)
	}
	defer keyboard.Close()

	fmt.Println("Press 'q' to quit...")

	for {
		char, _, err := keyboard.GetKey()
		if err != nil {
			log.Fatalf("Error getting key: %+v", err)
		}
		if char == 'q' || char == 'Q' {
			break
		}
	}

	fmt.Println("Exiting...")
}
