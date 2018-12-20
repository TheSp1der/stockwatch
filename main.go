package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/TheSp1der/goerror"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	sData := make(chan iex)

	// get current prices
	go updateStockData(sData)

	// output to console
	if !cmdLnNoConsole {
		go outputConsole(sData)
	}

	// output to webserver
	if cmdLnHTTPPort > 0 && cmdLnHTTPPort < 65535 {
		go webListener(sData, cmdLnHTTPPort)
	}

	// end of day e-mail
	if cmdLnEmailAddress != "" && cmdLnEmailFrom != "" && cmdLnEmailHost != "" {
		go notifyViaMail(sData)
	}

	// if everything is fine, loop indefinitely
	if !cmdLnNoConsole || (cmdLnHTTPPort > 0 && cmdLnHTTPPort < 65535) || (cmdLnEmailAddress != "" && cmdLnEmailFrom != "" && cmdLnEmailHost != "") {
		for {
			time.Sleep(time.Duration(time.Second * 5))
		}
	}
}

func notifyViaMail(sData chan iex) {
	for {
		open, sleepTime := marketStatus()
		if !open {
			time.Sleep(time.Duration(time.Minute * 5))
			if err := basicMailSend(cmdLnEmailHost+":"+strconv.Itoa(cmdLnEmailPort), cmdLnEmailAddress, cmdLnEmailFrom, "Stock Alert", displayHTML(<-sData)); err != nil {
				goerror.Warning(err)
			}
		}

		time.Sleep(time.Duration(sleepTime))
	}
}

func outputConsole(sData chan iex) {
	var hData iex

	for {
		cData := <-sData

		// set cursor to top left position
		if terminal.IsTerminal(int(os.Stdout.Fd())) {
			fmt.Printf("\033[0;0H")
		}

		// display output
		fmt.Print(displayTerminal(hData, cData))

		// update historical data
		hData = cData

		// sleep
		time.Sleep(time.Duration(time.Second * 5))
	}
}
