package main

//https://iexcloud.io/docs/api/#quote

import (
	"errors"
	"flag"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// global variables
var (
	timeFormat = "2006-01-02 15:04:05"
	log        = logrus.New()
)

// init configures the parameters the process needs to run.
func init() {
	// set default log level
	log.SetLevel(logrus.WarnLevel)
}

func main() {
	// read command line options
	swConfig, investments, stocks := readFlags()

	// display current configuration
	displayCurrentConfig(swConfig, stocks)

	err := validateStartParameters(swConfig, stocks)
	if err != nil {
		flag.PrintDefaults()
		log.Fatal(err)
	}

	swConfig, err = prepareStartParameters(swConfig, stocks)
	if err != nil {
		flag.PrintDefaults()
		log.Fatal(err.Error())
	}

	dReader := make(chan map[string]iexStock)
	sData := make(chan map[string]*stockData)

	// start processing market data
	go dataReader(dReader, swConfig)
	go dataDistributer(dReader, sData, swConfig)

	for i := 0; i < 5; i = i + 1 {
		for k, s := range <-sData {
			log.Infof("Stock: %v", k)
			log.Infof("  LastPrice: %v", s.StockDetail.LatestPrice)
		}

		time.Sleep(time.Second * 1)
	}

	log.Info(investments)

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

func validateStartParameters(config Configuration, stocks string) (err error) {
	// verify required data was provided
	if config.IexAPIKey == "" {
		return errors.New("IEX API Key was not provided")
	} else if stocks == "" {
		return errors.New("No stocks were defined")
	} else if config.HTTPPort <= 0 || config.HTTPPort > 65535 {
		return errors.New("Invalid port given for http server provided")
	} else if config.Mail.Port <= 0 || config.Mail.Port > 65535 {
		return errors.New("Invalid port for mail server provided")
	}

	return
}

func prepareStartParameters(config Configuration, stocks string) (Configuration, error) {
	// split provided stocks
	re := regexp.MustCompile(`(\s+)?,(\s+)?`)
	config.TrackedStocks = re.Split(stocks, -1)

	// add stocks from investments
	if len(config.Investments) > 0 {
		for _, i := range config.Investments {
			config.TrackedStocks = append(config.TrackedStocks, strings.ToLower(i.Ticker))
		}
	}

	// validate ticker
	re = regexp.MustCompile(`^[a-z0-9]+$`)
	for _, value := range config.TrackedStocks {
		if !re.Match([]byte(value)) {
			return config, errors.New("Stock ticker " + value + " is not a valid ticker name")
		}
	}

	// remove duplicate stock tickers
	config.TrackedStocks = uniqueString(config.TrackedStocks)

	log.Info("Initialization complete.")

	return config, nil
}
