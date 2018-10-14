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
	tplt += "| {{.CompanyName}} | {{.CurrentValue}} | {{.Change}} | {{.GL}} |\n"
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
			ival float64 // investment value
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
			Symbol:       strings.TrimSpace(strings.ToLower(stock[k].Company.Symbol)),
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
func displayHTML(stock iex) string {
	var (
		err            error
		outputTemplate *template.Template
		data           outputStructure
		output         bytes.Buffer
		gtol           float64
	)

	// create the template
	tplt := `<!doctype html>
<html xmlns="http://www.w3.org/1999/xhtml">
	<head>
		<title></title>"
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
	</head>
	<body>
		<span style="font-weight: bold;">Stock report as of {{.CurrentTime}}</span><br>
		<br>
		<table style="min-width: 700px;">
			<tr style="border-bottom: 4px solid gray;">
					<th style="text-align: left;">Company Name</th>
					<th style="text-align: right;">Market Value</th>
					<th style="text-align: right;">Today's Change</th>
				<th style="text-align: right;">Gain/Loss</th>
			</tr>
			{{- range .Stock}}
			<tr style="border-bottom: 1px solid gray;">
				<td style="text-align: left;">{{.CompanyName}}</td>
				<td style="text-align: right;">{{.CurrentValue}}</td>
				<td style="text-align: right;">{{.Change}}</td>
				<td style="text-align: right;">{{.GL}}</td>
			</tr>
			{{- end }}
		</table>
		<br>
		{{- if .TotalGainLoss}}
		<span style="font-weight: bold;">Overall Performance: {{.TotalGainLoss}}</span>
		{{- end}}
		<br>
		<br>
		{{- range .Stock}}
		<a href="https://finviz.com/quote.ashx?t={{.Symbol}}">{{.CompanyName}}</a><br>
		<img src="https://finviz.com/chart.ashx?t={{.Symbol}}"><br>
		{{- end}}
	</body>
</html>`

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
			ch = `<span style="color: red;">` + strconv.FormatFloat(stock[k].Quote.Change, 'f', 2, 64) + "</span>"
		} else if stock[k].Quote.Change > 0 {
			ch = `<span style="color: green;">` + strconv.FormatFloat(stock[k].Quote.Change, 'f', 2, 64) + "</span>"
		} else {
			ch = alignRight("", 14)
		}

		if totl < 0 {
			t = `<span style="color: red;">` + strconv.FormatFloat(totl, 'f', 2, 64) + "</span>"
		} else if totl > 0 {
			t = `<span style="color: green;">` + strconv.FormatFloat(totl, 'f', 2, 64) + "</span>"
		}

		data.Stock[stock[k].Company.Symbol] = stockData{
			CompanyName:  strings.TrimSpace(cn),
			CurrentValue: strings.TrimSpace(cv),
			Change:       ch,
			GL:           t,
			Symbol:       strings.TrimSpace(strings.ToLower(stock[k].Company.Symbol)),
		}
	}

	// set the date/time
	if m, _ := marketStatus(); m {
		data.CurrentTime = `<span style="color: green;">` + time.Now().Local().Format(timeFormat) + "</span>"
	} else {
		data.CurrentTime = `<span style="color: red;">` + time.Now().Local().Format(timeFormat) + "</span>"
	}

	// overall total/loss
	if gtol < 0 {
		data.TotalGainLoss = `<span style="color: red;">` + strconv.FormatFloat(gtol, 'f', 2, 64) + "</span>"
	} else if gtol > 0 {
		data.TotalGainLoss = `<span style="color: green;">` + strconv.FormatFloat(gtol, 'f', 2, 64) + "</span>"
	}

	outputTemplate = template.Must(template.New("console").Parse(tplt))

	if err = outputTemplate.Execute(&output, data); err != nil {
		goerror.Fatal(err)
	}

	return output.String()
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

