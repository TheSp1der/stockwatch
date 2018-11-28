package main

// httpHeader is a struct for http connections to submit multiple
// headers with requests/gets/posts/etc.
type httpHeader []struct {
	Name  string
	Value string
}

// investments is a struct for tracking investment positions for
// calculations of gains/losses.
type investments []investment
type investment struct {
	Ticker   string
	Quantity float64
	Price    float64
}

// outputStructure is the output structure available to templates
type outputStructure struct {
	CurrentTime   string
	MarketStatus  string
	TotalGainLoss string
	Stock         []stockData
}
type stockData struct {
	CompanyName  string
	CurrentValue string
	Change       string
	GL           string
	Symbol       string
}

// iex is a struct for the data returned by the iex api.
type iex map[string]iexData
type iexData struct {
	Quote   iexQuote   `json:"quote"`
	Price   float64    `json:"price"`
	Company iexCompany `json:"company"`
	Stats   iexStats   `json:"stats"`
	Ohlc    iexOhlc    `json:"ohlc"`
}
type iexQuote struct {
	AvgTotalVolume        int64   `json:"avgTotalVolume"`
	CalculationPrice      string  `json:"calculationPrice"`
	Change                float64 `json:"change"`
	ChangePercent         float64 `json:"changePercent"`
	Close                 float64 `json:"close"`
	CloseTime             int64   `json:"closeTime"`
	CompanyName           string  `json:"companyName"`
	DelayedPrice          float64 `json:"delayedPrice"`
	DelayedPriceTime      int64   `json:"delayedPriceTime"`
	ExtendedChange        float64 `json:"extendedChange"`
	ExtendedChangePercent float64 `json:"extendedChangePercent"`
	ExtendedPrice         float64 `json:"extendedPrice"`
	ExtendedPriceTime     int64   `json:"extendedPriceTime"`
	High                  float64 `json:"high"`
	IexAskPrice           float64 `json:"iexAskPrice"`
	IexAskSize            int64   `json:"iexAskSize"`
	IexBidPrice           float64 `json:"iexBidPrice"`
	IexBidSize            int64   `json:"iexBidSize"`
	IexLastUpdated        int64   `json:"iexLastUpdated"`
	IexMarketPercent      float64 `json:"iexMarketPercent"`
	IexRealtimePrice      float64 `json:"iexRealtimePrice"`
	IexRealtimeSize       int64   `json:"iexRealtimeSize"`
	IexVolume             int64   `json:"iexVolume"`
	LatestPrice           float64 `json:"latestPrice"`
	LatestSource          string  `json:"latestSource"`
	LatestTime            string  `json:"latestTime"`
	LatestUpdate          int64   `json:"latestUpdate"`
	LatestVolume          int64   `json:"latestVolume"`
	Low                   float64 `json:"low"`
	MarketCap             int64   `json:"marketCap"`
	Open                  float64 `json:"open"`
	OpenTime              int64   `json:"openTime"`
	PeRatio               float64 `json:"peRatio"`
	PreviousClose         float64 `json:"previousClose"`
	PrimaryExchange       string  `json:"primaryExchange"`
	Sector                string  `json:"sector"`
	Symbol                string  `json:"symbol"`
	Week52High            float64 `json:"week52High"`
	Week52Low             float64 `json:"week52Low"`
	YtdChange             float64 `json:"ytdChange"`
}
type iexCompany struct {
	Ceo         string   `json:"CEO"`
	CompanyName string   `json:"companyName"`
	Description string   `json:"description"`
	Exchange    string   `json:"exchange"`
	Industry    string   `json:"industry"`
	IssueType   string   `json:"issueType"`
	Sector      string   `json:"sector"`
	Symbol      string   `json:"symbol"`
	Tags        []string `json:"tags"`
	Website     string   `json:"website"`
}

type iexStats struct {
	Ebitda              int64       `json:"EBITDA"`
	EPSSurpriseDollar   interface{} `json:"EPSSurpriseDollar"`
	EPSSurprisePercent  float64     `json:"EPSSurprisePercent"`
	Beta                float64     `json:"beta"`
	Cash                int64       `json:"cash"`
	CompanyName         string      `json:"companyName"`
	ConsensusEPS        float64     `json:"consensusEPS"`
	Day200MovingAvg     float64     `json:"day200MovingAvg"`
	Day30ChangePercent  float64     `json:"day30ChangePercent"`
	Day50MovingAvg      float64     `json:"day50MovingAvg"`
	Day5ChangePercent   float64     `json:"day5ChangePercent"`
	Debt                int64       `json:"debt"`
	DividendRate        float64     `json:"dividendRate"`
	DividendYield       float64     `json:"dividendYield"`
	Float               int64       `json:"float"`
	GrossProfit         int64       `json:"grossProfit"`
	InsiderPercent      interface{} `json:"insiderPercent"`
	InstitutionPercent  float64     `json:"institutionPercent"`
	LatestEPS           float64     `json:"latestEPS"`
	Marketcap           int64       `json:"marketcap"`
	Month1ChangePercent float64     `json:"month1ChangePercent"`
	Month3ChangePercent float64     `json:"month3ChangePercent"`
	Month6ChangePercent float64     `json:"month6ChangePercent"`
	NumberOfEstimates   int64       `json:"numberOfEstimates"`
	PeRatioHigh         float64     `json:"peRatioHigh"`
	PeRatioLow          float64     `json:"peRatioLow"`
	PriceToBook         float64     `json:"priceToBook"`
	PriceToSales        float64     `json:"priceToSales"`
	ProfitMargin        float64     `json:"profitMargin"`
	ReturnOnAssets      float64     `json:"returnOnAssets"`
	ReturnOnCapital     interface{} `json:"returnOnCapital"`
	ReturnOnEquity      float64     `json:"returnOnEquity"`
	Revenue             int64       `json:"revenue"`
	RevenuePerEmployee  int64       `json:"revenuePerEmployee"`
	RevenuePerShare     int64       `json:"revenuePerShare"`
	SharesOutstanding   int64       `json:"sharesOutstanding"`
	ShortInterest       int64       `json:"shortInterest"`
	ShortRatio          float64     `json:"shortRatio"`
	Symbol              string      `json:"symbol"`
	TtmEPS              float64     `json:"ttmEPS"`
	Week52change        float64     `json:"week52change"`
	Week52high          float64     `json:"week52high"`
	Week52low           float64     `json:"week52low"`
	Year1ChangePercent  float64     `json:"year1ChangePercent"`
	Year2ChangePercent  float64     `json:"year2ChangePercent"`
	Year5ChangePercent  float64     `json:"year5ChangePercent"`
	YtdChangePercent    float64     `json:"ytdChangePercent"`
}
type iexOhlc struct {
	Close struct {
		Price float64 `json:"price"`
		Time  int64   `json:"time"`
	} `json:"close"`
	High float64 `json:"high"`
	Low  float64 `json:"low"`
	Open struct {
		Price float64 `json:"price"`
		Time  int64   `json:"time"`
	} `json:"open"`
}
