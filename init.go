package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// getEnvString returns string from environment variable.
func getEnvString(env string, def string) (val string) {
	val = os.Getenv(env)

	if len(val) == 0 {
		return def
	}

	return
}

// getEnvBool returns boolean from environment variable.
func getEnvBool(env string, def bool) (ret bool) {
	val := os.Getenv(env)

	if len(val) == 0 {
		return def
	}

	ret, err := strconv.ParseBool(val)
	if err != nil {
		log.Fatal(val + " environment variable is not boolean")
	}

	return
}

// getEnvInt returns int from environment variable.
func getEnvInt(env string, def int) (ret int) {
	val := os.Getenv(env)

	if len(val) == 0 {
		return def
	}

	ret, err := strconv.Atoi(val)
	if err != nil {
		log.Fatal(env + " environment variable is not numeric")
	}

	return
}

func getFloat64(env string, def float64) (ret float64) {
	val := os.Getenv(env)

	if len(val) == 0 {
		return def
	}

	ret, err := strconv.ParseFloat(val, 64)
	if err != nil {
		log.Fatal(env + " environment variable is not floating point.")
	}

	return
}

// String format flag value.
func (i *configInvestments) String() string {
	return fmt.Sprint(*i)
}

// Set set flag value.
func (i *configInvestments) Set(value string) error {
	if len(strings.Split(value, ",")) == 3 {
		var (
			err      error
			quantity float64
			price    float64
		)

		inv := strings.Split(value, ",")
		if quantity, err = strconv.ParseFloat(inv[1], 32); err != nil {
			return err
		}
		if price, err = strconv.ParseFloat(inv[2], 32); err != nil {
			return err
		}
		stockwatchConfig.Investments = append(stockwatchConfig.Investments, configInvestment{
			Ticker:   inv[0],
			Quantity: quantity,
			Price:    price,
		})
	}
	return nil
}

// remove duplicate values from slice
func uniqueString(inputString []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range inputString {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// global variables
var (
	stockwatchConfig Configuration
	timeFormat       = "2006-01-02 15:04:05"
)

// init configures the parameters the process needs to run.
func init() {
	var (
		cmdLnInvestments configInvestments
		cmdLnStocks      string
	)

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

	// verify required data was provided
	if stockwatchConfig.IexAPIKey == "" {
		flag.PrintDefaults()
		log.Fatal("IEX API Key was not provided.")
	} else if cmdLnStocks == "" {
		flag.PrintDefaults()
		log.Fatal("No stocks were defined.")
	} else if stockwatchConfig.HTTPPort < 0 && stockwatchConfig.HTTPPort > 65535 {
		flag.PrintDefaults()
		log.Fatal("Invalid port given for http server provided.")
	} else if stockwatchConfig.Mail.Port < 0 && stockwatchConfig.Mail.Port > 65534 {
		flag.PrintDefaults()
		log.Fatal("Invalid port for mail server provided.")
	}

	// split provided stocks
	re := regexp.MustCompile(`(\s+)?,(\s+)?`)
	stockwatchConfig.TrackedStocks = re.Split(cmdLnStocks, -1)

	// add stocks from investments
	if len(stockwatchConfig.Investments) > 0 {
		for _, i := range stockwatchConfig.Investments {
			stockwatchConfig.TrackedStocks = append(stockwatchConfig.TrackedStocks, strings.ToLower(i.Ticker))
		}
	}

	// validate ticker
	re = regexp.MustCompile(`^[a-z0-9]+$`)
	for _, value := range stockwatchConfig.TrackedStocks {
		if !re.Match([]byte(value)) {
			flag.PrintDefaults()
			log.Fatal("Stock ticker " + value + " is not a valid ticker name.")
		}
	}

	// remove duplicate stock tickers
	stockwatchConfig.TrackedStocks = uniqueString(stockwatchConfig.TrackedStocks)

	log.Println("Initialization complete.")
}
