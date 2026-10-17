package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

const (
	minAllowedTemperature = 15
	maxAllowedTemperature = 30
	impossibleTemperature = -1
)

var ErrUnknownOperator = errors.New("unknown operator")

type Operator string

const (
	operatorGreaterOrEqual Operator = ">="
	operatorLessOrEqual    Operator = "<="
	operatorEqual          Operator = "="
)

type TemperatureRange struct {
	min int
	max int
}

func NewTemperatureRange() TemperatureRange {
	return TemperatureRange{
		min: minAllowedTemperature,
		max: maxAllowedTemperature,
	}
}

func (r *TemperatureRange) Apply(operator Operator, threshold int) error {
	switch operator {
	case operatorGreaterOrEqual:
		r.min = max(r.min, threshold)
	case operatorLessOrEqual:
		r.max = min(r.max, threshold)
	case operatorEqual:
		r.min = max(r.min, threshold)
		r.max = min(r.max, threshold)
	default:
		return fmt.Errorf("%w: %q", ErrUnknownOperator, operator)
	}

	return nil
}

func (r *TemperatureRange) Result() int {
	if r.min > r.max {
		return impossibleTemperature
	}

	return r.min
}

func readToken(scanner *bufio.Scanner) (string, error) {
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("scan input: %w", err)
		}

		return "", fmt.Errorf("read input token: %w", io.EOF)
	}

	return scanner.Text(), nil
}

func readInt(scanner *bufio.Scanner) (int, error) {
	token, err := readToken(scanner)
	if err != nil {
		return 0, err
	}

	value, err := strconv.Atoi(token)
	if err != nil {
		return 0, fmt.Errorf("parse integer %q: %w", token, err)
	}

	return value, nil
}

func readOperator(scanner *bufio.Scanner) (Operator, error) {
	token, err := readToken(scanner)
	if err != nil {
		return "", err
	}

	operator := Operator(token)

	switch operator {
	case operatorGreaterOrEqual, operatorLessOrEqual, operatorEqual:
		return operator, nil
	default:
		return "", fmt.Errorf("%w: %q", ErrUnknownOperator, token)
	}
}

func processTestCases(
	scanner *bufio.Scanner,
	testCasesCount int,
) error {
	for range testCasesCount {
		operationsCount, err := readInt(scanner)
		if err != nil {
			return fmt.Errorf("read operations count: %w", err)
		}

		temperatureRange := NewTemperatureRange()

		for range operationsCount {
			operator, err := readOperator(scanner)
			if err != nil {
				return fmt.Errorf("read operator: %w", err)
			}

			threshold, err := readInt(scanner)
			if err != nil {
				return fmt.Errorf("read temperature: %w", err)
			}

			if err := temperatureRange.Apply(operator, threshold); err != nil {
				return err
			}

			fmt.Println(temperatureRange.Result())
		}
	}

	return nil
}

func run() error {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	testCasesCount, err := readInt(scanner)
	if err != nil {
		return fmt.Errorf("read test cases count: %w", err)
	}

	if err := processTestCases(scanner, testCasesCount); err != nil {
		return fmt.Errorf("process test cases: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
