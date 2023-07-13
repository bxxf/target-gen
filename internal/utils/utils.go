package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bxxf/tgen/internal/csv"
)

// ParseArgs parses command line arguments and returns a map of parameters.
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

var CountryToLocale = map[string]string{
	"US": "en-US",
	"EN": "en-US",
	"CA": "en-CA",
	"AU": "en-AU",
	"GB": "en-GB",
	"FI": "fi-FI",
	"PL": "pl-PL",
	"BR": "pt-BR",
	"IT": "it-IT",
	"NO": "no-NO",
	"DK": "da-DK",
	"NL": "nl-NL",
	"FR": "fr-FR",
	"DE": "de-DE",
	"ES": "es-ES",
	"SE": "sv-SE",
	"CN": "zh-CN",
	"TW": "zh-TW",
	"GR": "el-GR",
	"IL": "he-IL",
	"HU": "hu-HU",
	"ID": "id-ID",
	"JP": "ja-JP",
	"KR": "ko-KR",
	"MY": "ms-MY",
	"PT": "pt-PT",
	"RU": "ru-RU",
	"SK": "sk-SK",
	"TH": "th-TH",
	"TR": "tr-TR",
	"UA": "uk-UA",
	"VN": "vi-VN",
	"SA": "ar-SA",
	"MX": "es-MX",
	"CL": "es-CL",
	"CO": "es-CO",
	"PE": "es-PE",
	"VE": "es-VE",
	"ZA": "en-ZA",
	"IN": "en-IN",
	"RS": "sr-RS",
	"CZ": "cs-CZ",
	"SV": "sv-SE",
	"DA": "da-DK",
	"ZH": "zh-CN",
	"EL": "el-GR",
	"HE": "he-IL",
	"JA": "ja-JP",
	"KO": "ko-KR",
	"MS": "ms-MY",
	"RO": "ro-RO",
	"UK": "uk-UA",
	"VI": "vi-VN",
	"AR": "ar-SA",
	"BG": "bg-BG",
	"HR": "hr-HR",
	"CS": "cs-CZ",
}

func GetLangCountryMapping() map[string]string {
	return CountryToLocale
}

func ConvertCountryToLocale(countries []string) []string {
	locales := make([]string, len(countries))

	for i, country := range countries {
		if locale, ok := CountryToLocale[country]; ok {
			locales[i] = locale
		} else {
			log.Fatal("Error: Country code not found: " + country)
		}
	}

	return locales
}

// Contains checks if a string is present in a slice.
func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

// CheckError prints an error message and exits if the error is not nil.
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

func GetLanguagesFromLocFile(locFilePath string) ([]string, error) {
	fileData, err := csv.ReadCSVFile(locFilePath)
	if err != nil {
		if !strings.HasSuffix(locFilePath, ".csv") {
			csvFilePath := locFilePath + ".csv"
			fileData, err = csv.ReadCSVFile(csvFilePath)
		}
		if err != nil {
			if !strings.HasSuffix(locFilePath, ".txt") && !strings.HasSuffix(locFilePath, ".csv") {
				txtFilePath := locFilePath + ".txt"
				fileData, err = csv.ReadCSVFile(txtFilePath)
			}
			if err != nil {
				return nil, err
			}
		}
	}

	languages := make([]string, 0)

	// Extract the language codes from the first column of fileData
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
