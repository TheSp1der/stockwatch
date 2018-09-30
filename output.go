package main

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

// displayTermnal returns a string for display in the terminal window of
// calculated and tracked stocks and the overall gains/losses of provided
// investments.
func displayTerminal(stockData iex) string {
	var (
		gtol    float32
		gtolStr string
		output  string
	)

	// sorting stocks for display
	keys := make([]string, 0, len(stockData))
	for k := range stockData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// timestamp of the current run, color is determined by the status
	// of the market (open vs. closed)
	if m, _ := marketStatus(); m {
		output += color.GreenString(time.Now().Format(timeFormat)) + "\n"
	} else {
		output += color.YellowString(time.Now().Format(timeFormat)) + "\n"
	}

	// table header
	output += ".--------------------------------.--------------.----------------.------------.\n"
	output += "| Company Name                   | Market Value | Today's Change | Gain/Loss  |\n"
	output += "|--------------------------------|--------------|----------------|------------|\n"

	// individual stock lines
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

		// calculate the total for the ticker in the event a stock
		// has multiple investments
		for _, i := range cmdLnInvestments {
			if strings.TrimSpace(strings.ToLower(stockData[k].Company.Symbol)) == strings.TrimSpace(strings.ToLower(i.Ticker)) {
				ival = i.Quantity * i.Price
				cval = i.Quantity * float32(stockData[k].Price)
				diff = cval - ival
				totl = totl + diff
			}
		}

		// update the grand total loss/gain
		gtol = gtol + totl

		// update the variable for each displayed table data
		cmpy = alignLeft(stockData[k].Company.CompanyName, 30)
		prce = alignRight(strconv.FormatFloat(stockData[k].Price, 'f', 2, 64), 12)
		if stockData[k].Quote.Change < 0 {
			chge = color.RedString(alignRight(strconv.FormatFloat(stockData[k].Quote.Change, 'f', 2, 64), 14))
		} else if stockData[k].Quote.Change > 0 {
			chge = color.GreenString(alignRight(strconv.FormatFloat(stockData[k].Quote.Change, 'f', 2, 64), 14))
		} else {
			chge = alignRight("", 14)
		}

		if totl < 0 {
			totlStr = color.RedString(alignRight(strconv.FormatFloat(float64(totl), 'f', 2, 64), 10))
		} else if totl > 0 {
			totlStr = color.GreenString(alignRight(strconv.FormatFloat(float64(totl), 'f', 2, 64), 10))
		} else {
			totlStr = alignRight("", 10)
		}

		// record that lines output
		output += "| " + cmpy + " | " + prce + " | " + chge + " | " + totlStr + " |\n"
	}

	// update the grand total loss/gain if it has value
	if gtol < 0 {
		gtolStr = color.RedString(alignRight(strconv.FormatFloat(float64(gtol), 'f', 2, 64), 10))
	} else if gtol > 0 {
		gtolStr = color.GreenString(alignRight(strconv.FormatFloat(float64(gtol), 'f', 2, 64), 10))
	}

	// record footer (differ if grand total had value)
	if gtol != 0 {
		output += "|--------------------------------'--------------'----------------'------------|\n"
		output += "| Total Investment Value:                                          " + gtolStr + " |\n"
		output += "`-----------------------------------------------------------------------------'\n"
	} else {
		output += "`--------------------------------'--------------'----------------'------------'\n"
	}

	return output
}

