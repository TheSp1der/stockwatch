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
	CompanyData iexCompany
	StockDetail iexStock
}

// iexStock is the json structure returned from IEX for stock data
type iexStock struct {
	LatestPrice           float64 `json:"latestPrice"`
	LatestVolume          float64 `json:"latestVolume"`
	LatestUpdate          float64 `json:"latestUpdate"`
	LatestTime            string  `json:"latestTime"`
	CalculationPrice      string  `json:"calculationPrice"`
	LatestSource          string  `json:"latestSource"`
	Change                float64 `json:"change"`
	ChangePercent         float64 `json:"changePercent"`
	Open                  float64 `json:"open"`
	OpenTime              float64 `json:"openTime"`
	Close                 float64 `json:"close"`
	CloseTime             float64 `json:"closeTime"`
	High                  float64 `json:"high"`
	Low                   float64 `json:"low"`
	ExtendedPrice         float64 `json:"extendedPrice"`
	ExtendedChange        float64 `json:"extendedChange"`
	ExtendedChangePercent float64 `json:"extendedChangePercent"`
	ExtendedPriceTime     float64 `json:"extendedPriceTime"`
	DelayedPrice          float64 `json:"delayedPrice"`
	DelayedPriceTime      float64 `json:"delayedPriceTime"`
	MarketCap             float64 `json:"marketCap"`
	AvgTotalVolume        float64 `json:"avgTotalVolume"`
	Week52High            float64 `json:"week52High"`
	Week52Low             float64 `json:"week52Low"`
	YTDChange             float64 `json:"ytdChange"`
	IEXRealtimePrice      float64 `json:"iexRealtimePrice"`
	IEXRealtimeSize       float64 `json:"iexRealtimeSize"`
	IEXLastUpdated        float64 `json:"iexLastUpdated"`
	IEXMarketPercent      float64 `json:"iexMarketPercent"`
	IEXVolume             float64 `json:"iexVolume"`
	IEXBidPrice           float64 `json:"iexBidPrice"`
	IEXBidSize            float64 `json:"iexBidSize"`
	IEXAskPrice           float64 `json:"iexAskPrice"`
	IEXAskSize            float64 `json:"iexAskSize"`
	Symbol                string  `json:"symbol"`
	CompanyName           string  `json:"companyName"`
	PeRatio               float64 `json:"peRatio"`
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
