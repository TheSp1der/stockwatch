package main

//https://iexcloud.io/docs/api/#quote

import (
	"errors"
	"flag"
	"log"
	"regexp"
	"strings"
	"time"
)

func validateStartParameters() error {
	// verify required data was provided
	if stockwatchConfig.IexAPIKey == "" {
		return errors.New("IEX API Key was not provided")
	} else if cmdLnStocks == "" {
		return errors.New("No stocks were defined")
	} else if stockwatchConfig.HTTPPort < 0 && stockwatchConfig.HTTPPort > 65535 {
		return errors.New("Invalid port given for http server provided")
	} else if stockwatchConfig.Mail.Port < 0 && stockwatchConfig.Mail.Port > 65534 {
		return errors.New("Invalid port for mail server provided")
	}

	return nil
}

func prepareStartParameters() error {
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
			return errors.New("Stock ticker " + value + " is not a valid ticker name")
		}
	}

	// remove duplicate stock tickers
	stockwatchConfig.TrackedStocks = uniqueString(stockwatchConfig.TrackedStocks)

	log.Println("Initialization complete.")

	return nil
}

func main() {
	if err := validateStartParameters(); err != nil {
		flag.PrintDefaults()
		log.Fatal(err.Error())
	}

	if err := prepareStartParameters(); err != nil {
		flag.PrintDefaults()
		log.Fatal(err.Error())
	}

	dReader := make(chan iexTop)
	sData := make(chan map[string]*stockData)

	// start processing market data
	go dataReader(dReader)
	go dataDistributer(dReader, sData)

	for i := 0; i < 5; i = i + 1 {
		for k, s := range <-sData {
			log.Printf("Stock: %v", k)
			log.Printf("  Bid: %v", s.Bid)
			log.Printf("  Ask: %v", s.Ask)
			log.Printf("  Last: %v", s.Last)
		}

		time.Sleep(time.Second * 1)
	}

	/*
		// start outputting to the console
		if !stockwatchConfig.NoConsole {
			go outputConsole(sData)
		}

		// start the web listener
		if stockwatchConfig.HTTPPort != 0 {
			go webListener(sData, stockwatchConfig.HTTPPort)
		}

		// start e-mail notifier
		if stockwatchConfig.Mail.Address != "" && stockwatchConfig.Mail.From != "" && stockwatchConfig.Mail.Host != "" {
			go notifyViaMail(sData)
		}

		// run the program infinitely
		if !stockwatchConfig.NoConsole ||
			stockwatchConfig.HTTPPort != 0 ||
			(stockwatchConfig.Mail.Address != "" && stockwatchConfig.Mail.From != "" && stockwatchConfig.Mail.Host != "" && stockwatchConfig.Mail.Port != 0) {
			for {
				time.Sleep(time.Duration(time.Second * 5))
			}
		}
	*/
}
