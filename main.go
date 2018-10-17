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

	go updateStockData()

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
