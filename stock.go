/******************************************************************************
	stock.go
	market interaction functrions for retrieving data, determining open
	trading hours, etc

	data is provided for free by the API from IEX, you can read	more about
	the API here:
	https://iextrading.com/developer/

	please read IEX's terms of use prior to using this program:
	https://iextrading.com/api-exhibit-a/
******************************************************************************/
package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/TheSp1der/goerror"
)

func stockMonitor() {
	var (
		err       error
		sleepTime time.Duration
	)

	for {
		if o, s := marketStatus(); o {
			sleepTime = time.Duration(time.Second * 5)
		} else {
			var stockData iex
			sleepTime = s

			if stockData, err = getPrices(); err != nil {
				goerror.Warning(err)
			}

			if err = sendMail(cmdLnEmailHost+":"+strconv.Itoa(cmdLnEmailPort), cmdLnEmailAddress, cmdLnEmailFrom, "Stock Alert", printPrices(stockData, false)); err != nil {
				goerror.Warning(err)
			}

			fmt.Println("Market is currently closed.")
			fmt.Println("Script will resume at " + time.Now().Add(time.Duration(sleepTime)).Format(timeFormat) + " which is in " + strconv.FormatFloat(sleepTime.Seconds(), 'f', 0, 64) + " seconds.")
		}

		stockCurrent()
		time.Sleep(sleepTime)
	}
}

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

	if ct.After(open) && ct.Before(close) {
		return true, 0
	}
	return false, open.Sub(ct)
}

func stockCurrent() {
	var (
		err       error
		stockData iex
	)

	if stockData, err = getPrices(); err != nil {
		goerror.Warning(err)
	}

	fmt.Println(printPrices(stockData, true))
}

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
