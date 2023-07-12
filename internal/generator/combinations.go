package generator

import (
	"log"
	"strings"

	"github.com/bxxf/tgen/internal/config"
	"github.com/bxxf/tgen/internal/rediscli"
	"github.com/bxxf/tgen/internal/utils"
)

func getBrandCountries(langs []string, flags map[string]string, includeEN bool) ([]string, bool, error) {
	isCountryFormat, countries, err := determineFormat(langs, flags)
	if err != nil {
		return nil, isCountryFormat, err
	}

	if includeEN {
		countries = append(countries, EN_COUNTRIES...)
	}
	countries = utils.RemoveDuplicates(countries)
	for i := range countries {
		countries[i] = strings.ToUpper(countries[i])
	}
	return countries, isCountryFormat, nil
}

func determineFormat(langs []string, flags map[string]string) (bool, []string, error) {
	redisCli := rediscli.NewUpstashClient(config.Config.URL, config.Config.Token)
	isCountryFormat := false
	formatFlag := flags["format"]
	if formatFlag != "" {
		if strings.ToLower(formatFlag) == "countryiso" {
			isCountryFormat = true
		}
		delete(flags, "format")
	}
	var countries []string
	if len(langs) > 0 {
		if strings.ToLower(langs[0]) == "avg" {
			isCountryFormat = true
		}

		countriesMapped, err := redisCli.GetCountries(strings.ToLower(langs[0]))
		log.Printf("cm: %v", countriesMapped)
		if err != nil {
			return isCountryFormat, countries, err
		}

		countries = countriesMapped
	}
	return isCountryFormat, countries, nil
}

func generateHeader(params map[string][]string, isCountryFormat bool) ([]string, []string) {
	paramKeys := utils.GetParamKeys(params)
	header := []string{"email", "locale"}
	if isCountryFormat {
		header[1] = "country_iso"
	}
	var relevantParamKeys []string
	for _, key := range paramKeys {
		value := params[key]
		// Exclude parameters with empty values
		if len(value) > 0 {
			header = append(header, key)
			relevantParamKeys = append(relevantParamKeys, key)
		}
	}
	return relevantParamKeys, header
}

func createRecords(header []string) [][]string {
	var records [][]string
	records = append(records, header)
	return records
}

func assembleRecords(resultsCh chan []string, records [][]string) [][]string {
	for record := range resultsCh {
		if record != nil {
			records = append(records, record)
		}
	}
	return records
}

func createCombinations(flags map[string]string, paramKeys []string, lang, country string, isCountryFormat bool, resultsCh chan<- []string, params map[string][]string) {
	iterateCombinations(0, []string{}, flags, paramKeys, func(comb []string) {
		email := constructEmail(lang, country, comb)
		record := []string{email, lang}
		if isCountryFormat {
			record[1] = country
		}
		record = append(record, comb...)

		resultsCh <- record
	}, params)
}

func iterateCombinations(idx int, current []string, flags map[string]string, paramKeys []string, callback func([]string), params map[string][]string) {
	if idx == len(paramKeys) {
		callback(current)
		return
	}

	key := paramKeys[idx]
	value := params[key]
	if len(value) > 0 {
		for _, v := range value {
			newCurrent := append(current, v)
			iterateCombinations(idx+1, newCurrent, flags, paramKeys, callback, params)
		}
	} else {
		iterateCombinations(idx+1, current, flags, paramKeys, callback, params)
	}
}

func constructEmail(lang, country string, paramValues []string) string {
	normalizedLang := strings.ReplaceAll(lang, "-", "_")
	params := strings.Join(paramValues, "_")
	delimiter := ""
	if params != "" {
		delimiter = "_"
	}
	removeWhitespace := strings.NewReplacer(" ", "_", "\n", "__", "\t", "__", "\r", "")
	params = removeWhitespace.Replace(params)

	var b strings.Builder
	b.WriteString("ttgen_")
	b.WriteString(strings.ToLower(normalizedLang))
	b.WriteString(delimiter)
	b.WriteString(params)
	b.WriteString("@example.com")
	return b.String()
}
