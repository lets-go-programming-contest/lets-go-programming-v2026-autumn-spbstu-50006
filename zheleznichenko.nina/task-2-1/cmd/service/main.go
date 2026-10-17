package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	minAllowedTemperature = 15
	maxAllowedTemperature = 30
	noSuitableTemperature = -1
)

type Operator string

const (
	OperatorLessThan     Operator = "<"
	OperatorGreaterThan  Operator = ">"
	OperatorLessEqual    Operator = "<="
	OperatorGreaterEqual Operator = ">="
)

var (
	errUnknownOperator = errors.New("unknown operator")
	errInvalidInput    = errors.New("invalid input")
)

type TemperatureRange struct {
	min int
	max int
}

func NewTemperatureRange() *TemperatureRange {
	return &TemperatureRange{
		min: minAllowedTemperature,
		max: maxAllowedTemperature,
	}
}

func (t *TemperatureRange) Update(operator Operator, threshold int) error {
	switch operator {
	case OperatorLessThan:
		if threshold-1 < t.max {
			t.max = threshold - 1
		}
	case OperatorGreaterThan:
		if threshold+1 > t.min {
			t.min = threshold + 1
		}
	case OperatorLessEqual:
		if threshold < t.max {
			t.max = threshold
		}
	case OperatorGreaterEqual:
		if threshold > t.min {
			t.min = threshold
		}
	default:
		return fmt.Errorf("%w: %q", errUnknownOperator, operator)
	}

	return nil
}

func (t *TemperatureRange) Value() int {
	if t.min > t.max {
		return noSuitableTemperature
	}

	return t.max
}

func parseCondition(line string) (Operator, int, error) {
	fields := strings.Fields(line)
	if len(fields) != 2 {
		return "", 0, fmt.Errorf("%w: expected operator and temperature", errInvalidInput)
	}

	operator := Operator(fields[0])
	switch operator {
	case OperatorLessThan, OperatorGreaterThan,
		OperatorLessEqual, OperatorGreaterEqual:
	default:
		return "", 0, fmt.Errorf("%w: %q", errUnknownOperator, fields[0])
	}

	threshold, err := strconv.Atoi(fields[1])
	if err != nil {
		return "", 0, fmt.Errorf("parse temperature: %w", err)
	}

	return operator, threshold, nil
}

func run(scanner *bufio.Scanner) error {
	if !scanner.Scan() {
		return fmt.Errorf("%w: missing department count", errInvalidInput)
	}

	departments, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		return fmt.Errorf("parse department count: %w", err)
	}
	if departments < 1 {
		return fmt.Errorf("%w: department count must be positive", errInvalidInput)
	}

	if !scanner.Scan() {
		return fmt.Errorf("%w: missing employee count", errInvalidInput)
	}

	employees, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		return fmt.Errorf("parse employee count: %w", err)
	}
	if employees < 1 {
		return fmt.Errorf("%w: employee count must be positive", errInvalidInput)
	}

	ranges := make([]*TemperatureRange, departments)
	for i := range ranges {
		ranges[i] = NewTemperatureRange()
	}

	for employee := 0; employee < employees; employee++ {
		for department := 0; department < departments; department++ {
			if !scanner.Scan() {
				return fmt.Errorf("%w: missing condition", errInvalidInput)
			}

			operator, threshold, err := parseCondition(scanner.Text())
			if err != nil {
				return err
			}

			if err := ranges[department].Update(operator, threshold); err != nil {
				return err
			}

			fmt.Println(ranges[department].Value())
		}
	}

	return scanner.Err()
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024), 1024*1024)

	if err := run(scanner); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
