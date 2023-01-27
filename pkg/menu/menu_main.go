package menu

import (
	"fmt"
	"log"
	"os"

	"github.com/lazybark/go-img-downloader/config"
	"github.com/lazybark/go-img-downloader/pkg/actions/urldownload"
	"github.com/lazybark/go-img-downloader/pkg/clif"
)

// Vars & function below are moved here to avoid doubles im main() functions of different versions
// with same menu.
var (
	GreeTingText = "Welcome to image downloader (v.%s)\n\n"

	MainMenuActions = clif.ActionList{
		Message: "Please, choose a mode",
		Actions: []clif.Action{
			DownloadURLAction,
			DownloadFromFileAction,
			DownloadRecursivelyAction,
		},
		WrongInputMessage: "unknown command",
		BackKey:           "B",
		CanExit:           true,
		ExitKey:           "E",
	}

	DownloadURLAction         = clif.Action{Key: "D", Text: " to download images by single link"}
	DownloadFromFileAction    = clif.Action{Key: "F", Text: " to download images by list of links in a file"}
	DownloadRecursivelyAction = clif.Action{Key: "R", Text: " to recursively download images by a link"}

	MenuActionsFile = clif.ActionList{
		Message:           "Please, provide full path to a text/csv file. File must contain links, divided by comma or like breaks",
		WrongInputMessage: "unknown command",
		CanGoBack:         true,
		BackKey:           "B",
		CanExit:           true,
		ExitKey:           "E",
	}
)

// PromtMainMenu should be called in main() of v1 packages
func PromtMainMenu(cfg config.Config) {
	fmt.Printf(GreeTingText, config.Ver)
	fmt.Printf("Images will be saved in: %s\n", cfg.OutputPath)

	var err error
	for {
		MainMenuActions.Promt()
		userCommand := MainMenuActions.AwaitCommand()
		if userCommand.Key == MainMenuActions.ExitKey {
			fmt.Println("Exiting")
			os.Exit(0)
		}

		if userCommand == DownloadURLAction {
			err = urldownload.AllImgsByURL(cfg)
			if err != nil {
				log.Fatal("[ERROR]", err)
				return
			}
		} else if userCommand == DownloadFromFileAction {
			err = urldownload.GetAllImgsFile(cfg)
			if err != nil {
				log.Fatal("[ERROR]", err)
				return
			}
		} else if userCommand == DownloadRecursivelyAction {
			err = urldownload.GetAllImgsRecursively(cfg)
			if err != nil {
				log.Fatal("[ERROR]", err)
				return
			}
		}
	}
}
