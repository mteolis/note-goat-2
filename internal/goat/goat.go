package goat

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	// "github.com/mteolis/note-goat-2/internal/constants"
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
	fmt.Printf("Excel Comms File: %s\n", excelCommsFile)

	excelPurchasesFile = excelPurchasesPath
	// excelOutputFile = excelOutputPath
	excelOutputFile = "C:\\Users\\Michael\\git\\note-goat-2\\test\\email-attachments\\original\\V1_Script_Template.xlsx"
	fmt.Printf("Excel Output File: %s\n", excelOutputFile)
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
		// remove row number printing
		// fmt.Printf("row %d: %+v\n", i+1, row)
		fmt.Printf("row %d\n", i+1)
		if i == 0 {
			fmt.Printf("skipping header row\n")
			continue // skip header row
		}
		clientId := row[1]
		// TODO potentially change logic to check forward client id change
		if i == 1 || (i > 0 && clientId != rows[i-1][1]) {
			if i > 1 && clientId != rows[i-1][1] {
				utils.MakeOutputDir()

				fmt.Printf("created output dir\n")

				savePreviousClientExcelFile(outputFile, rows, i)

				fmt.Printf("opening file for %s\n", clientId)
				outputFile, err = excelize.OpenFile(excelOutputFile)
			}
		}

		// TODO - implement logic to write data to file per row

		// test writing to same output file
		oldSheetName := sheets[0]
		clientName := row[5]
		outputFile.SetSheetName(oldSheetName, clientName)

		// commsTransfers := constants.CommsFxToTemplate
		// for j, transfer := range commsTransfers.Transfers {
		// 	fmt.Printf("transfer %d: %s -> %s\n", j, transfer.SrcCol, transfer.DstFromCell)
		// 	for _, srcCol := range transfer.SrcCol {
		// 		// TODO - check if is range dst from dst to, if dst to (last cell) full, append to overflow cell
		// 		sc := fmt.Sprintf("%s%d", srcCol, i+1) // +1 to skip header row
		// 		fmt.Printf("sc: %s\n", sc)
		// 		getCell := fmt.Sprintf("%s%d", srcCol, i+1) // +1 to skip header row
		// 		getCellValue, _ := file.GetCellValue(sheets[0], getCell)
		// 		if getCellValue == "" {
		// 			fmt.Printf("skipping getCellValue is empty for %s\n", getCell)
		// 			continue // skip empty cells
		// 		}
		// 		fmt.Printf("getCell(%s): %s\n", getCell, getCellValue)

		// 		// if dstToCell empty, write to dstFromCell
		// 		if transfer.DstToCell == "" {
		// 			fmt.Printf("writing to %s on sheet %s target cell %s with value %s\n", outputFile.Path, clientName, transfer.DstFromCell, getCellValue)
		// 			outputFile.SetCellValue(clientName, transfer.DstFromCell, getCellValue)
		// 		} else {
		// 			// loop dstFrom to dstTo and add on to empty, if full, append to overflow cell
		// 		}
		// 		// fmt.Printf("utils.ColumnToIndex(%s): %d\n", srcCol, utils.ColumnToIndex(srcCol))
		// 		// scv := row[utils.ColumnToIndex(srcCol)]
		// 		// fmt.Printf("scv: %s\n", scv)
		// 	}
		// 	// valueCol := fmt.Sprintf("%s%d", transfer.SrcCol[0], i)
		// 	// fmt.Printf("valueCol: %s\n", valueCol)
		// 	// outputFile.SetCellValue(clientName, transfer.DstFromCell)
		// }

		// saveCurrentClientExcelFile(outputFile, rows, i)
		// outputFile.Close()
		// os.Exit(0) // TODO - remove test exit
		// fmt.Printf("Exiting script - End of testing\n")

		// currentValue, _ := outputFile.GetCellValue(sheet, "A1")
		// newValue := currentValue + "," + strconv.Itoa(i+1)
		// outputFile.SetCellValue(sheet, "A1", newValue)
		// ^^ test writing to same output file

		// save last client file if last row
		if i == len(rows)-1 {
			saveCurrentClientExcelFile(outputFile, rows, i)
			fmt.Printf("saved file for %s\n", clientId)
		}
	}

	fmt.Printf("Exiting script - End of testing\n")
	os.Exit(0)
}

func saveCurrentClientExcelFile(outputFile *excelize.File, rows [][]string, i int) {
	clientId := rows[i][1]
	clientName := rows[i][1]

	outputFile.SetSheetName(outputFile.GetSheetName(0), clientName)

	if err := outputFile.SaveAs(fmt.Sprintf("output/%s.xlsx", clientId)); err != nil {
		logger.Error("Error saving output file for client %s with id %s: %+v\n", "clientName", clientName, "clientId", clientId, "err", err)
		log.Printf("Error saving output file for client %s with id %s: %+v\n", clientName, clientId, err)
		return
	}
	outputFile.Close()
	fmt.Printf("saved %s\n", outputFile.Path)
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
	// TODO - remove test console print
	fmt.Printf("saved file for %s\n", previousClientId)
}

func AggregateExcelDataPurchases() {
	logger.Debug("Aggregating Excel data from purchases file...")
	log.Println("Aggregating Excel data from purchases file...")
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
