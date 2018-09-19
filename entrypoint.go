/******************************************************************************
	entrypoint.go
	This is the entrypoint for the go process. All other functions are
	initiated from this main function.

	Global variables are also defined here as a matter of convenience.
******************************************************************************/

package main

// import external libaries
import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"net/url"

	"github.com/TheSp1der/goerror"
	"github.com/fatih/color"
)

// global variables
var (
	cmdLnStocks       string
	cmdLnEmailAddress string
	cmdLnEmailHost    string
	cmdLnEmailPort    int
	cmdLnEmailFrom    string

	trackedTickers []string

	timeFormat = "2006-01-02 15:04:05"
)

// init
// ----
// process initialization
//
// input:
//
// return:
func init() {
	var (
		err     error
		tickers string
	)

	// read command line options
	flag.StringVar(&cmdLnStocks, "ticker", "", "Comma seperated list of stocks to report (TICKER)")
	flag.StringVar(&cmdLnEmailAddress, "email", "", "To address or EOD message (EMAIL_ADDR)")
	flag.StringVar(&cmdLnEmailHost, "host", "localhost", "E-Mail server hostname (EMAIL_HOST)")
	flag.IntVar(&cmdLnEmailPort, "port", 25, "E-Mail server port (EMAIL_PORT)")
	flag.StringVar(&cmdLnEmailFrom, "from", "StockWatch <noreply@localhost>", "Address to send mail from (EMAIL_FROM)")
	flag.Parse()

	// read options from environment variables
	if cmdLnStocks == "" && len(os.Getenv("TICKER")) > 0 {
		tickers = strings.ToLower(os.Getenv("TICKER"))
	}
	if cmdLnEmailAddress == "" && len(os.Getenv("EMAIL_ADDR")) > 0 {
		cmdLnEmailAddress = os.Getenv("EMAIL_ADDR")
	}
	if cmdLnEmailHost == "" && len(os.Getenv("EMAIL_HOST")) > 0 {
		cmdLnEmailHost = os.Getenv("EMAIL_HOST")
	}
	if cmdLnEmailPort == 25 && len(os.Getenv("EMAIL_PORT")) > 0 {
		if cmdLnEmailPort, err = strconv.Atoi(os.Getenv("EMAIL_PORT")); err != nil {
			goerror.Fatal(errors.New("EMAIL_PORT must be numeric."))
		}
	}
	if cmdLnEmailFrom == "StockWatch <noreply@localhost>" && len(os.Getenv("EMAIL_FROM")) > 0 {
		cmdLnEmailFrom = os.Getenv("EMAIL_FROM")
	}

	// get tickers from command line (this overrides environment variables)
	if len(cmdLnStocks) > 0 {
		tickers = strings.ToLower(cmdLnStocks)
	}

	// convert input to struct
	re := regexp.MustCompile(`(\s+)?,(\s+)?`)
	trackedTickers = re.Split(tickers, -1)

	// verify stocks were provided
	if len(trackedTickers) == 1 && trackedTickers[0] == "" {
		goerror.Fatal(errors.New("no Stocks defined"))
	} else {
		re = regexp.MustCompile(`^[a-z0-9]+$`)
		for _, value := range trackedTickers {
			if !re.Match([]byte(value)) {
				goerror.Fatal(errors.New("stock ticker format error"))
			}
		}
	}
}

// main
// ----
// process entrypoint
//
// input:
//
// return:
func main() {
	if cmdLnEmailAddress != "" && cmdLnEmailFrom != "" && cmdLnEmailHost != "" {
		stockMonitor()
	} else {
		stockCurrent()
	}
}

