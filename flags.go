package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

func readFlags() (config Configuration, cmdLnInvestment configInvestments, cmdLnStocks string) {
	// read log out options
	flag.StringVar(&config.LogLevel, "log", getEnvString("LOG_LEVEL", "warn"), "(LOG_LEVEL)\nVerbosity of log output.")
	// read stock options & lists
	flag.Var(&cmdLnInvestment, "invest", "Formatted investment in the form of \"Ticker,Quantity,Price\".")
	flag.StringVar(&cmdLnStocks, "ticker", getEnvString("TICKERS", ""), "(TICKERS)\nComma saperated list of stocks to report.")
	// read mail configuration options
	flag.StringVar(&config.Mail.Address, "mail-to", getEnvString("EMAIL_TO", ""), "(EMAIL_TO)\nDestination e-mail address that will receive the end of day summary.")
	flag.StringVar(&config.Mail.Host, "mail-host", getEnvString("EMAIL_HOST", ""), "(EMAIL_HOST)\nE-Mail server host.")
	flag.IntVar(&config.Mail.Port, "mail-port", getEnvInt("EMAIL_PORT", 25), "(EMAIL_PORT)\nE-Mail server port.")
	flag.StringVar(&config.Mail.From, "mail-from", getEnvString("EMAIL_FROM", "noreply@localhost"), "(EMAIL_FROM)\nAddress the message will be sent from.")
	// read output options
	flag.BoolVar(&config.NoConsole, "no-console", getEnvBool("NO_CONSOLE", true), "(NO_CONSOLE)\nDon't display stock data in the console.")
	flag.IntVar(&config.HTTPPort, "web-port", getEnvInt("WEB_PORT", 8080), "(WEB_PORT)\nWeb server listen port.")
	// IEX connection options
	flag.StringVar(&config.IexAPIKey, "api-key", getEnvString("API_KEY", ""), "(API_KEY)\nIEX API Key for data retrieval.")
	flag.IntVar(&config.PollFrequency, "refresh", getEnvInt("REFRESH", 2), "(REFRESH)\nTime in seconds between stock data retrieval.")

	flag.Parse()

	return
}

func displayCurrentConfig(config Configuration, stocks string) {
	log.Infof("Configuration value for %v= %v", "LOG_LEVEL  ", config.LogLevel)
	log.Infof("Configuration value for %v= %v", "TICKERS    ", stocks)
	log.Infof("Configuration value for %v= %v", "EMAIL_TO   ", config.Mail.Address)
	log.Infof("Configuration value for %v= %v", "EMAIL_HOST ", config.Mail.Host)
	log.Infof("Configuration value for %v= %v", "EMAIL_PORT ", config.Mail.Port)
	log.Infof("Configuration value for %v= %v", "EMAIL_FROM ", config.Mail.From)
	log.Infof("Configuration value for %v= %v", "NO_CONSOLE ", config.NoConsole)
	log.Infof("Configuration value for %v= %v", "WEB_PORT   ", config.HTTPPort)
	log.Infof("Configuration value for %v= %v", "API_KEY    ", config.IexAPIKey)
	log.Infof("Configuration value for %v= %v", "REFRESH    ", config.PollFrequency)
}

// String format flag value.
func (i *configInvestments) String() string {
	return fmt.Sprint(*i)
}

// Set set flag value.
func (i *configInvestments) Set(value string) error {
	if len(strings.Split(value, ",")) == 3 {
		inv := strings.Split(value, ",")

		quantity, err := strconv.ParseFloat(inv[1], 32)
		if err != nil {
			return err
		}

		price, err := strconv.ParseFloat(inv[2], 32)
		if err != nil {
			return err
		}

		*i = append(*i, configInvestment{
			Ticker:   inv[0],
			Quantity: quantity,
			Price:    price,
		})
	}

	return nil
}
