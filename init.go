package main

import (
	"flag"
)

// global variables
var (
	stockwatchConfig Configuration
	timeFormat       = "2006-01-02 15:04:05"
	cmdLnInvestments configInvestments
	cmdLnStocks      string
)

// init configures the parameters the process needs to run.
func init() {
	// read command line options
	flag.Var(&cmdLnInvestments, "invest", "Formatted investment in the form of \"Ticker,Quantity,Price\".")
	flag.StringVar(&cmdLnStocks, "ticker", getEnvString("TICKERS", ""), "(TICKERS)\nComma saperated list of stocks to report.")
	flag.StringVar(&stockwatchConfig.Mail.Address, "mail-to", getEnvString("EMAIL_TO", ""), "(EMAIL_TO)\nDestination e-mail address that will receive the end of day summary.")
	flag.StringVar(&stockwatchConfig.Mail.Host, "mail-host", getEnvString("EMAIL_HOST", ""), "(EMAIL_HOST)\nE-Mail server host.")
	flag.IntVar(&stockwatchConfig.Mail.Port, "mail-port", getEnvInt("EMAIL_PORT", 25), "(EMAIL_PORT)\nE-Mail server port.")
	flag.StringVar(&stockwatchConfig.Mail.From, "mail-from", getEnvString("EMAIL_FROM", "noreply@localhost"), "(EMAIL_FROM)\nAddress the message will be sent from.")
	flag.BoolVar(&stockwatchConfig.NoConsole, "no-console", getEnvBool("NO_CONSOLE", false), "(NO_CONSOLE)\nDon't display stock data in the console.")
	flag.IntVar(&stockwatchConfig.HTTPPort, "web-port", getEnvInt("WEB_PORT", 0), "(WEB_PORT)\nWeb server listen port.")
	flag.StringVar(&stockwatchConfig.IexAPIKey, "api-key", getEnvString("API_KEY", ""), "(API_KEY)\nIEX API Key for data retrieval.")
	flag.IntVar(&stockwatchConfig.PollFrequency, "refresh", getEnvInt("REFRESH", 2), "(REFRESH)\nTime in seconds between stock data retrieval.")
	flag.Parse()
}
