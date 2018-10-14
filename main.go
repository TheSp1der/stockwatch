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

	// continually update stock data
	go func() {
		var (
			err     error
			runTime = time.Now()
		)

		for {
			if time.Now().After(runTime) || time.Now().Equal(runTime) {
				sData, err = getPrices()
				if err != nil {
					goerror.Warning(err)
				}
			}

			open, openTime := marketStatus()
			if open {
				runTime = time.Now().Add(time.Duration(time.Second * 5))
			} else {
				if time.Now().After(runTime) {
					runTime = time.Now().Add(time.Duration(time.Minute * 60))
					if time.Now().Add(openTime).Before(runTime) {
						runTime = time.Now().Add(openTime)
					}
				}
			}

			time.Sleep(time.Duration(time.Millisecond * 100))
		}
	}()

	if !cmdLnNoConsole {
		go func() {
			for {
				if terminal.IsTerminal(int(os.Stdout.Fd())) {
					fmt.Printf("\033[0;0H")
				}
				fmt.Println(displayTerminal(sData))
				time.Sleep(time.Duration(time.Second * 5))
			}
		}()
	}

	if cmdLnHTTPPort > 0 && cmdLnHTTPPort < 65535 {
		go webListener(cmdLnHTTPPort)
	}

	if cmdLnEmailAddress != "" && cmdLnEmailFrom != "" && cmdLnEmailHost != "" {
		go func() {
			for {
				open, sleepTime := marketStatus()
				if !open {
					time.Sleep(time.Duration(time.Minute * 5))
					if err := basicMailSend(cmdLnEmailHost+":"+strconv.Itoa(cmdLnEmailPort), cmdLnEmailAddress, cmdLnEmailFrom, "Stock Alert", displayHTML(sData)); err != nil {
						goerror.Warning(err)
					}
				}

				time.Sleep(time.Duration(sleepTime))
			}
		}()
	}

	if !cmdLnNoConsole || (cmdLnHTTPPort > 0 && cmdLnHTTPPort < 65535) || (cmdLnEmailAddress != "" && cmdLnEmailFrom != "" && cmdLnEmailHost != "") {
		for {
			time.Sleep(time.Duration(time.Second * 5))
		}
	}
}
