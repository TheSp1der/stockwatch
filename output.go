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
		if m, _ := marketStatus(); m {
			output += "            " + color.GreenString(time.Now().Format(timeFormat)) + "\n"
		} else {
			output += "             " + color.YellowString(time.Now().Format(timeFormat)) + "\n"
		}
		output += ".---------------------.------------.------------.\n"
		output += "| Company             | Price      | Change     |\n"
		output += "|---------------------|------------|------------|\n"
		for _, k := range keys {
			c := blockedOutput(stockData[k].Company.CompanyName, 19)
			p := blockedOutput(strconv.FormatFloat(stockData[k].Price, 'f', 2, 64), 10)
			var ch string

			if stockData[k].Quote.Change < 0 {
				ch = color.RedString(blockedOutput(strconv.FormatFloat(stockData[k].Quote.Change, 'f', 2, 64), 10))
			} else if stockData[k].Quote.Change > 0 {
				ch = color.GreenString(blockedOutput(strconv.FormatFloat(stockData[k].Quote.Change, 'f', 2, 64), 10))
			} else {
				ch = blockedOutput(strconv.FormatFloat(stockData[k].Quote.Change, 'f', 2, 64), 10)
			}
			output += "| " + c + " | " + p + " | " + ch + " |\n"
		}
		output += "`---------------------'------------'------------'\n"
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
