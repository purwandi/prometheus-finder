package finder

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/purwandi/prometheus-finder/config"
	"github.com/sirupsen/logrus"
)

type Result struct {
	Found bool
	Line  int
}

func ReadFile(path, query string, mode config.FinderMode) (Result, error) {
	var result Result

	file, err := os.Open(path)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	length := 0
	for i := 1; scanner.Scan(); i++ {
		read_line := strings.TrimSpace(scanner.Text())

		if strings.Contains(read_line, query) {
			result.Found = true
			result.Line = i
		}

		length = i
	}

	if err := scanner.Err(); err != nil {
		return result, err
	}

	if mode == config.ModeFirst {
		if result.Found && result.Line == 1 {
			return result, nil
		}
		return result, fmt.Errorf("finder is not found in first line")
	}

	if mode == config.ModeLast {
		if result.Found && result.Line == length {
			return result, nil
		}
		return result, fmt.Errorf("finder is not found in last line")
	}

	if result.Found {
		return result, nil
	}

	return Result{}, errors.New("not found")
}
