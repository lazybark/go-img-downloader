package urldownload

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/lazybark/go-helpers/cli/clf"
	"github.com/lazybark/go-img-downloader/config"
	"github.com/lazybark/go-img-downloader/pkg/clif"
	"github.com/lazybark/go-img-downloader/pkg/parsers"
)

var (
	menuActionsFile = clif.ActionList{
		Message:           "Please, provide full path to a text/csv file. File must contain links, divided by comma or like breaks",
		WrongInputMessage: "unknown command",
		CanGoBack:         true,
		BackKey:           "B",
		CanExit:           true,
		ExitKey:           "E",
	}
)

// GetAllImgsFile runs sequence to get all links from the file and to download all images on each page.
func GetAllImgsFile(cfg config.Config) error {
	var linksList []string
	for {
		menuActionsFile.Promt()
		userCommand, pathToFile := menuActionsFile.AwaitInput(parsers.ParseURL)
		if userCommand.Key == menuActionsFile.ExitKey {
			fmt.Println("Exiting")
			os.Exit(0)
		}
		if userCommand.Key == menuActionsFile.BackKey {
			return nil
		}
		if pathToFile == "" {
			fmt.Println(clf.Red("no path found"))
			continue
		}
		f, err := os.Open(pathToFile)
		if err != nil {
			fmt.Println(clf.Red("can not open file:"), err)
			continue
		}
		fileScanner := bufio.NewScanner(f)
		fileScanner.Split(bufio.ScanLines)
		var lineSkuArr = []string{}
		var txt string
		for fileScanner.Scan() {
			txt = fileScanner.Text()
			if txt == "" {
				continue
			}
			lineSkuArr = strings.Split(txt, ",")
			for _, s := range lineSkuArr {
				if s != "" {
					linksList = append(linksList, s)
				}
			}
		}
		f.Close()

		for _, l := range linksList {
			if err := PerformDownloadByURL(l, cfg); err != nil {
				fmt.Println(clf.Red(err))
			}
			fmt.Println()
		}
		fmt.Println()

	}
}
