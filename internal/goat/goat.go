package goat

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/mteolis/note-goat-2/internal/constants"
	"github.com/mteolis/note-goat-2/internal/gemini"
	"github.com/mteolis/note-goat-2/internal/utils"
	"github.com/xuri/excelize/v2"
)

var (
	logger             *slog.Logger
	excelCommsFile     string
	excelPurchasesFile string
	excelOutputFile    string
	promptFile         string
)

func InitGoat(slogger *slog.Logger, excelCommsPath string, excelPurchasesPath string, excelOutputPath string, promptPath string, geminiApiKey string) {
	logger = slogger
	// excelCommsFile = excelCommsPath
	// TODO - remove testing paths
	excelCommsFile = "C:\\Users\\Michael\\git\\note-goat-2\\test\\email-attachments\\new\\Reviewed_Section 1 & 2 - Communications and FX.xlsx"
	// fmt.Printf("Excel Comms File: %s\n", excelCommsFile)

	// excelPurchasesFile = excelPurchasesPath
	excelPurchasesFile = "C:\\Users\\Michael\\git\\note-goat-2\\test\\email-attachments\\new\\Section 3 - Note Purchases.xlsx"
	// excelOutputFile = excelOutputPath
	excelOutputFile = "C:\\Users\\Michael\\git\\note-goat-2\\test\\email-attachments\\original\\V1_Script_Template.xlsx"
	// fmt.Printf("Excel Output File: %s\n", excelOutputFile)
	promptFile = promptPath
	gemini.InitModel(logger, geminiApiKey)
}

func AggregateExcelDataComms() {
	logger.Debug("Aggregating Excel data from communications file...")
	log.Println("Aggregating Excel data from communications file...")

	file, err := excelize.OpenFile(excelCommsFile)
	if err != nil {
		logger.Error("Error opening communications file: %+v\n", "err", err)
		log.Printf("Error opening communications file: %+v\n", err)
		return
	}
	defer file.Close()

	sheets := file.GetSheetList()
	if len(sheets) == 0 {
		logger.Error("No sheets found in the Excel file: %s", "file", excelCommsFile)
		log.Printf("No sheets found in the Excel file: %s\n", excelCommsFile)
		return
	}

	rows, err := file.GetRows(sheets[0])
	if err != nil {
		logger.Error("Error getting rows from sheet: %+v\n", "err", err)
		log.Printf("Error getting rows from sheet: %+v\n", err)
		return
	}

	outputFile, err := excelize.OpenFile(excelOutputFile)
	if err != nil {
		logger.Error("Error opening output file: %+v\n", "err", err)
		log.Printf("Error opening output file: %+v\n", err)
		return
	}
	defer outputFile.Close()

	// assume sorted by client (col B)
	for i, row := range rows {
		if i == 0 {
			continue // skip header row
		}
		clientId := row[1]
		if i == 1 || (i > 0 && clientId != rows[i-1][1]) {
			if i > 1 && clientId != rows[i-1][1] {
				utils.MakeOutputDir()

				savePreviousClientExcelFile(outputFile, rows, i)

				outputFile, err = excelize.OpenFile(excelOutputFile)
			}
		}

		commsTransfers := constants.CommsFxToTemplate
		for _, transfer := range commsTransfers.Transfers {
			for _, srcCol := range transfer.SrcCol {
				getCell := fmt.Sprintf("%s%d", srcCol, i+1) // +1 to skip header row
				fileCellValue, _ := file.GetCellValue(sheets[0], getCell)
				if srcCol == "AA" && fileCellValue == "" {
					fileCellValue, _ = file.GetCellValue(sheets[0], fmt.Sprintf("%s%d", "AB", i+1))
					if fileCellValue == "" {
						fileCellValue = "N/A"
					}
				}
				if fileCellValue == "" {
					continue // skip empty cells
				}

				if transfer.DstToCell == "" {
					outputFile.SetCellValue(outputFile.GetSheetName(0), transfer.DstFromCell, fileCellValue)
				} else {
					dstFromLetter, dstFromNumber := parseCell(transfer.DstFromCell)
					_, dstToNumber := parseCell(transfer.DstToCell)
					for row := dstFromNumber; row <= dstToNumber; row++ {
						outputCellValue, _ := outputFile.GetCellValue(outputFile.GetSheetName(0), fmt.Sprintf("%s%d", dstFromLetter, row))
						if srcCol == "G" && outputCellValue != "" {
							if len(strings.Split(outputCellValue, " ")) >= 2 {
								continue
							}
							outputCellValue += fmt.Sprintf(" %s", fileCellValue)
							outputFile.SetCellValue(outputFile.GetSheetName(0), fmt.Sprintf("%s%d", dstFromLetter, row), outputCellValue)
							break
						}
						if outputCellValue == "" {
							outputFile.SetCellValue(outputFile.GetSheetName(0), fmt.Sprintf("%s%d", dstFromLetter, row), fileCellValue)
							break
						} else if row == dstToNumber {
							overflowCellData, _ := outputFile.GetCellValue(outputFile.GetSheetName(0), transfer.OverflowCell)
							overflowCellData += fmt.Sprintf(", %s", fileCellValue)
							outputFile.SetCellValue(outputFile.GetSheetName(0), transfer.OverflowCell, overflowCellData)
							break
						}
					}
				}
			}
		}

		// save last client file if last row
		if i == len(rows)-1 {
			saveCurrentClientExcelFile(outputFile, rows, i)
		}
	}
}

