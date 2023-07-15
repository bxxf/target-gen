package generator

import (
	"fmt"
	"log"
	"sync"

	"github.com/bxxf/tgen/internal/constants"
)

func Generate(languages []string, flags map[string]string, parameters map[string][]string) ([][]string, error) {
	langCountryMapping := constants.CountryToLocale
	enAll := flags["en-all"] == "true"
	countries, isAvg, err := getBrandCountries(languages, flags, enAll)
	if err != nil {
		return nil, err
	}

	params, header := generateHeader(parameters, isAvg)
	records := createRecords(header)

	var wg sync.WaitGroup
	resultCh := make(chan []string)
	errorCh := make(chan error)

	for _, country := range countries {
		wg.Add(1)
		go func(country string) {
			defer wg.Done()
			lang, ok := langCountryMapping[country]
			if !ok {
				errorCh <- fmt.Errorf("Warning: language mapping not found for country %s", country)
				resultCh <- nil
				return
			}
			createCombinations(flags, params, lang, country, isAvg, resultCh, parameters)
		}(country)
	}
	go func() {
		wg.Wait()
		close(resultCh)
	}()
	go func() {
		for err := range errorCh {
			log.Println(err)
		}
	}()
	records = assembleRecords(resultCh, records)
	records = removeDuplicateRecords(records)

	log.Printf("Successfully generated %v records with %v languages.", len(records)-1, len(countries))

	return records, nil
}

func removeDuplicateRecords(records [][]string) [][]string {
	seen := make(map[string][]string)
	var unique [][]string

	for _, record := range records {
		key := fmt.Sprint(record)
		if _, ok := seen[key]; !ok {
			seen[key] = record
			unique = append(unique, record)
		}
	}
	return unique
}
