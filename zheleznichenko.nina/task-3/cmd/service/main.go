package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  int    `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	ValueStr string `xml:"Value"`
	ParsedVal float64
}

// Структура для записи в JSON
type CurrencyResult struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	configFile, err := os.ReadFile(*configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to read config file: %v", err))
	}

	var cfg Config
	err = yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to parse config yaml: %v", err))
	}

	xmlData, err := os.ReadFile(strings.TrimSpace(cfg.InputFile))
	if err != nil {
		panic(fmt.Sprintf("failed to read input xml file: %v", err))
	}

	var valCurs ValCurs
	err = xml.Unmarshal(xmlData, &valCurs)
	if err != nil {
		panic(fmt.Sprintf("failed to parse xml data: %v", err))
	}

	for i := range valCurs.Valutes {
		valClean := strings.ReplaceAll(valCurs.Valutes[i].ValueStr, ",", ".")
		valFloat, err := strconv.ParseFloat(strings.TrimSpace(valClean), 64)
		if err != nil {
			panic(fmt.Sprintf("failed to convert currency value: %v", err))
		}
		valCurs.Valutes[i].ParsedVal = valFloat
	}

	sort.Slice(valCurs.Valutes, func(i, j int) bool {
		return valCurs.Valutes[i].ParsedVal > valCurs.Valutes[j].ParsedVal
	})

	var results []CurrencyResult
	for _, v := range valCurs.Valutes {
		results = append(results, CurrencyResult{
			NumCode:  v.NumCode,
			CharCode: v.CharCode,
			Value:    v.ParsedVal,
		})
	}

	outputDir := filepath.Dir(cfg.OutputFile)
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		panic(fmt.Sprintf("failed to create output directory: %v", err))
	}

	jsonOutput, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(fmt.Sprintf("failed to marshal json: %v", err))
	}

	err = os.WriteFile(cfg.OutputFile, jsonOutput, 0644)
	if err != nil {
		panic(fmt.Sprintf("failed to write output file: %v", err))
	}
}