func parseCell(cell string) (string, int) {
	var letters string
	var numbersStr string

	for _, r := range cell {
		if unicode.IsLetter(r) {
			letters += string(r)
		} else if unicode.IsDigit(r) {
			numbersStr += string(r)
		}
	}
	numbers, err := strconv.Atoi(numbersStr)
	if err != nil {
		logger.Error("Error converting cell number to int: %+v\n", "err", err)
		log.Printf("Error converting cell number to int: %+v\n", err)
		numbers = -1
	}
	return letters, numbers
}

func saveCurrentClientExcelFile(outputFile *excelize.File, rows [][]string, i int) {
	clientId := rows[i][1]
	clientName := rows[i][5]

	outputFile.SetSheetName(outputFile.GetSheetName(0), clientName)

	if err := outputFile.SaveAs(fmt.Sprintf("output/%s.xlsx", clientId)); err != nil {
		logger.Error("Error saving output file for client %s with id %s: %+v\n", "clientName", clientName, "clientId", clientId, "err", err)
		log.Printf("Error saving output file for client %s with id %s: %+v\n", clientName, clientId, err)
		return
	}
	outputFile.Close()
}

func savePreviousClientExcelFile(outputFile *excelize.File, rows [][]string, i int) {
	previousClientId := rows[i-1][1]
	previousClientName := rows[i-1][5]

	outputFile.SetSheetName(outputFile.GetSheetName(0), previousClientName)

	if err := outputFile.SaveAs(fmt.Sprintf("output/%s.xlsx", previousClientId)); err != nil {
		logger.Error("Error saving output file for client %s with id %s: %+v\n", "clientName", previousClientName, "clientId", previousClientId, "err", err)
		log.Printf("Error saving output file for client %s with id %s: %+v\n", previousClientName, previousClientId, err)
		return
	}
	outputFile.Close()
}

