package generator

import (
	"fmt"
	"strings"

	"github.com/bxxf/target-gen/internal/csv"
)

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
