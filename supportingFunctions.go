package main

import (
	"os"
	"strconv"
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

func getEnvFloat64(env string, def float64) (ret float64) {
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
