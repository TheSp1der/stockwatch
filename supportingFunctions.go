package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// getEnvString returns string from environment variable.
func getEnvString(env string, def string) (val string) {
	val = os.Getenv(env)

	if len(val) == 0 {
		return def
	}

	return
}

// getEnvBool returns boolean from environment variable.
func getEnvBool(env string, def bool) (ret bool) {
	val := os.Getenv(env)

	if len(val) == 0 {
		return def
	}

	ret, err := strconv.ParseBool(val)
	if err != nil {
		log.Fatal(val + " environment variable is not boolean")
	}

	return
}

// getEnvInt returns int from environment variable.
func getEnvInt(env string, def int) (ret int) {
	val := os.Getenv(env)

	if len(val) == 0 {
		return def
	}

	ret, err := strconv.Atoi(val)
	if err != nil {
		log.Fatal(env + " environment variable is not numeric")
	}

	return
}

func getFloat64(env string, def float64) (ret float64) {
	val := os.Getenv(env)

	if len(val) == 0 {
		return def
	}

	ret, err := strconv.ParseFloat(val, 64)
	if err != nil {
		log.Fatal(env + " environment variable is not floating point.")
	}

	return
}

// remove duplicate values from slice
func uniqueString(inputString []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range inputString {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// String format flag value.
func (i *configInvestments) String() string {
	return fmt.Sprint(*i)
}

// Set set flag value.
func (i *configInvestments) Set(value string) error {
	if len(strings.Split(value, ",")) == 3 {
		var (
			err      error
			quantity float64
			price    float64
		)

		inv := strings.Split(value, ",")
		if quantity, err = strconv.ParseFloat(inv[1], 32); err != nil {
			return err
		}
		if price, err = strconv.ParseFloat(inv[2], 32); err != nil {
			return err
		}
		stockwatchConfig.Investments = append(stockwatchConfig.Investments, configInvestment{
			Ticker:   inv[0],
			Quantity: quantity,
			Price:    price,
		})
	}
	return nil
}
