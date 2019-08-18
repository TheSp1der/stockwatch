package main

import (
	"os"
	"strconv"
	"testing"
)

func TestGetEnvString(t *testing.T) {
	exp := "Good"
	got := "Bad"

	os.Setenv("TEST", exp)

	res := getEnvString("TEST", got)
	if res == got {
		t.Errorf("Expected return invalid, expected %v, got %v", exp, res)
	}
}

func TestGetEnvBool(t *testing.T) {
	exp := true
	got := false

	os.Setenv("TEST", strconv.FormatBool(exp))

	res := getEnvBool("TEST", got)
	if res == got {
		t.Errorf("Expected return invalid, expected %v, got %v", exp, res)
	}
}

func TestGetEnvInt(t *testing.T) {
	exp := 1
	got := 0

	os.Setenv("TEST", strconv.Itoa(exp))

	res := getEnvInt("TEST", got)
	if res == got {
		t.Errorf("Expected return invalid, expected %v, got %v", exp, res)
	}
}

func TestGetEnvFloat64(t *testing.T) {
	exp := float64(1)
	got := float64(0)

	os.Setenv("TEST", strconv.FormatFloat(exp, 'f', 0, 64))

	res := getEnvFloat64("TEST", got)
	if res == got {
		t.Errorf("Expected return invalid, expected %v, got %v", exp, res)
	}
}

func TestUniqueString(t *testing.T) {
	res := uniqueString([]string{"test", "test", "test", "test123"})

	if res[0] != "test" || res[1] != "test123" || len(res) != 2 {
		t.Errorf("Un-expected result form uniqueString, got %v", res)
	}
}
