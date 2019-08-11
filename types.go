package main

// Configuration is a struct to hold the runtime configuration.
type Configuration struct {
	TrackedStocks []string
	Investments   configInvestments
	Mail          configMail
	NoConsole     bool
	HTTPPort      int
	IexAPIKey     string
	PollFrequency int
}

type configInvestments []configInvestment
type configInvestment struct {
	Ticker   string
	Quantity float64
	Price    float64
}

type configMail struct {
	Address string
	Host    string
	Port    int
	From    string
}

// outputStructure is the output structure available to templates
type outputStructure struct {
	CurrentTime   string
	MarketStatus  string
	TotalGainLoss string
	Stock         stockData
}

type stockData struct {
	CompanyName iexCompany
	Ask         float64
	Bid         float64
	Last        float64
}

// iexTop struct for the data returned by the iex api.
type iexTop []struct {
	Symbol        string  `json:"symbol"`
	BidSize       float64 `json:"bidSize"`
	BidPrice      float64 `json:"bidPrice"`
	AskSize       float64 `json:"askSize"`
	AskPrice      float64 `json:"askPrice"`
	Volume        float64 `json:"volume"`
	LastSalePrice float64 `json:"lastSalePrice"`
	LastSaleSize  float64 `json:"lastSaleSize"`
	LastSaleTime  int64   `json:"lastSaleTime"`
	LastUpdated   int64   `json:"lastUpdated"`
	Sector        string  `json:"sector"`
	SecurityType  string  `json:"securityType"`
}

type iexCompany struct {
	Symbol       string   `json:"symbol"`
	CompanyName  string   `json:"companyName"`
	Exchange     string   `json:"exchange"`
	Industry     string   `json:"industry"`
	Website      string   `json:"website"`
	Description  string   `json:"description"`
	CEO          string   `json:"CEO"`
	SecurityName string   `json:"securityName"`
	IssueType    string   `json:"issueType"`
	Sector       string   `json:"sector"`
	Employees    int      `json:"employees"`
	Tags         []string `json:"tags"`
}
