package main

import (
	"log"
	"strings"
	"time"

	"encoding/json"
	"net/url"
)

// updateStockData manages updates to the reader channel
func dataReader(sData chan<- iexTop) {
	var (
		runTime = time.Now()
		s       iexTop
	)

	for {
		if time.Now().After(runTime) || time.Now().Equal(runTime) {
			var err error
			if s, err = getPrices(); err != nil {
				time.Sleep(time.Millisecond * 500)
				continue
			}
		}

		if time.Now().After(runTime) {
			open, openTime := marketStatus()
			if open {
				runTime = time.Now().Add(time.Millisecond * time.Duration(stockwatchConfig.PollFrequency*1000))
			} else {
				runTime = time.Now().Add(time.Duration(time.Minute * 120))
				if time.Now().Add(openTime).Before(runTime) {
					runTime = time.Now().Add(openTime)
				}
			}
		}

		// send the updated data to the channel
		sData <- s

		// sleep a little to keep cpu usage down
		time.Sleep(time.Millisecond * 50)
	}
}

func dataDistributer(newData <-chan iexTop, dataSender chan<- map[string]*stockData) {
	sData := make(map[string]*stockData)

	// get company data for stocks
	for _, stock := range stockwatchConfig.TrackedStocks {
		log.Printf("Getting company data for: %v", strings.ToUpper(stock))
		cN, err := getCompanyName(strings.ToUpper(stock))
		if err != nil {
			log.Fatal("Unable to obtain the company name for " + strings.ToUpper(stock))
		}

		sData[strings.ToUpper(stock)] = &stockData{
			CompanyName: cN,
		}
	}

	for {
		select {
		// read
		case tmp := <-newData:
			// update sData holder
			for _, stock := range tmp {
				sData[strings.ToUpper(stock.Symbol)] = &stockData{
					Ask: stock.AskPrice,
					Bid: stock.BidPrice,
				}
			}

		// send
		case dataSender <- sData:
		default:
		}

		time.Sleep(time.Millisecond * 50)
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

	// get num of days until next open day
	switch ct.Weekday() {
	// Sun
	case 0:
		o = 1
		c = 1
	// Mon, Tues, Wed, Thur
	case 1, 2, 3, 4:
		if ct.After(time.Date(ct.Year(), ct.Month(), ct.Day(), 16, 0, 0, 0, est)) {
			o = 1
			c = 1
		}
	// Fri
	case 5:
		if ct.After(time.Date(ct.Year(), ct.Month(), ct.Day(), 16, 0, 0, 0, est)) {
			o = 3
			c = 3
		}
	// Sat
	case 6:
		o = 2
		c = 2
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

// get company name from ticker
func getCompanyName(ticker string) (iexCompany, error) {
	// prepare the url
	var newURL url.URL
	newURL.Scheme = "https"
	newURL.Host = "sandbox.iexapis.com"
	newURL.Path = "v1/stock/" + ticker + "/company"

	// url parameters
	params := newURL.Query()
	params.Add("token", stockwatchConfig.IexAPIKey)
	params.Add("format", "json")
	newURL.RawQuery = params.Encode()

	// connect and retrieve data from remote source
	resp, err := httpGet(
		newURL.String(),
		[]struct {
			Name, Value string
		}{
			{
				Name:  "Content-Type",
				Value: "application/json",
			},
		})
	if err != nil {
		return iexCompany{}, err
	}

	// unmarshal response
	var jsonIexCompany iexCompany
	if err = json.Unmarshal(resp, &jsonIexCompany); err != nil {
		log.Println(err.Error())
	}

	return jsonIexCompany, nil
}

// getPrices will get the current stock data.
func getPrices() (iexTop, error) {
	// prepare the url
	var newURL url.URL
	newURL.Scheme = "https"
	newURL.Host = "sandbox.iexapis.com"
	newURL.Path = "v1/tops"

	// url parameters
	params := newURL.Query()
	params.Add("token", stockwatchConfig.IexAPIKey)
	params.Add("format", "json")
	params.Add("symbols", strings.Join(stockwatchConfig.TrackedStocks, ","))
	newURL.RawQuery = params.Encode()

	log.Printf("URL: %v", newURL.String())

	// connect and retrieve data from remote source
	resp, err := httpGet(
		newURL.String(),
		[]struct {
			Name, Value string
		}{
			{
				Name:  "Content-Type",
				Value: "application/json",
			},
		})
	if err != nil {
		return iexTop{}, err
	}

	// unmarshal response
	var jsonIexTop iexTop
	if err = json.Unmarshal(resp, &jsonIexTop); err != nil {
		log.Println(err.Error())
	}

	return jsonIexTop, nil
}
