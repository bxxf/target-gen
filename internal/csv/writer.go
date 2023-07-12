package csv

import (
	"encoding/csv"
	"os"
)

func WriteToCsv(records [][]string, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = '\t'

	if err := writer.WriteAll(records); err != nil {
		return err
	}

	return nil
}
