package main

import (
	"fmt"
)

const (
	initMinT = 15
	initMaxT = 30
)

func main() {
	var (
		numDepartments int
		numEmployees   int
	)

	_, err := fmt.Scan(&numDepartments)
	if err != nil || numDepartments <= 0 {
		fmt.Println("Invalid number of departments")

		return
	}

	for range numDepartments {
		_, err = fmt.Scan(&numEmployees)
		if err != nil || numEmployees < 0 {
			fmt.Println("Invalid number of employees")

			return
		}

		if !processDepartment(numEmployees) {
			return
		}
	}
}

func processDepartment(numEmployees int) bool {
	minTemp := initMinT
	maxTemp := initMaxT

	for range numEmployees {
		var (
			sign string
			temp int
		)

		_, err := fmt.Scan(&sign, &temp)
		if err != nil {
			fmt.Println("Invalid value")

			return false
		}

		switch sign {
		case ">=":
			if temp > minTemp {
				minTemp = temp
			}
		case "<=":
			if temp < maxTemp {
				maxTemp = temp
			}
		default:
			fmt.Println("Invalid character")

			return false
		}

		if minTemp > maxTemp {
			fmt.Println(-1)
		} else {
			fmt.Println(minTemp)
		}
	}

	return true
}
