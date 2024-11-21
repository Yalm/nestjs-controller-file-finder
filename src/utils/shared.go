package utils

import (
	"log"
	"os"
	"strconv"
)

func Getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func GetIntenv(key, fallback string) int {
	value := Getenv(key, fallback)
	convertedValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Failed to convert %s to interger: %s", key, err)
	}
	return convertedValue
}

func GetBoolenv(key, fallback string) bool {
	value := Getenv(key, fallback)
	convertedValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Fatalf("Failed to convert %s to boolean: %s", key, err)
	}
	return convertedValue
}
