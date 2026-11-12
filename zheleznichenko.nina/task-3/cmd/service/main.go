package main

import (
	"flag"
	"log"
	"strings"

	"zheleznichenko.nina/task-3/internal/config"
	"zheleznichenko.nina/task-3/internal/currency"
	"zheleznichenko.nina/task-3/internal/storage"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.ReadConfig(*configPath)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	currencies, err := currency.ReadCurrencies(strings.TrimSpace(cfg.InputFile))
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if err := currency.ParseValues(currencies); err != nil {
		log.Fatalf("error: %v", err)
	}

	results := currency.Convert(currencies)

	if err := storage.WriteResults(cfg.OutputFile, results); err != nil {
		log.Fatalf("error: %v", err)
	}
}