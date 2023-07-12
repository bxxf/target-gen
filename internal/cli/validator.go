package cli

import (
	"errors"
)

func validate(parameters map[string][]string) error {
	// add more validation logic here
	languages := parameters["languages"]
	locFile := parameters["--loc-file"]

	if (languages == nil || len(languages) == 0) && (locFile == nil || len(locFile) == 0) {
		return errors.New("languages parameter or loc file is required")
	}

	return nil
}
