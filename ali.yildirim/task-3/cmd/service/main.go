package main

import (
	"flag"

	"github.com/aliyildirim/task-3/internal/config"
	"github.com/aliyildirim/task-3/internal/currency"
	"github.com/aliyildirim/task-3/internal/storage"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		panic(err)
	}

	items, err := currency.Load(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	currency.SortByValueDesc(items)

	if err := storage.WriteJSON(cfg.OutputFile, items); err != nil {
		panic(err)
	}
}