func stockMonitor() {
	for {

		var (
			err       error
			o         int
			c         int
			sleepTime time.Duration
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
			sleepTime = time.Duration(time.Second * 5)
		} else {
			sleepTime = open.Sub(ct)

			var stockData iex
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

func stockCurrent() {
	var (
		err       error
		stockData iex
	)

	if stockData, err = getPrices(); err != nil {
		goerror.Warning(err)
	}

	fmt.Println(printPrices(stockData, true))

	var shares float64 = 865
	investment := (shares * 44.51)
	value := (shares * stockData["BP"].Price)

	fmt.Println("\n" + strconv.FormatFloat(value-investment, 'f', 2, 64))
}

func printPrices(stockData iex, text bool) string {
	var (
		output string
	)

	keys := make([]string, 0, len(stockData))
	for k := range stockData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	if text {
		output += "Stock report as of " + color.BlueString(time.Now().Format(timeFormat)) + "\n"
		output += "Company" + "\t\t" + "Current Price" + "\t" + "Change" + "\n"
		for _, k := range keys {
			if stockData[k].Quote.Change < 0 {
				output += stockData[k].Company.Symbol + "\t\t" + strconv.FormatFloat(stockData[k].Price, 'f', 2, 64) + "\t\t" + color.RedString(strconv.FormatFloat(stockData[k].Quote.Change, 'f', 2, 64)) + "\t\t" + stockData[k].Company.CompanyName + "\n"
			} else if stockData[k].Quote.Change > 0 {
				output += stockData[k].Company.Symbol + "\t\t" + strconv.FormatFloat(stockData[k].Price, 'f', 2, 64) + "\t\t" + color.GreenString(strconv.FormatFloat(stockData[k].Quote.Change, 'f', 2, 64)) + "\t\t" + stockData[k].Company.CompanyName + "\n"
			} else {
				output += stockData[k].Company.Symbol + "\t\t" + strconv.FormatFloat(stockData[k].Price, 'f', 2, 64) + "\t\t\t\t" + stockData[k].Company.CompanyName + "\n"
			}
		}
	} else {
		output += "<span style=\"font-weight: bold;\">Stock report as of " + time.Now().Format(timeFormat) + "</span><br>\n"
		output += "<table>\n"
		output += "\t<tr>\n"
		output += "\t\t<th style=\"text-align: left;\">Company</th>\n"
		output += "\t\t<th style=\"text-align: left;\">Closing Price</th>\n"
		output += "\t\t<th style=\"text-align: left;\">Change</th>\n"
		output += "\t</tr>\n"
		for _, k := range keys {
			output += "\t<tr>\n"
			output += "\t\t<td><a href=\"" + stockData[k].Company.Website + "\">" + stockData[k].Company.CompanyName + "</a></td>\n"
			output += "\t\t<td>" + strconv.FormatFloat(stockData[k].Price, 'f', 2, 64) + "</td>\n"
			if stockData[k].Quote.Change < 0 {
				output += "\t\t<td><span style=\"color: red;\">" + strconv.FormatFloat(stockData[k].Quote.Change, 'f', 2, 64) + "</span></td>\n"
			} else if stockData[k].Quote.Change > 0 {
				output += "\t\t<td><span style=\"color: green;\">" + strconv.FormatFloat(stockData[k].Quote.Change, 'f', 2, 64) + "</span></td>\n"
			} else {
				output += "\t\t<td>" + strconv.FormatFloat(stockData[k].Quote.Change, 'f', 2, 64) + "</td>\n"
			}
			output += "\t</tr>\n"
		}
		output += "</table>"
		output += "<br>"
		output += "<span style=\"font-weight: bold;\">Graphs:</span><br>"
		for _, k := range keys {
			output += "<img src=\"https://finviz.com/chart.ashx?t=" + stockData[k].Company.Symbol + "\"><br>"
		}
	}

	return output
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

func httpGet(url string, header httpHeader) ([]byte, error) {
	var (
		err    error
		client http.Client
		req    *http.Request
		resp   *http.Response
		output []byte
	)

	// set timeouts
	client = http.Client{
		Timeout: time.Duration(2 * time.Second),
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: time.Duration(2 * time.Second),
			}).Dial,
			TLSHandshakeTimeout: time.Duration(2 * time.Second),
		},
	}

	// set up request
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		return []byte(""), err
	}

	// set headers
	if len(header) > 0 {
		for _, h := range header {
			req.Header.Set(h.Name, h.Value)
		}
	}

	// perform the request
	if resp, err = client.Do(req); err != nil {
		return []byte(""), err
	}

	// convert the response to byte
	output, err = ioutil.ReadAll(resp.Body)

	// close the connection
	defer resp.Body.Close()

	return output, nil
}

func sendMail(host string, to string, from string, subject string, body string) error {
	var message string

	// connect to the remote server
	client, err := smtp.Dial(host)
	if err != nil {
		return err
	}
	defer client.Close()

	// set sender and and recipient
	client.Mail(from)
	client.Rcpt(to)

	// send the body
	mailContent, err := client.Data()
	if err != nil {
		return err
	}
	defer mailContent.Close()

	message = "From: " + from + "\n"
	message += "To: " + to + "\n"
	message += "Subject: " + subject + "\n"
	message += "MIME-Version: 1.0\n"
	message += "Content-Type: text/html; charset=UTF-8\n"
	message += "<html>\n"
	message += "<body>\n"
	message += body
	message += "</body>\n"
	message += "</html>\n"
	message += "\n"

	buf := bytes.NewBufferString(message)
	if _, err = buf.WriteTo(mailContent); err != nil {
		return err
	}

	return nil
}