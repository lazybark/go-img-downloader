package main

import (
	"log"

	"github.com/lazybark/go-img-downloader/config"
	"github.com/lazybark/go-img-downloader/pkg/menu"
)

func main() {
	cfg, err := config.InitApp()
	if err != nil {
		log.Fatal("[APP INIT ERROR]", err)
		return
	}

	menu.PromtMainMenu(cfg)
}
