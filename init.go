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
	cmdLnInvestments  Investments
	cmdLnEmailAddress string
	cmdLnEmailHost    string
	cmdLnEmailPort    int
	cmdLnEmailFrom    string

	trackedTickers []string

	timeFormat = "2006-01-02 15:04:05"
)

// getEnvString returns string from environment variable
func getEnvString(env string, def string) string {
	val := os.Getenv(env)
	if len(val) == 0 {
		return def
	}
	return val
}

// getEnvPort returns int from environment variable
func getEnvInt(env string, def int) int {
	val, _ := strconv.Atoi(os.Getenv(env))
	if len(strconv.Itoa(val)) == 0 {
		return def
	}
	return val
}

// String format flag value
func (i *Investments) String() string {
	return fmt.Sprint(*i)
}

// Set set flag value
func (i *Investments) Set(value string) error {
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
		cmdLnInvestments = append(cmdLnInvestments, Investment{
			Ticker:   inv[0],
			Quantity: float32(quantity),
			Price:    float32(price),
		})
	}
	return nil
}

// init configures the parameters the process needs to run
func init() {
	// read command line options
	flag.Var(&cmdLnInvestments, "invest", "Formatted investment in the form of \"Ticker,Quantity,Price\".")
	stocks := flag.String("ticker", getEnvString("TICKERS", ""), "(TICKERS)\nComma seperated list of stocks to report.")
	mailAddress := flag.String("email", getEnvString("EMAIL_ADDR", ""), "(EMAIL_ADDR)\nE-mail address that will receive the end of day summary.")
	mailHost := flag.String("host", getEnvString("EMAIL_HOST", ""), "(EMAIL_HOST)\nE-Mail server hostname.")
	mailPort := flag.Int("port", getEnvInt("EMAIL_PORT", 25), "(EMAIL_PORT)\nE-Mail server port.")
	mailFrom := flag.String("from", getEnvString("EMAIL_FROM", "StockWatch <noreply@localhost>"), "(EMAIL_FROM)\nAddress to send mail from.")
	flag.Parse()

	cmdLnStocks = *stocks
	cmdLnEmailAddress = *mailAddress
	cmdLnEmailHost = *mailHost
	cmdLnEmailPort = *mailPort
	cmdLnEmailFrom = *mailFrom

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
}
