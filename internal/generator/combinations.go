package generator

import (
	"fmt"
	"strings"

	"github.com/bxxf/tgen/internal/utils"
)

func parseCountries(languages []string, flags map[string]string, enAll bool) ([]string, bool, error) {
	countryFormat, brandCountries, err := handleFormat(languages, flags)
	if err != nil {
		return nil, countryFormat, err
	}
	countries := brandCountries

	if enAll {
		countries = append(countries, EN_COUNTRIES...)
	}
	countries = utils.RemoveDuplicates(countries)
	for i := range countries {
		countries[i] = strings.ToUpper(countries[i])
	}
	return countries, countryFormat, nil
}

func handleFormat(languages []string, flags map[string]string) (bool, []string, error) {
	countryFormat := false
	format := flags["format"]
	if format != "" {
		if strings.ToLower(format) == "countryiso" {
			countryFormat = true
		}
		delete(flags, "format")
	}
	var err error
	var brandCountries []string
	if len(languages) > 0 {
		if strings.ToLower(languages[0]) == Avg {
			countryFormat = true
		}
		brandCountriesK, ok := BrandToCountries[languages[0]]
		if !ok {
			brandCountriesK = languages
		}

		brandCountries = brandCountriesK

	}
	return countryFormat, brandCountries, err
}

func parseHeader(parameters map[string][]string, countryFormat bool) ([]string, []string) {
	paramKeys := utils.GetParamKeys(parameters)
	header := []string{"email", "locale"}
	if countryFormat {
		header[1] = "country_iso"
	}
	var dynamicParamKeys []string
	for _, key := range paramKeys {
		value := parameters[key]
		// Exclude parameters with empty values
		if len(value) > 0 {
			header = append(header, key)
			dynamicParamKeys = append(dynamicParamKeys, key)
		}
	}
	return dynamicParamKeys, header
}

func initializeRecords(header []string) [][]string {
	var records [][]string
	records = append(records, header)
	return records
}

func collectRecords(resultCh chan []string, records [][]string) [][]string {
	for record := range resultCh {
		if record != nil {
			records = append(records, record)
		}
	}
	return records
}
func generateCombinations(flags map[string]string, paramKeys []string, lang, country string, countryFormat bool, resultCh chan<- []string, parameters map[string][]string) {
	combinationGenerator(0, []string{}, flags, paramKeys, func(comb []string) {
		email := generateEmail(lang, country, comb)
		record := []string{email, lang}
		if countryFormat {
			record[1] = country
		}
		record = append(record, comb...)

		resultCh <- record
	}, parameters)
}

func combinationGenerator(index int, current []string, flags map[string]string, paramKeys []string, callback func([]string), parameters map[string][]string) {
	if index == len(paramKeys) {
		callback(current)
		return
	}

	key := paramKeys[index]
	value := parameters[key]
	if len(value) > 0 {
		for _, v := range value {
			newCurrent := append(current, v)
			combinationGenerator(index+1, newCurrent, flags, paramKeys, callback, parameters)
		}
	} else {
		combinationGenerator(index+1, current, flags, paramKeys, callback, parameters)
	}
}

func generateEmail(lang, country string, paramValues []string) string {
	normalizedLang := strings.ReplaceAll(lang, "-", "_")
	params := strings.Join(paramValues, "_")
	delimiter := ""
	if params != "" {
		delimiter = "_"
	}
	return fmt.Sprintf("ttgen_%s%s%s@example.com", strings.ToLower(normalizedLang), delimiter, params)
}
