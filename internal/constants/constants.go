package constants

var (
	AppName         = "NoteGoat"
	Version         = "v2.0.0"
	GeminiApiKeyVar = "GEMINI_API_KEY"

	CommsClientNameColumn  = "F4"
	CommsAdvisorNameColumn = "E6"
	CommsFxToTemplate      = TransferMapping{
		Transfers: []SrcToDst{
			// Comms
			{SrcCol: []string{"F"}, DstFromCell: "F4"},
			{SrcCol: []string{"U"}, DstFromCell: "E6"},
			{SrcCol: []string{"S"}, DstFromCell: "E10"},
			{SrcCol: []string{"T"}, DstFromCell: "H10"},
			{SrcCol: []string{"Y"}, DstFromCell: "D31"},
			{SrcCol: []string{"E", "G"}, DstFromCell: "B34", DstToCell: "B45", OverflowCell: "D47"},
			{SrcCol: []string{"W"}, DstFromCell: "E34", DstToCell: "E45", OverflowCell: "D47"},
			{SrcCol: []string{"X"}, DstFromCell: "H34", DstToCell: "H45", OverflowCell: "D47"},
			// FX
			{SrcCol: []string{"AC"}, DstFromCell: "D55"},
			{SrcCol: []string{"E", "G"}, DstFromCell: "B58", DstToCell: "B69", OverflowCell: "D71"},
			{SrcCol: []string{"AA"}, DstFromCell: "E58", DstToCell: "E69", OverflowCell: "D71"},
		},
	}
	FxToTemplate = TransferMapping{
		Transfers: []SrcToDst{
			{SrcCol: []string{"AC"}, DstFromCell: "D55"},
			{SrcCol: []string{"E", "G"}, DstFromCell: "B58", DstToCell: "B69"},
			{SrcCol: []string{"AA", "AB"}, DstFromCell: "E58", DstToCell: "E69"},
		},
	}
	PurchasesToTemplate = TransferMapping{
		Transfers: []SrcToDst{
			{SrcCol: []string{"F"}, DstFromCell: "F4"},
			{SrcCol: []string{"E", "G"}, DstFromCell: "B82", DstToCell: "B93"},
			{SrcCol: []string{"?", "?"}, DstFromCell: "B82", DstToCell: "B93"},
		},
	}
	CommsToTemplateMapping = []map[string]string{
		{"F": "F4"},
		{"U": "E6"},
		{"S": "E10"},
		{"T": "H10"},
		{"Y": "D31"},
		{"E,G": "B34:B45"},
		{"W": "E34:E45"},
		{"X": "H34:H45"},
		{"AC": "D55"},
		{"E,G": "D55"},
	}
	ColumnMappings = map[string]int{
		"A":  0,
		"B":  1,
		"C":  2,
		"D":  3,
		"E":  4,
		"F":  5,
		"G":  6,
		"H":  7,
		"I":  8,
		"J":  9,
		"K":  10,
		"L":  11,
		"M":  12,
		"N":  13,
		"O":  14,
		"P":  15,
		"Q":  16,
		"R":  17,
		"S":  18,
		"T":  19,
		"U":  20,
		"V":  21,
		"W":  22,
		"X":  23,
		"Y":  24,
		"Z":  25,
		"AA": 26,
		"AB": 27,
	}
)

type SrcToDst struct {
	SrcCol       []string
	DstFromCell  string
	DstToCell    string
	OverflowCell string
}

type TransferMapping struct {
	Transfers []SrcToDst
}
