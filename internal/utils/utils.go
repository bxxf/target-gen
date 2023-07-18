package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bxxf/tgen/internal/constants"
	"github.com/bxxf/tgen/internal/csv"
)


func GetParamKeys(parameters map[string][]string) []string {
	keys := make([]string, 0, len(parameters))
	for k := range parameters {
		keys = append(keys, k)
	}
	return keys
}


func RemoveDuplicates(s []string) []string {
	unique := make([]string, 0, len(s))
	seen := make(map[string]struct{})

	for _, item := range s {
		if _, ok := seen[item]; !ok {
			unique = append(unique, item)
			seen[item] = struct{}{}
		}
	}
	return unique
}

func ParseParams(params []string) map[string][]string {
	parameters := make(map[string][]string, len(params))
	for _, param := range params {
		parts := strings.SplitN(param, "=", 2)
		if len(parts) == 2 {
			parameters[parts[0]] = strings.Split(parts[1], ",")
		} else {
			parameters[param] = []string{}
		}
	}
	return parameters
}

func GenerateFileName(output string) string {
	timestamp := time.Now().Format("20060102150405")

	filename := "tgen-" + timestamp + ".csv"
	if output != "" {
		filename = output
		if !strings.HasSuffix(filename, ".csv") {
			filename += ".csv"
		}
	}

	return filename
}

func GetLanguagesFromLocFile(locFilePath string) ([]string, error) {
	fileData, err := csv.ReadCSVFile(locFilePath)
	if err != nil {
		fmt.Printf("Error while reading LOC FILE: %s", err)
		return nil, err
	}

	languages := make([]string, 0)
	for _, row := range fileData {
		if len(row) < 1 {
			continue
		}
		langCode := row[0]
		if len(langCode) == 2 {
			languages = append(languages, langCode)
		}
	}

	return languages, nil
}
