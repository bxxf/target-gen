package csv

import (
	"bufio"
	"encoding/csv"
	"os"
)

func WriteToCsv(records [][]string, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(bufio.NewWriter(file))
	writer.Comma = '\t'

	writer.WriteAll(records[:0]) // Preallocate capacity by writing an empty slice

	if err := writer.WriteAll(records); err != nil {
		return err
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return err
	}

	return nil
}
