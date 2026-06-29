package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Yalm/nestjs-controller-file-finder/src/utils"
)

func FindRoutes(rootDir string, ignoredPaths map[string]bool) ([]utils.Route, error) {

	log.Println("Searching for routes in", rootDir)

	var extractedRoutes []utils.Route

	controllerRegex := regexp.MustCompile(`@Controller\(\s*(\[[^\]]*\]|['"][^'"]*['"])?\s*\)`)
	controllerValueRegex := regexp.MustCompile(`['"]([^'"]*)['"]`)
	methodRegex := regexp.MustCompile(`@(Get|Post|Put|Delete|Patch)\(['"]?([^'"]*)['"]?\)`)
	pathRegex := regexp.MustCompile(`:(\w+)`)

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".controller.ts") {
			extractRoutesFromFile(path, controllerRegex, controllerValueRegex, methodRegex, pathRegex, ignoredPaths, &extractedRoutes)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return extractedRoutes, nil
}

func extractRoutesFromFile(
	filePath string,
	controllerRegex, controllerValueRegex, methodRegex, pathRegex *regexp.Regexp,
	ignoredPaths map[string]bool,
	routes *[]utils.Route) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Error opening file", filePath, ":", err)
		return
	}
	defer file.Close()

	paths := make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if matches := controllerRegex.FindStringSubmatch(line); len(matches) > 1 {
			if matches[1] == "" {
				paths[""] = true
				continue
			}
			values := controllerValueRegex.FindAllStringSubmatch(matches[1], -1)
			for _, v := range values {
				paths[v[1]] = true
			}
		}

		if matches := methodRegex.FindStringSubmatch(line); len(matches) > 2 {
			methodRoute := matches[2]
			for baseRoute := range paths {
				fullRoute := filepath.Join("/", baseRoute, methodRoute)
				fullRoute = strings.ReplaceAll(fullRoute, "\\", "/")

				if ignoredPaths[fullRoute] {
					log.Println("Ignoring route", fullRoute)
					continue
				}

				fullRoute = pathRegex.ReplaceAllString(fullRoute, `{$1}`)
				*routes = append(*routes, utils.NewRoute(strings.ToUpper(matches[1]), fullRoute))
			}

		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error reading file", filePath, ":", err)
	}
}
