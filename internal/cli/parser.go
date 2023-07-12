package cli

import (
	"strings"
)

func ParseArgs(args []string) (map[string][]string, error) {
	parameters := make(map[string][]string)
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 2 {
			parameters[parts[0]] = strings.Split(parts[1], ",")
		} else {
			parameters[arg] = []string{}
		}
	}
	return parameters, validate(parameters)
}
