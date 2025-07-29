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
	"github.com/xuri/excelize/v2"
)

var (
	logger *slog.Logger
)

func InitUtils(slogger *slog.Logger) {
	logger = slogger
}

func SqweekExcelComms() string {
	logger.Debug("Waiting for excel communications and FX file selection...\n")
	log.Printf("Waiting for excel communications and FX file selection...\n")

	return sqweekExcel("Select an Excel Communications and FX file to parse data from")
}

func SqweekExcelPurchases() string {
	logger.Debug("Waiting for excel purchases file selection...\n")
	log.Printf("Waiting for excel purchases file selection...\n")

	return sqweekExcel("Select an Excel Purchases file to parse data from")
}

func SqweekOutputExcel() string {
	logger.Debug("Waiting for excel template file output selection...\n")
	log.Printf("Waiting for excel template file output selection...\n")

	return sqweekExcel("Select an Excel Template file to save the output as")
}

func sqweekExcel(title string) string {
	var excelExtensions = []string{"xlsx", "xls", "xlsm", "xlsb", "csv", "xltx", "xltm", "ods"}
	filePath, err := dialog.File().
		Filter("Excel Files", excelExtensions...).
		Title(title).
		Load()
	if err != nil {
		logger.Error("Error selecting file: %+v\n", "err", err)
		log.Fatalf("Error selecting file: %+v\n", err)
	}
	return filePath
}

func SqweekInputPrompt() string {
	logger.Debug("Waiting for prompt text file input selection...\n")
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
		logger.Debug(fmt.Sprintf("Waiting for %s text file input selection...\n", constants.GeminiApiKeyVar))
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

func MakeOutputDir() {
	if err := os.MkdirAll("output", os.ModePerm); err != nil {
		logger.Error("Error creating output directory: %+v\n", "err", err)
		log.Printf("Error creating output directory: %+v\n", err)
		return
	}
}

func ReadExcelFileContents(filePath string) (string, error) {
	xl, err := excelize.OpenFile(filePath)
	if err != nil {
		logger.Error("Error reading file %s: %+v\n", "filePath", filePath, "err", err)
		log.Printf("Error reading file %s: %+v\n", filePath, err)
		return "", err
	}

	contents := ""
	rows, err := xl.GetRows(xl.GetSheetName(0))
	for _, row := range rows {
		if len(row) == 0 {
			continue // Skip empty rows
		}
		contents += strings.Join(row, " ") + "\n"
	}

	return contents, err
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
