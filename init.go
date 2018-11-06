package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/TheSp1der/goerror"
)

// global variables
var (
	cmdLnStocks       string
	cmdLnInvestments  investments
	cmdLnEmailAddress string
	cmdLnEmailHost    string
	cmdLnEmailPort    int
	cmdLnEmailFrom    string
	cmdLnNoConsole    bool
	cmdLnHTTPPort     int

	trackedTickers []string

	timeFormat = "2006-01-02 15:04:05"
)

// getEnvString returns string from environment variable.
func getEnvString(env string, def string) string {
	val := os.Getenv(env)
	if len(val) == 0 {
		return def
	}
	return val
}

// getEnvBool returns boolean from environment variable.
func getEnvBool(env string, def bool) bool {
	var (
		err error
		val = os.Getenv(env)
		ret bool
	)

	if len(val) == 0 {
		return def
	}

	if ret, err = strconv.ParseBool(val); err != nil {
		goerror.Fatal(errors.New(val + " environment variable is not boolean"))
	}

	return ret
}

// getEnvInt returns int from environment variable.
func getEnvInt(env string, def int) int {
	var (
		err error
		val = os.Getenv(env)
		ret int
	)

	if len(val) == 0 {
		return def
	}

	if ret, err = strconv.Atoi(val); err != nil {
		goerror.Fatal(errors.New(env + " environment variable is not numeric"))
	}

	return ret
}

// String format flag value.
func (i *investments) String() string {
	return fmt.Sprint(*i)
}

// Set set flag value.
func (i *investments) Set(value string) error {
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
		cmdLnInvestments = append(cmdLnInvestments, investment{
			Ticker:   inv[0],
			Quantity: quantity,
			Price:    price,
		})
	}
	return nil
}

// init configures the parameters the process needs to run.
func init() {
	// read command line options
	flag.Var(&cmdLnInvestments, "invest", "Formatted investment in the form of \"Ticker,Quantity,Price\".")
	stocks := flag.String("ticker", getEnvString("TICKERS", ""), "(TICKERS)\nComma saperated list of stocks to report.")
	mailAddress := flag.String("mailto", getEnvString("EMAIL_TO", ""), "(EMAIL_TO)\nDestination e-mail address that will receive the end of day summary.")
	mailHost := flag.String("mailhost", getEnvString("EMAIL_HOST", ""), "(EMAIL_HOST)\nE-Mail server host.")
	mailPort := flag.Int("mailport", getEnvInt("EMAIL_PORT", 25), "(EMAIL_PORT)\nE-Mail server port.")
	mailFrom := flag.String("mailfrom", getEnvString("EMAIL_FROM", "noreply@localhost"), "(EMAIL_FROM)\nAddress the message will be sent from.")
	noConsole := flag.Bool("noconsole", getEnvBool("NO_CONSOLE", false), "(NO_CONSOLE)\nDon't display stock data in the console.")
	webPort := flag.Int("webport", getEnvInt("WEB_PORT", 0), "(WEB_PORT)\nWeb server listen port.")
	flag.Parse()

	// set global variables
	cmdLnStocks = strings.ToLower(*stocks)
	cmdLnEmailAddress = *mailAddress
	cmdLnEmailHost = *mailHost
	cmdLnEmailPort = *mailPort
	cmdLnEmailFrom = *mailFrom
	cmdLnNoConsole = *noConsole
	cmdLnHTTPPort = *webPort

	// convert input to struct
	re := regexp.MustCompile(`(\s+)?,(\s+)?`)
	trackedTickers = re.Split(cmdLnStocks, -1)

	// verify stocks were provided
	if len(trackedTickers) == 1 && trackedTickers[0] == "" {
		goerror.Fatal(errors.New("no Stocks defined"))
	} else {
		re = regexp.MustCompile(`^[a-z0-9]+$`)
		for _, value := range trackedTickers {
			if !re.Match([]byte(value)) {
				goerror.Fatal(errors.New("stock ticker format error"))
			}
		}
	}

	// add stocks from investments
	if len(cmdLnInvestments) > 0 {
		for _, i := range cmdLnInvestments {
			trackedTickers = append(trackedTickers, strings.ToLower(i.Ticker))
		}
	}
}
