package generator

import (
	"fmt"
	"log"
	"sync"

	"github.com/bxxf/target-gen/internal/utils"
)

func Generate(languages []string, flags map[string]string, parameters map[string][]string) ([][]string, error) {
	langCountryMapping := utils.GetLangCountryMapping()
	enAll := flags["en-all"] == "true"
	countries, isAvg, err := parseCountries(languages, flags, enAll)
	if err != nil {
		return nil, err
	}

	params, header := parseHeader(parameters, isAvg)
	records := initializeRecords(header)

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
			generateCombinations(flags, params, lang, country, isAvg, resultCh, parameters)
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
	records = collectRecords(resultCh, records)
	log.Printf("Successfully generated %v records with %v languages.", len(records)-1, len(languages))

	return records, nil
}
