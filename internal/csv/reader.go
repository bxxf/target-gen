package csv

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
)

func ReadCSVFile(filepath string) ([][]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	estimatedRecords := int(fileInfo.Size()) / 50

	records := make([][]string, 0, estimatedRecords)
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}