// displayHTML returns a string for e-mail messages of calculated and
// tracked stocks and the overall gains/losses of provided investments.
func displayHTML(stockData iex) string {
	var (
		gtol   float32
		output string
	)

	// sorting stocks for display
	keys := make([]string, 0, len(stockData))
	for k := range stockData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// write html header data (my html is always tab indented)
	output += "<span style=\"font-weight: bold;\">Stock report as of " + time.Now().Format(timeFormat) + "</span><br>\n"
	output += "<table>\n"
	output += "\t<tr>\n"
	output += "\t\t<th style=\"text-align: left;\">Company Name</th>\n"
	output += "\t\t<th style=\"text-align: left;\">Market Value</th>\n"
	output += "\t\t<th style=\"text-align: left;\">Today's Change</th>\n"
	output += "\t\t<th style=\"text-align: left;\">Gain/Loss</th>\n"
	output += "\t</tr>\n"
	for _, k := range keys {
		var (
			ival float32 // investment value
			cval float32 // current value
			diff float32 // difference
			totl float32 // total difference
		)

		// calculate the total for the ticker in the event a stock
		// has multiple investments
		for _, i := range cmdLnInvestments {
			if strings.TrimSpace(strings.ToLower(stockData[k].Company.Symbol)) == strings.TrimSpace(strings.ToLower(i.Ticker)) {
				ival = i.Quantity * i.Price
				cval = i.Quantity * float32(stockData[k].Price)
				diff = cval - ival
				totl = totl + diff
			}
		}
		
		// update the grand total loss/gain
		gtol = gtol + totl

		// record the table row and cell data
		output += "\t<tr>\n"
		output += "\t\t<td>" + stockData[k].Company.CompanyName + "</td>\n"
		output += "\t\t<td style=\"text-align: right;\">" + strconv.FormatFloat(stockData[k].Price, 'f', 2, 64) + "</td>\n"


		if stockData[k].Quote.Change < 0 {
			output += "\t\t<td style=\"text-align: right; color: red;\">" + strconv.FormatFloat(stockData[k].Quote.Change, 'f', 2, 64) + "</td>\n"
		} else if stockData[k].Quote.Change > 0 {
			output += "\t\t<td style=\"text-align: right; color: green;\">" + strconv.FormatFloat(stockData[k].Quote.Change, 'f', 2, 64) + "</td>\n"
		} else {
			output += "\t\t<td></td>\n"
		}

		if totl < 0 {
			output += "\t\t<td style=\"text-align: right; color: red;\">" + strconv.FormatFloat(float64(totl), 'f', 2, 64) + "</td>\n"
		} else if totl > 0 {
			output += "\t\t<td style=\"text-align: right; color: green;\">" + strconv.FormatFloat(float64(totl), 'f', 2, 64) + "</td>\n"
		} else {
			output += "\t\t<td></td>\n"
		}
		output += "\t</tr>\n"
	}

	// close the table
	output += "</table>\n"
	output += "<br>\n"

	// record the grand total loss/gain if it has value
	if gtol < 0 {
		output += "<span style=\"font-weight: bold;\">Overall Performance: <span style=\"color: red;\">" + strconv.FormatFloat(float64(gtol), 'f', 2, 64) + "</span></span>\n"
	} else if gtol > 0 {
		output += "<span style=\"font-weight: bold;\">Overall Performance: <span style=\"color: green;\">" + strconv.FormatFloat(float64(gtol), 'f', 2, 64) + "</span></span>\n"
	}
	output += "<br>\n"
	output += "<span style=\"font-weight: bold;\">Graphs:</span><br>\n"
	for _, k := range keys {
		output += "<img src=\"https://finviz.com/chart.ashx?t=" + stockData[k].Company.Symbol + "\"><br>\n"
	}

	return output
}

// alignLeft will format the table data to the left of the cell
// and will trim off characters in the event the output is too
// long.
func alignLeft(input string, width int) string {
	r := []rune(input)

	if len(r) > width {
		return string(r[0:width])
	} else if len(r) < width {
		s := width - len(r)
		return string(r) + strings.Repeat(" ", s)
	}
	return input
}

// alignRight will format the table data to the right of the cell
// and will trim off characters in the event the output is too
// long.
func alignRight(input string, width int) string {
	r := []rune(input)

	if len(r) > width {
		return string(r[0:width])
	} else if len(r) < width {
		s := width - len(r)
		return strings.Repeat(" ", s) + string(r)
	}
	return input
}
