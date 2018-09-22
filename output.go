/******************************************************************************
	output.go
	prepares and formates output for display
******************************************************************************/
package main

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

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
		var (
			gtol    float32
			gtolStr string
		)

		if m, _ := marketStatus(); m {
			output += color.GreenString(time.Now().Format(timeFormat)) + "\n"
		} else {
			output += color.YellowString(time.Now().Format(timeFormat)) + "\n"
		}

		output += ".---------------------.------------.------------.------------.\n"
		output += "| Company             | Price      | Change     | Investment |\n"
		output += "|---------------------|------------|------------|------------|\n"
		for _, k := range keys {
			var (
				cmpy    string  // company name
				prce    string  // price
				chge    string  // change
				ival    float32 // investment value
				cval    float32 // current value
				diff    float32 // difference
				totl    float32 // total difference
				totlStr string  // total difference (string)
			)

			for _, i := range cmdLnInvestments {
				if strings.TrimSpace(strings.ToLower(stockData[k].Company.Symbol)) == strings.TrimSpace(strings.ToLower(i.Ticker)) {
					ival = i.Quantity * i.Price
					cval = i.Quantity * float32(stockData[k].Price)
					diff = cval - ival
				}
				totl = totl + diff
			}

			gtol = gtol + totl
			cmpy = blockedOutput(stockData[k].Company.CompanyName, 19)
			prce = blockedOutput(strconv.FormatFloat(stockData[k].Price, 'f', 2, 64), 10)
			if stockData[k].Quote.Change < 0 {
				chge = color.RedString(blockedOutput(strconv.FormatFloat(stockData[k].Quote.Change, 'f', 2, 64), 10))
			} else if stockData[k].Quote.Change > 0 {
				chge = color.GreenString(blockedOutput(strconv.FormatFloat(stockData[k].Quote.Change, 'f', 2, 64), 10))
			} else {
				chge = blockedOutput(strconv.FormatFloat(stockData[k].Quote.Change, 'f', 2, 64), 10)
			}

			if totl < 0 {
				totlStr = color.RedString(blockedOutput(strconv.FormatFloat(float64(totl), 'f', 2, 64), 10))
			} else if totl > 0 {
				totlStr = color.GreenString(blockedOutput(strconv.FormatFloat(float64(totl), 'f', 2, 64), 10))
			} else {
				totlStr = blockedOutput(strconv.FormatFloat(float64(totl), 'f', 2, 64), 10)
			}

			output += "| " + cmpy + " | " + prce + " | " + chge + " | " + totlStr + " |\n"
		}
		if gtol < 0 {
			gtolStr = color.RedString(blockedOutput(strconv.FormatFloat(float64(gtol), 'f', 2, 64), 10))
		} else if gtol > 0 {
			gtolStr = color.GreenString(blockedOutput(strconv.FormatFloat(float64(gtol), 'f', 2, 64), 10))
		} else {
			gtolStr = blockedOutput(strconv.FormatFloat(float64(gtol), 'f', 2, 64), 10)
		}

		output += "|---------------------'------------'------------'------------|\n"
		output += "| Total Investment Value:                         " + gtolStr + " |\n"
		output += "`------------------------------------------------------------'\n"
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

func blockedOutput(input string, width int) string {
	r := []rune(input)

	if len(r) > width {
		return string(r[0:width])
	} else if len(r) < width {
		s := width - len(r)
		return string(r) + strings.Repeat(" ", s)
	}
	return input
}
