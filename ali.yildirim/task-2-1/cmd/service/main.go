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
	MinTemp    = 15
	MaxTemp    = 30
	NoSolution = -1
)

type Operator string

const (
	OpGreaterOrEqual Operator = ">="
	OpLessOrEqual    Operator = "<="
)

var (
	ErrUnknownOperator = errors.New("unknown operator")
	ErrUnexpectedEOF   = errors.New("unexpected end of input")
)

type TemperatureRange struct {
	low  int
	high int
}

func NewTemperatureRange() *TemperatureRange {
	return &TemperatureRange{
		low:  MinTemp,
		high: MaxTemp,
	}
}

func (temperature *TemperatureRange) Update(operator Operator, val int) error {
	switch operator {
	case OpGreaterOrEqual:
		if val > temperature.low {
			temperature.low = val
		}
	case OpLessOrEqual:
		if val < temperature.high {
			temperature.high = val
		}
	default:
		return fmt.Errorf("%w: %q", ErrUnknownOperator, operator)
	}

	return nil
}

func (temperature *TemperatureRange) Optimal() int {
	if temperature.low <= temperature.high {
		return temperature.low
	}

	return NoSolution
}

type TokenReader struct {
	scanner *bufio.Scanner
}

func NewTokenReader(reader io.Reader) *TokenReader {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)

	return &TokenReader{scanner: scanner}
}

func (reader *TokenReader) Next() (string, error) {
	if !reader.scanner.Scan() {
		if err := reader.scanner.Err(); err != nil {
			return "", fmt.Errorf("scanner error: %w", err)
		}

		return "", ErrUnexpectedEOF
	}

	return reader.scanner.Text(), nil
}

func (reader *TokenReader) NextInt() (int, error) {
	token, err := reader.Next()
	if err != nil {
		return 0, err
	}

	val, err := strconv.Atoi(token)
	if err != nil {
		return 0, fmt.Errorf("invalid integer format %q: %w", token, err)
	}

	return val, nil
}

func processEmployee(reader *TokenReader, rangeTracker *TemperatureRange) (int, error) {
	opStr, err := reader.Next()
	if err != nil {
		return NoSolution, err
	}

	val, err := reader.NextInt()
	if err != nil {
		return NoSolution, err
	}

	if err := rangeTracker.Update(Operator(opStr), val); err != nil {
		return NoSolution, err // DÜZELTME: Hata varsa hatayı dön
	}

	return rangeTracker.Optimal(), nil
}

func processDepartment(reader *TokenReader, out io.Writer) error {
	employees, err := reader.NextInt()
	if err != nil {
		return err
	}

	rangeTracker := NewTemperatureRange()

	for range employees {
		res, empErr := processEmployee(reader, rangeTracker)
		if empErr != nil {
			return empErr
		}

		if _, printErr := fmt.Fprintln(out, res); printErr != nil {
			return fmt.Errorf("print result: %w", printErr)
		}
	}

	return nil
}

func RunService(in io.Reader, out io.Writer) error {
	tokenReader := NewTokenReader(in)

	departments, err := tokenReader.NextInt()
	if err != nil {
		if errors.Is(err, ErrUnexpectedEOF) {
			return nil
		}

		return err
	}

	for range departments {
		if err := processDepartment(tokenReader, out); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if err := RunService(os.Stdin, os.Stdout); err != nil {
		os.Exit(1)
	}
}
