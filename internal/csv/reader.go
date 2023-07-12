package csv

import (
	"encoding/csv"
	"os"
)

func ReadCSVFile(filepath string) ([][]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'         // Set tab as the delimiter
	reader.FieldsPerRecord = -1 // Allow varying number of fields per record
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}