func AggregateExcelDataPurchases() {
	logger.Debug("Aggregating Excel data from purchases file...")
	log.Println("Aggregating Excel data from purchases file...")

	file, err := excelize.OpenFile(excelPurchasesFile)
	if err != nil {
		logger.Error("Error opening purchases file: %+v\n", "err", err)
		log.Printf("Error opening purchases file: %+v\n", err)
		return
	}

	outputFile, err := excelize.OpenFile(excelOutputFile)
	if err != nil {
		logger.Error("Error opening output file: %+v\n", "err", err)
		log.Printf("Error opening output file: %+v\n", err)
		return
	}
	defer outputFile.Close()

	for sheetIndex, sheetName := range file.GetSheetList() {
		rows, err := file.GetRows(file.GetSheetName(sheetIndex))
		for i := range rows {
			if i < 15 {
				continue // skip header rows
			}
			previousClientId, _ := file.GetCellValue(sheetName, fmt.Sprintf("B%d", i))
			clientId, _ := file.GetCellValue(sheetName, fmt.Sprintf("B%d", i+1))
			clientFileName := fmt.Sprintf("output/%s.xlsx", clientId)
			if _, err := os.Stat(clientFileName); os.IsNotExist(err) {
				logger.Warn("Warning: skipping opening output file %s for client with id %s: %+v\n", "fileName", clientFileName, "clientId", clientId, "err", err)
				log.Printf("Warning: skipping opening output file %s for client with id %s: %+v\n", clientFileName, clientId, err)
				continue // skip if file does not exist
			}
			if previousClientId != clientId {
				date, _ := file.GetCellValue(sheetName, constants.PurchasesSection.Date.From)
				outputFile.SetCellValue(outputFile.GetSheetName(0), constants.PurchasesSection.Date.To, date)
				outputFile.Save()
				outputFile.Close()

				outputFile, err = excelize.OpenFile(clientFileName)
				if err != nil {
					logger.Warn("Warning: skipping opening output file %s for client with id %s: %+v\n", "fileName", clientFileName, "clientId", clientId, "err", err)
					log.Printf("Warning: skipping opening output file %s for client with id %s: %+v\n", clientFileName, clientId, err)
					continue // skip if file does not exist
				}
			}

			colSIndex := colLetterToIndex("S")
			colBPIndex := colLetterToIndex("BP")
			for j := colSIndex; j <= colBPIndex; j++ {
				amountPosition := fmt.Sprintf("%s%d", colIndexToLetter(j), i+1)
				amountValue, _ := file.GetCellValue(sheetName, amountPosition)

				if amountValue == "" {
					continue
				}

				colToRange := strings.Split(constants.PurchasesSection.AccountNumber.Transfer.To, ":")
				_, k := parseCell(colToRange[0])
				_, colToNumber := parseCell(colToRange[1])
				for ; k <= colToNumber; k++ {
					accountTargetCell := fmt.Sprintf("%s%d", "B", k)
					tickerTargetCell := fmt.Sprintf("%s%d", "E", k)
					amountTargetCell := fmt.Sprintf("%s%d", "H", k)

					accountTargetValue, _ := outputFile.GetCellValue(outputFile.GetSheetName(0), accountTargetCell)
					tickerTargetValue, _ := outputFile.GetCellValue(outputFile.GetSheetName(0), tickerTargetCell)
					amountTargetValue, _ := outputFile.GetCellValue(outputFile.GetSheetName(0), amountTargetCell)

					accountNumberPosition := fmt.Sprintf("%s%d", "E", i+1)
					accountTypePosition := fmt.Sprintf("%s%d", "G", i+1)
					tickerPosition := fmt.Sprintf("%s%s", colIndexToLetter(j), "12")

					accountNumberValue, _ := file.GetCellValue(sheetName, accountNumberPosition)
					accountTypeValue, _ := file.GetCellValue(sheetName, accountTypePosition)
					accountValue := fmt.Sprintf("%s %s", accountNumberValue, accountTypeValue)
					tickerValue, _ := file.GetCellValue(sheetName, tickerPosition)

					if accountTargetValue == "" && tickerTargetValue == "" && amountTargetValue == "" {
						outputFile.SetCellValue(outputFile.GetSheetName(0), accountTargetCell, accountValue)
						outputFile.SetCellValue(outputFile.GetSheetName(0), tickerTargetCell, tickerValue)
						outputFile.SetCellValue(outputFile.GetSheetName(0), amountTargetCell, amountValue)
						break
					} else if k == colToNumber && accountTargetValue != "" && tickerTargetValue != "" && amountTargetValue != "" {
						overflowValue, _ := outputFile.GetCellValue(outputFile.GetSheetName(0), constants.PurchasesSection.AccountNumber.Transfer.Overflow)
						overflowValue += fmt.Sprintf("; %s %s %s", accountValue, tickerValue, amountValue)
						outputFile.SetCellValue(outputFile.GetSheetName(0), constants.PurchasesSection.AccountNumber.Transfer.Overflow, overflowValue)
						break
					}
				}
			}
			if i == len(rows)-1 {
				date, _ := file.GetCellValue(sheetName, constants.PurchasesSection.Date.From)
				outputFile.SetCellValue(outputFile.GetSheetName(0), constants.PurchasesSection.Date.To, date)
				outputFile.Save()
				outputFile.Close()
			}
		}
	}
}

func colIndexToLetter(index int) string {
	letters := ""
	for index > 0 {
		index--
		letters = string(rune('A'+index%26)) + letters
		index /= 26
	}
	return letters
}

func colLetterToIndex(col string) int {
	col = strings.ToUpper(col)
	result := 0
	for i := 0; i < len(col); i++ {
		result *= 26
		result += int(col[i]-'A') + 1
	}
	return result
}

func AddAISummary() {
	logger.Debug("Adding AI summary...")
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

// TODO: remove output test function
func saveStringToFile(str string) {
	err := os.WriteFile("test/output.txt", []byte(str), 0644)
	if err != nil {
		logger.Error("Error writing to file: %+v", "err", err)
		log.Fatalf("Error writing to file: %+v", err)
	}
}
