package currency

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

type Currency struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

type valCurs struct {
	Valutes []valute `xml:"Valute"`
}

type valute struct {
	NumCode  int    `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

func Load(path string) ([]Currency, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read xml file: %w", err)
	}

	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = charset.NewReaderLabel

	var doc valCurs
	if err := decoder.Decode(&doc); err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	items := make([]Currency, 0, len(doc.Valutes))

	for _, v := range doc.Valutes {
		raw := strings.ReplaceAll(strings.TrimSpace(v.Value), ",", ".")

		value, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return nil, fmt.Errorf("parse value %q: %w", v.Value, err)
		}

		items = append(items, Currency{
			NumCode:  v.NumCode,
			CharCode: v.CharCode,
			Value:    value,
		})
	}

	return items, nil
}

func SortByValueDesc(items []Currency) {
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Value > items[j].Value
	})
}
