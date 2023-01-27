// v1m is a special version for my own use. It has some pre-defined settings and nothing more.
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

	cfg.ForceChrome = true

	menu.PromtMainMenu(cfg)
}
