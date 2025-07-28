package constants

var (
	AppName         = "NoteGoat"
	Version         = "v2.0.0"
	GeminiApiKeyVar = "GEMINI_API_KEY"

	CommsFxToTemplate = TransferMapping{
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
	PurchasesSection = Purchases{
		Date: Transfer{
			From:     "A2",
			To:       "D79",
			Overflow: "",
		},
		AccountNumber: AccountNumber{
			Transfer: Transfer{
				From:     "E",
				To:       "B82:B93",
				Overflow: "D95",
			},
		},
		AccountType: AccountType{
			Transfer: Transfer{
				From:     "G",
				To:       "B82:B93",
				Overflow: "D95",
			},
		},
		Ticker: Ticker{
			Transfer: Transfer{
				From:     "S12:BP12",
				To:       "E82:E93",
				Overflow: "D95",
			},
		},
		AmountPerTicker: AmountPerTicker{
			Transfer: Transfer{
				From:     "S:BP",
				To:       "H82:H93",
				Overflow: "D95",
			},
		},
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

type Purchases struct {
	Date            Transfer
	AccountNumber   AccountNumber
	AccountType     AccountType
	Ticker          Ticker
	AmountPerTicker AmountPerTicker
}

type AmountPerTicker struct {
	Transfer Transfer
}

type Ticker struct {
	Transfer Transfer
}

type AccountType struct {
	Transfer Transfer
}

type AccountNumber struct {
	Transfer Transfer
}

type Transfer struct {
	From     string
	To       string
	Overflow string
}
