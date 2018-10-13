package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"encoding/json"
	"net/url"

	"github.com/TheSp1der/goerror"
	"golang.org/x/crypto/ssh/terminal"
)

// stockCurrent will retrieve the market values and write the
// results to the terminal.
func stockCurrent() {
	var (
		err       error
		stockData iex
	)

	if stockData, err = getPrices(); err != nil {
		goerror.Warning(err)
	}

	// output the formatted stock information
	fmt.Println(displayTerminal(stockData))
}

// stockMonitor is the entrypoint for monitoring the market in an
// infinite method, must have e-mail configured for EOD messages.
func stockMonitor() {
	var (
		err       error
		sleepTime time.Duration
	)

	// begin loop
	for {
		// if market is open, sleep for 5 seconds
		if o, s := marketStatus(); o {
			sleepTime = time.Duration(time.Second * 5)
			// if market is closed send EOD message and sleep until it opens
		} else {
			var stockData iex
			sleepTime = s

			if stockData, err = getPrices(); err != nil {
				goerror.Warning(err)
			}

			if err = basicMailSend(cmdLnEmailHost+":"+strconv.Itoa(cmdLnEmailPort), cmdLnEmailAddress, cmdLnEmailFrom, "Stock Alert", displayHTML(stockData)); err != nil {
				goerror.Warning(err)
			}

			fmt.Println("Market is currently closed.")
			fmt.Println("Script will resume at " + time.Now().Add(time.Duration(sleepTime)).Format(timeFormat) + " which is in " + strconv.FormatFloat(sleepTime.Seconds(), 'f', 0, 64) + " seconds.")
		}

		// if verbose update terminal
		if cmdLnVerbose {

			// check to see if we are running in a terminal
			if terminal.IsTerminal(int(os.Stdout.Fd())) {
				// reset the location of the cursor
				fmt.Printf("\033[0;0H")
			}

			// print the current prices to the screen
			stockCurrent()
		}

		// sleep
		time.Sleep(sleepTime)
	}
}

// marketStatus will determine if the market is open, if it is closed
// it will return the time until it is open again.
func marketStatus() (bool, time.Duration) {
	var (
		o int
		c int
	)

	// get current new york time
	est, _ := time.LoadLocation("America/New_York")
	ct := time.Now().In(est)

	// get open and closed times
	switch ct.Weekday() {
	case 1, 2, 3, 4:
		if ct.After(time.Date(ct.Year(), ct.Month(), ct.Day(), 16, 0, 0, 0, est)) {
			o = 1
			c = 1
		}
	case 5:
		if ct.After(time.Date(ct.Year(), ct.Month(), ct.Day(), 16, 0, 0, 0, est)) {
			o = 3
			c = 3
		}
	case 6:
		o = 2
		c = 2
	default:
		o = 1
		o = 1
	}
	open := time.Date(ct.Year(), ct.Month(), ct.Day()+o, 9, 30, 0, 0, est)
	close := time.Date(ct.Year(), ct.Month(), ct.Day()+c, 16, 0, 0, 0, est)

	// if the market is open return true
	if ct.After(open) && ct.Before(close) {
		return true, 0
	}

	// otherwise return false with time until it is open
	return false, open.Sub(ct)
}

// getPrices will get the current stock data.
func getPrices() (iex, error) {
	var (
		err       error
		newURL    url.URL
		params    url.Values
		headers   httpHeader
		resp      []byte
		stockData iex
	)

	// prepare the url
	newURL.Scheme = "https"
	newURL.Host = "api.iextrading.com"
	newURL.Path = "1.0/stock/market/batch"

	// url parameters
	params = newURL.Query()
	params.Add("symbols", strings.Join(trackedTickers, ","))
	params.Add("types", "quote,price,company,stats,ohlc")
	newURL.RawQuery = params.Encode()

	// connect and retrieve data from remote source
	if resp, err = httpGet(newURL.String(), headers); err != nil {
		return stockData, err
	}

	// unmarshal response
	if err = json.Unmarshal(resp, &stockData); err != nil {
		goerror.Info(err)
	}

	return stockData, nil
}