func displayWeb(stock iex) string {
	var (
		err            error
		outputTemplate *template.Template
		data           outputStructure
		output         bytes.Buffer
		gtol           float64
	)

	// create the template
	tplt := `
	<!doctype html>
	<html lang="en">
		<head>
			<!-- Required meta tags -->
			<meta charset="utf-8">
			<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
			<meta http-equiv="refresh" content="5">
			<!-- Bootstrap CSS -->
			<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.1.3/css/bootstrap.min.css" crossorigin="anonymous">
			<title>Stockwatch</title>
		</head>
		<body>
			<main role="main">
				<div class="jumbotron">
					<div class="container">
						<h1 class="display-5">Stockwatch</h1>
						<p class="lead">A simple stock tracker written in Go</p>
					</div>
				</div>
				<div class="container">
					<table class="table-sm table-striped mx-auto">
						<thead>
							<tr>
								<th>Company Name</th>
								<th class="text-right">Market Value</th>
								<th class="text-right">Today's Change</th>
								<th class="text-right">Gain/Loss</th>
							</tr>
						</thead>
						<tbody>
							{{- range .Stock}}
							<tr>
								<td>{{.CompanyName}}</td>
								<td class="text-right">{{.CurrentValue}}</td>
								<td class="text-right">{{.Change}}</td>
								<td class="text-right">{{.GL}}</td>
							</tr>
							{{- end }}
							{{- if .TotalGainLoss}}
							<tr>
								<td colspan="2" class="text-left font-weight-bold">Overall Performance</td>
								<td colspan="2" class="text-right">{{.TotalGainLoss}}</td>
							</tr>
							{{- end}}
						</tbody>
					</table>
				</div>
			</main>
			<footer class="container mt-2">
				<hr>
			</footer>
			<script src="https://code.jquery.com/jquery-3.3.1.slim.min.js" crossorigin="anonymous"></script>
			<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.3/umd/popper.min.js" crossorigin="anonymous"></script>
			<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.1.3/js/bootstrap.min.js" crossorigin="anonymous"></script>
		</body>
</html>`

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
			ch = `<span style="color: red;">` + strconv.FormatFloat(stock[k].Quote.Change, 'f', 2, 64) + "</span>"
		} else if stock[k].Quote.Change > 0 {
			ch = `<span style="color: green;">` + strconv.FormatFloat(stock[k].Quote.Change, 'f', 2, 64) + "</span>"
		} else {
			ch = alignRight("", 14)
		}

		if totl < 0 {
			t = `<span style="color: red;">` + strconv.FormatFloat(totl, 'f', 2, 64) + "</span>"
		} else if totl > 0 {
			t = `<span style="color: green;">` + strconv.FormatFloat(totl, 'f', 2, 64) + "</span>"
		}

		data.Stock[stock[k].Company.Symbol] = stockData{
			CompanyName:  strings.TrimSpace(cn),
			CurrentValue: strings.TrimSpace(cv),
			Change:       ch,
			GL:           t,
			Symbol:       strings.TrimSpace(strings.ToLower(stock[k].Company.Symbol)),
		}
	}

	// set the date/time
	if m, _ := marketStatus(); m {
		data.CurrentTime = `<span style="color: green;">` + time.Now().Local().Format(timeFormat) + "</span>"
	} else {
		data.CurrentTime = `<span style="color: red;">` + time.Now().Local().Format(timeFormat) + "</span>"
	}

	// overall total/loss
	if gtol < 0 {
		data.TotalGainLoss = `<span style="color: red;">` + strconv.FormatFloat(gtol, 'f', 2, 64) + "</span>"
	} else if gtol > 0 {
		data.TotalGainLoss = `<span style="color: green;">` + strconv.FormatFloat(gtol, 'f', 2, 64) + "</span>"
	}

	outputTemplate = template.Must(template.New("console").Parse(tplt))

	if err = outputTemplate.Execute(&output, data); err != nil {
		goerror.Fatal(err)
	}

	return output.String()
}
