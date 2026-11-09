package main

import (
	"log"
	"os"

	"github.com/sohan-14/go-starter/internal/platform/config"
)

func main() {
	cfg, err := config.Load(os.Getenv)

	if err != nil {
		log.Fatalf("%v", err)
	}

	println(cfg.App.Envrionment)
}
