package main

import (
	"testing"
)

func TestValidateStartParameters(t *testing.T) {
	w := Configuration{
		IexAPIKey: "",
		HTTPPort:  0,
		Mail: configMail{
			Port: 0,
		},
	}

	e := validateStartParameters(w, "")
	if e.Error() != "IEX API Key was not provided" {
		t.Errorf("Error expected and not detected, got %v", e)
	}

	w.IexAPIKey = "The quick fox jumped over the lazy dog."
	e = validateStartParameters(w, "")
	if e.Error() != "No stocks were defined" {
		t.Errorf("Error expected and not detected, got %v", e)
	}

	e = validateStartParameters(w, "example")
	if e.Error() != "Invalid port given for http server provided" {
		t.Errorf("Error expected and not detected, got %v", e)
	}

	w.HTTPPort = 1
	e = validateStartParameters(w, "example")
	if e.Error() != "Invalid port for mail server provided" {
		t.Errorf("Error expected and not detected, got %v", e)
	}

	w.Mail.Port = 1
	e = validateStartParameters(w, "example")
	if e != nil {
		t.Errorf("No error expected but got one, got %v", e)
	}
}

func TestPrepareStartParameters(t *testing.T) {
	var (
		e error
		w = Configuration{
			Investments: configInvestments{
				configInvestment{
					Ticker:   "amd",
					Quantity: 1,
					Price:    1,
				}, {
					Ticker:   "ftd",
					Quantity: 1,
					Price:    1,
				},
			},
		}
		s = "amd,ati,googl*"
	)

	w, e = prepareStartParameters(w, s)
	if e.Error() != "Stock ticker googl* is not a valid ticker name" {
		t.Errorf("Error expected and not detected, got %v", e)
	}

	s = "amd,ati,googl"
	w, e = prepareStartParameters(w, s)
	if w.TrackedStocks[0] != "amd" || w.TrackedStocks[1] != "ati" || w.TrackedStocks[2] != "googl" || w.TrackedStocks[3] != "ftd" || len(w.TrackedStocks) != 4 {
		t.Errorf("Expected to get specific list of stocks, wanted %v, got %v", [4]string{"amd", "ati", "googl", "ftd"}, w.TrackedStocks)
	}
}
