package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	initMinT       = 15
	initMaxT       = 30
	signG          = ">="
	signL          = "<="
	impossibleTemp = -1
)

var ErrInvalidSign = errors.New("invalid sign")

type TempRange struct {
	minTemp int
	maxTemp int
}

func NewTempRange() *TempRange {
	return &TempRange{
		minTemp: initMinT,
		maxTemp: initMaxT,
	}
}

func (tempRange *TempRange) Update(sign string, temp int) error {
	switch sign {
	case signG:
		if temp > tempRange.minTemp {
			tempRange.minTemp = temp
		}
	case signL:
		if temp < tempRange.maxTemp {
			tempRange.maxTemp = temp
		}
	default:
		return fmt.Errorf("%w: %s", ErrInvalidSign, sign)
	}

	return nil
}

func (tempRange *TempRange) GetOptimal() int {
	if tempRange.minTemp > tempRange.maxTemp {
		return impossibleTemp
	}

	return tempRange.minTemp
}

func processDepartment(numEmployees int, out io.Writer) error {
	tempRange := NewTempRange()

	for range numEmployees {
		var (
			sign string
			temp int
		)

		_, err := fmt.Scan(&sign, &temp)
		if err != nil {
			return fmt.Errorf("invalid value: %w", err)
		}

		if err := tempRange.Update(sign, temp); err != nil {
			return err
		}

		_, err = fmt.Fprintln(out, tempRange.GetOptimal())
		if err != nil {
			return fmt.Errorf("failed to write result: %w", err)
		}
	}

	return nil
}

func main() {
	var numDepartments int

	_, err := fmt.Scan(&numDepartments)
	if err != nil || numDepartments <= 0 {
		fmt.Println("Invalid number of departments")

		return
	}

	for range numDepartments {
		var numEmployees int

		_, err = fmt.Scan(&numEmployees)
		if err != nil || numEmployees < 0 {
			fmt.Println("Invalid number of employees")

			return
		}

		if err := processDepartment(numEmployees, os.Stdout); err != nil {
			fmt.Println(err)

			return
		}
	}
}
