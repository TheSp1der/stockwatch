package main

import (
	"strings"
	"time"

	"encoding/json"
	"net/url"

	"github.com/TheSp1der/goerror"
)

// updateStockData maintains up to date stock data (dependent on market status)
func updateStockData(sData chan<- iex) {
	var (
		err     error
		s       iex
		runTime = time.Now()
	)

	for {
		if time.Now().After(runTime) || time.Now().Equal(runTime) {
			s, err = getPrices()
			if err != nil {
				time.Sleep(time.Duration(time.Millisecond * 500))
				continue
			}
		}

		if time.Now().After(runTime) {
			open, openTime := marketStatus()
			if open {
				runTime = time.Now().Add(time.Duration(time.Second * 5))
			} else {
				runTime = time.Now().Add(time.Duration(time.Minute * 60))
				if time.Now().Add(openTime).Before(runTime) {
					runTime = time.Now().Add(openTime)
				}
			}
		}

		// non blocking channel read
		select {
		case sData <- s:
		default:
		}

		// slow the loop down so it does not use un-necessary cpu
		// this will also limit reads to once every 100 milliseconds
		time.Sleep(time.Duration(time.Millisecond * 100))
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
