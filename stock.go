package main

import (
	"strings"
	"time"

	"encoding/json"
	"net/url"

	"github.com/TheSp1der/httpclient"
)

// updateStockData manages updates to the reader channel
func dataReader(sData chan<- map[string]iexStock) {
	var (
		runTime = time.Now()
		s       = make(map[string]iexStock)
	)

	for {
		if time.Now().After(runTime) || time.Now().Equal(runTime) {
			for _, stock := range stockwatchConfig.TrackedStocks {
				i, err := getPrices(stock)
				if err != nil {
					time.Sleep(time.Millisecond * 500)
					continue
				}
				log.Printf("Symbol: %v", i.Symbol)
				s[i.Symbol] = i
			}
			// send the updated data to the channel
			sData <- s
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

		// sleep a little to keep cpu usage down
		time.Sleep(time.Millisecond * 50)
	}
}

func dataDistributer(newData <-chan map[string]iexStock, dataSender chan<- map[string]*stockData) {
	sData := make(map[string]*stockData)

	// get company data for stocks
	for _, stock := range stockwatchConfig.TrackedStocks {
		log.Printf("Getting company data for: %v", strings.ToUpper(stock))
		cN, err := getCompanyData(strings.ToUpper(stock))
		if err != nil {
			log.Fatal("Unable to obtain the company name for " + strings.ToUpper(stock))
		}

		sData[strings.ToUpper(stock)] = &stockData{
			CompanyData: cN,
		}
	}

	for {
		select {
		// read
		case tmp := <-newData:
			// update sData holder
			for _, stock := range tmp {
				sData[strings.ToUpper(stock.Symbol)] = &stockData{
					StockDetail: stock,
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
func getCompanyData(ticker string) (iexCompany, error) {
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
	req := httpclient.DefaultClient()
	req.SetHeader("Content-Type", "application/json")
	resp, err := req.Get(newURL.String())
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
func getPrices(stock string) (iexStock, error) {
	// prepare the url
	var newURL url.URL
	newURL.Scheme = "https"
	newURL.Host = "sandbox.iexapis.com"
	newURL.Path = "v1/stock/" + stock + "/quote"

	// url parameters
	params := newURL.Query()
	params.Add("token", stockwatchConfig.IexAPIKey)
	params.Add("format", "json")
	newURL.RawQuery = params.Encode()

	// connect and retrieve data from remote source
	req := httpclient.DefaultClient()
	req.SetHeader("Content-Type", "application/json")
	resp, err := req.Get(newURL.String())
	if err != nil {
		return iexStock{}, err
	}

	// unmarshal response
	var jsonIexStock iexStock
	if err = json.Unmarshal(resp, &jsonIexStock); err != nil {
		log.Println(err.Error())
	}

	return jsonIexStock, nil
}
