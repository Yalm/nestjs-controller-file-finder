package utils

import (
	"log"
	"os"
	"regexp"
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

func ExtracParamNames(route string) []string {
	paramNames := make([]string, 0)

	routePattern := regexp.MustCompile(`\{(\w+)\}`)

	matches := routePattern.FindAllStringSubmatch(route, -1)
	if matches == nil {
		return paramNames
	}

	for _, match := range matches {
		paramNames = append(paramNames, match[1])
	}

	return paramNames
}

func ConvertParamNamesToMapping(paramNames []string) map[string]string {
	paramMapping := make(map[string]string)

	for _, paramName := range paramNames {
		paramMapping["integration.request.path."+paramName] = "method.request.path." + paramName
	}

	return paramMapping
}

func ConvertParamNamesToMappingWithPrefix(paramNames []string, prefix string) map[string]bool {
	paramMapping := make(map[string]bool)

	for _, paramName := range paramNames {
		paramMapping[prefix+paramName] = true
	}

	return paramMapping
}
