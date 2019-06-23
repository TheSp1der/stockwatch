package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
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

func notifyViaMail(sData chan stockData) {
	for {
		open, sleepTime := marketStatus()
		if !open {
			time.Sleep(time.Duration(time.Minute * 5))
			if err := basicMailSend(stockwatchConfig.Mail.Host+":"+strconv.Itoa(stockwatchConfig.Mail.Port), stockwatchConfig.Mail.Address, stockwatchConfig.Mail.From, "Stock Alert", displayHTML(<-sData)); err != nil {
				log.Println(err.Error())
			}
		}

		time.Sleep(time.Duration(sleepTime))
	}
}

func outputConsole(sData chan stockData) {
	for {
		if terminal.IsTerminal(int(os.Stdout.Fd())) {
			fmt.Printf("\033[0;0H")
		}
		fmt.Print(displayTerminal(<-sData))
		time.Sleep(time.Duration(time.Second * 5))
	}
}
