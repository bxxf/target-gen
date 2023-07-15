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

func ParseArgs(args []string) map[string][]string {
	parameters := make(map[string][]string, len(args))
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 2 {
			parameters[parts[0]] = strings.Split(parts[1], ",")
		} else {
			parameters[arg] = []string{}
		}
	}
	return parameters
}

func GetParamKeys(parameters map[string][]string) []string {
	keys := make([]string, 0, len(parameters))
	for k := range parameters {
		keys = append(keys, k)
	}
	return keys
}

func ConvertCountryToLocale(countries []string) []string {
	locales := make([]string, len(countries))

	for i, country := range countries {
		if locale, ok := constants.CountryToLocale[country]; ok {
			locales[i] = locale
		} else {
			log.Fatal("Error: Country code not found: " + country)
		}
	}

	return locales
}

func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func CheckError(msg string, err error) {
	if err != nil {
		fmt.Println(msg+":", err)
		os.Exit(1)
	}
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

func GenerateEmail(lang, country string, paramValues []string) string {
	normalizedLang := strings.ReplaceAll(lang, "-", "_")
	params := strings.Join(paramValues, "_")
	delimiter := ""
	if params != "" {
		delimiter = "_"
	}
	return fmt.Sprintf("example_%s%s%s@example.com", strings.ToLower(normalizedLang), delimiter, params)
}
