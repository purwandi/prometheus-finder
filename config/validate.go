package config

import (
	"fmt"
)

func isModeisValid(m FinderMode) bool {
	switch m {
	case ModeFull, ModeFirst, ModeLast:
		return true
	default:
		return false
	}
}

func IsOk(configs []Config) (bool, error) {
	ok := true
	for _, item := range configs {
		if item.Query == "" {
			return false, fmt.Errorf("query string can't be empty at metric %s", item.Name)
		}

		if ok = isModeisValid(item.Mode); !ok {
			return ok, fmt.Errorf("query mode is must be one of (full_search, first_of_line, last_of_line) at metric %s", item.Name)
		}
	}

	return ok, nil
}
