package currency

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/transform"
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode   int     `xml:"NumCode"`
	CharCode  string  `xml:"CharCode"`
	ValueStr  string  `xml:"Value"`
	ParsedVal float64 `xml:"-"`
}

type CurrencyResult struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func ReadCurrencies(path string) ([]Valute, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read input xml file: %w", err)
	}

	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		encoding, err := htmlindex.Get(charset)
		if err != nil {
			return nil, fmt.Errorf("failed to detect XML encoding: %w", err)
		}

		return transform.NewReader(input, encoding.NewDecoder()), nil
	}

	var valCurs ValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		return nil, fmt.Errorf("failed to parse xml data: %w", err)
	}

	return valCurs.Valutes, nil
}

func ParseValues(currencies []Valute) error {
	for index := range currencies {
		value := strings.ReplaceAll(currencies[index].ValueStr, ",", ".")

		parsedValue, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
		if err != nil {
			return fmt.Errorf("failed to convert currency value: %w", err)
		}

		currencies[index].ParsedVal = parsedValue
	}

	return nil
}

func Convert(currencies []Valute) []CurrencyResult {
	sort.Slice(currencies, func(first, second int) bool {
		return currencies[first].ParsedVal > currencies[second].ParsedVal
	})

	results := make([]CurrencyResult, 0, len(currencies))
	for _, currency := range currencies {
		results = append(results, CurrencyResult{
			NumCode:  currency.NumCode,
			CharCode: currency.CharCode,
			Value:    currency.ParsedVal,
		})
	}

	return results
}