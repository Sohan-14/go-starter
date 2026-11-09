package main

import (
	"log"
	"os"

	"github.com/sohan-14/go-starter/internal/platform/config"
	"github.com/sohan-14/go-starter/internal/platform/logger"
)

func main() {
	cfg, err := config.Load(os.Getenv)

	if err != nil {
		log.Fatalf("%v", err)
	}

	logger.Setup(cfg.Logger)

	println(cfg.App.Envrionment)
}
