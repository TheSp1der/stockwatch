package main

import (
	"bytes"
	"sort"
	"strconv"
	"strings"
	"time"

	"text/template"

	"github.com/TheSp1der/goerror"
	"github.com/fatih/color"
)

// displayTermnal returns a string for display in the terminal window of
// calculated and tracked stocks and the overall gains/losses of provided
// investments.
func displayTerminal(stock iex) string {
	var (
		err            error
		outputTemplate *template.Template
		data           outputStructure
		output         bytes.Buffer
		gtol           float64
	)

	// create the template
	tplt := ".-----------------------------------------------------------------------------.\n"
	tplt += "| Current Time: {{.CurrentTime}} Market Status: {{.MarketStatus}} |\n"
	tplt += "|--------------------------------.--------------.----------------.------------|\n"
	tplt += "| Company Name                   | Market Value | Today's Change | Gain/Loss  |\n"
	tplt += "|--------------------------------|--------------|----------------|------------|\n"
	tplt += "{{- range .Stock}}\n"
	tplt += "| {{ .CompanyName}} | {{.CurrentValue}} | {{.Change}} | {{.GL}} |\n"
	tplt += "{{- end }}\n"
	tplt += "{{- if .TotalGainLoss}}\n"
	tplt += "|--------------------------------'--------------'----------------'------------|\n"
	tplt += "| Total Investment Value: {{.TotalGainLoss}} |\n"
	tplt += "`-----------------------------------------------------------------------------'\n"
	tplt += "{{- else}}\n"
	tplt += "`--------------------------------'--------------'----------------'------------'\n"
	tplt += "{{- end}}"

	// initialize data stock map
	data.Stock = make(map[string]stockData)

	// sort stocks for display
	keys := make([]string, 0, len(stock))
	for k := range stock {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		var (
			cn   string  // company name
			cv   string  // current value
			ch   string  // change
			ival float64 // // investment value
			cval float64 // current value
			diff float64 // difference
			totl float64 // total difference
			t    string  // total investment (string output)
		)

		// calculate the total for the ticker in the event a stock
		// has multiple investments
		for _, i := range cmdLnInvestments {
			if strings.TrimSpace(strings.ToLower(stock[k].Company.Symbol)) == strings.TrimSpace(strings.ToLower(i.Ticker)) {
				ival = i.Quantity * i.Price
				cval = i.Quantity * stock[k].Price
				diff = cval - ival
				totl = totl + diff
			}
		}

		// update the grand total loss/gain
		gtol = gtol + totl

		// start setting values for template data struct
		cn = alignLeft(stock[k].Company.CompanyName, 30)
		cv = alignRight(strconv.FormatFloat(stock[k].Price, 'f', 2, 64), 12)
		if stock[k].Quote.Change < 0 {
			ch = color.RedString(alignRight(strconv.FormatFloat(stock[k].Quote.Change, 'f', 2, 64), 14))
		} else if stock[k].Quote.Change > 0 {
			ch = color.GreenString(alignRight(strconv.FormatFloat(stock[k].Quote.Change, 'f', 2, 64), 14))
		} else {
			ch = alignRight("", 14)
		}

		if totl < 0 {
			t = color.RedString(alignRight(strconv.FormatFloat(totl, 'f', 2, 64), 10))
		} else if totl > 0 {
			t = color.GreenString(alignRight(strconv.FormatFloat(totl, 'f', 2, 64), 10))
		} else {
			t = alignRight("", 10)
		}

		data.Stock[stock[k].Company.Symbol] = stockData{CompanyName: cn,
			CurrentValue: cv,
			Change:       ch,
			GL:           t,
		}
	}

	// set the date/time and market status
	data.CurrentTime = alignLeft(time.Now().Local().Format(timeFormat), 38)
	if m, _ := marketStatus(); m {
		data.MarketStatus = color.GreenString(alignLeft("OPEN", 7))
	} else {
		data.MarketStatus = color.YellowString(alignRight("CLOSED", 7))
	}
	if gtol < 0 {
		data.TotalGainLoss = color.RedString(alignRight(strconv.FormatFloat(gtol, 'f', 2, 64), 51))
	} else if gtol > 0 {
		data.TotalGainLoss = color.GreenString(alignRight(strconv.FormatFloat(gtol, 'f', 2, 64), 51))
	}

	outputTemplate = template.Must(template.New("console").Parse(tplt))

	if err = outputTemplate.Execute(&output, data); err != nil {
		goerror.Fatal(err)
	}

	return output.String()
}

// displayHTML returns a string for e-mail messages of calculated and
// tracked stocks and the overall gains/losses of provided investments.
func displayHTML(stockData iex) string {
	var (
		gtol   float64
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
			ival float64 // investment value
			cval float64 // current value
			diff float64 // difference
			totl float64 // total difference
		)

		// calculate the total for the ticker in the event a stock
		// has multiple investments
		for _, i := range cmdLnInvestments {
			if strings.TrimSpace(strings.ToLower(stockData[k].Company.Symbol)) == strings.TrimSpace(strings.ToLower(i.Ticker)) {
				ival = i.Quantity * i.Price
				cval = i.Quantity * stockData[k].Price
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
			output += "\t\t<td style=\"text-align: right; color: red;\">" + strconv.FormatFloat(totl, 'f', 2, 64) + "</td>\n"
		} else if totl > 0 {
			output += "\t\t<td style=\"text-align: right; color: green;\">" + strconv.FormatFloat(totl, 'f', 2, 64) + "</td>\n"
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
		output += "<span style=\"font-weight: bold;\">Overall Performance: <span style=\"color: red;\">" + strconv.FormatFloat(gtol, 'f', 2, 64) + "</span></span>\n"
	} else if gtol > 0 {
		output += "<span style=\"font-weight: bold;\">Overall Performance: <span style=\"color: green;\">" + strconv.FormatFloat(gtol, 'f', 2, 64) + "</span></span>\n"
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
