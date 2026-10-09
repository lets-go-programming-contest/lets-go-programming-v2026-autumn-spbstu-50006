package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
	"gopkg.in/yaml.v3"
)

var (
	errConfigRequired = errors.New("config file flag -config is required")
	errInputRequired  = errors.New("input-file is required in config")
	errOutputRequired = errors.New("output-file is required in config")
)

type Config struct {
	InputFile  string `yaml:"inputFile"`
	OutputFile string `yaml:"outputFile"`
}

type ValCurs struct {
	Valute []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  int    `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	ValueStr string `xml:"Value"`
}

type CurrencyResult struct {
	NumCode  int     `json:"numCode"`
	CharCode string  `json:"charCode"`
	Value    float64 `json:"value"`
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	if strings.TrimSpace(*configPath) == "" {
		return errConfigRequired
	}

	cfg, err := loadConfig(*configPath)
	if err != nil {
		return err
	}

	results, err := processCurrencyData(cfg.InputFile)
	if err != nil {
		return err
	}

	return saveResults(cfg.OutputFile, results)
}

func loadConfig(path string) (Config, error) {
	var cfg Config

	configData, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("read config file: %w", err)
	}

	var values map[string]string
	if err := yaml.Unmarshal(configData, &values); err != nil {
		return cfg, fmt.Errorf("unmarshal yaml: %w", err)
	}

	cfg.InputFile = strings.TrimSpace(values["input-file"])
	if cfg.InputFile == "" {
		cfg.InputFile = strings.TrimSpace(values["inputFile"])
	}

	cfg.OutputFile = strings.TrimSpace(values["output-file"])
	if cfg.OutputFile == "" {
		cfg.OutputFile = strings.TrimSpace(values["outputFile"])
	}

	if cfg.InputFile == "" {
		return cfg, errInputRequired
	}
	if cfg.OutputFile == "" {
		return cfg, errOutputRequired
	}

	return cfg, nil
}

func processCurrencyData(inputFile string) ([]CurrencyResult, error) {
	xmlData, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, fmt.Errorf("read xml file: %w", err)
	}

	decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs ValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	results := make([]CurrencyResult, 0, len(valCurs.Valute))

	for _, item := range valCurs.Valute {
		cleanValue := strings.ReplaceAll(
			strings.TrimSpace(item.ValueStr),
			",",
			".",
		)

		value, err := strconv.ParseFloat(cleanValue, 64)
		if err != nil {
			return nil, fmt.Errorf(
				"parse currency value %q: %w",
				item.ValueStr,
				err,
			)
		}

		results = append(results, CurrencyResult{
			NumCode:  item.NumCode,
			CharCode: item.CharCode,
			Value:    value,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Value > results[j].Value
	})

	return results, nil
}

func saveResults(outputFile string, results []CurrencyResult) error {
	outDir := filepath.Dir(outputFile)

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	output := make([]map[string]any, 0, len(results))
	for _, item := range results {
		output = append(output, map[string]any{
			"num_code":  item.NumCode,
			"char_code": item.CharCode,
			"value":     item.Value,
		})
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	jsonData = append(jsonData, '\n')

	if err := os.WriteFile(outputFile, jsonData, 0o600); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	return nil
}
