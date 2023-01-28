package urldownload

import (
	"fmt"
	"os"

	"github.com/lazybark/go-helpers/cli/clf"
	"github.com/lazybark/go-img-downloader/config"
	"github.com/lazybark/go-img-downloader/pkg/clif"
	"github.com/lazybark/go-img-downloader/pkg/parser"
)

var (
	menuActions = clif.ActionList{
		Message:           "Please, provide a link",
		WrongInputMessage: "unknown command",
		CanGoBack:         true,
		BackKey:           "B",
		CanExit:           true,
		ExitKey:           "E",
	}
)

// AllImgsByURL runs sequence to download all images on the page
func AllImgsByURL(cfg config.Config) error {
	for {
		menuActions.Promt()
		userCommand, userLink := menuActions.AwaitInput(parser.ParseURL)
		if userCommand.Key == menuActions.ExitKey {
			fmt.Println("Exiting")
			os.Exit(0)
		}
		if userCommand.Key == menuActions.BackKey {
			return nil
		}
		if userLink == "" {
			fmt.Println(clf.Red("no link found"))
			continue
		}

		if err := PerformDownloadByURL(userLink, cfg); err != nil {
			fmt.Println(clf.Red(err))
			continue
		}
		fmt.Println()
	}
}
