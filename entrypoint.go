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
	"fmt"
	"net"
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

	if len(cmdLnInvestments) > 0 {
		for _, i := range cmdLnInvestments {
			fmt.Println(i.Ticker)
			fmt.Println(i.Quantity)
			fmt.Println(i.Price)
		}
	}
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

func httpGet(url string, headers httpHeader) ([]byte, error) {
	var (
		err    error          // error handler
		client http.Client    // http client
		req    *http.Request  // http request
		res    *http.Response // http response
		output []byte         // output
	)

	// set timeouts
	client = http.Client{
		Timeout: time.Duration(time.Second * 2),
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: time.Duration(time.Second * 2),
			}).Dial,
			TLSHandshakeTimeout: time.Duration(time.Second * 2),
		},
	}

	// setup request
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		return output, err
	}

	// setup headers
	if len(headers) > 0 {
		for _, header := range headers {
			req.Header.Set(header.Name, header.Value)
		}
	}

	// perform the request
	if res, err = client.Do(req); err != nil {
		return output, err
	}

	// close the connection upon function closure
	defer res.Body.Close()

	// extract response body
	if output, err = ioutil.ReadAll(res.Body); err != nil {
		return output, err
	}

	// check status
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return output, errors.New("non-successful status code received [" + strconv.Itoa(res.StatusCode) + "]")
	}

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
