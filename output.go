/******************************************************************************
	output.go
	prepares and formates output for display
******************************************************************************/
package main

import (
	"sort"
	"strconv"
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
