package urldownload

import (
	"fmt"
	"net/url"
	"os"

	"github.com/lazybark/go-helpers/cli/clf"
	"github.com/lazybark/go-img-downloader/config"
	"github.com/lazybark/go-img-downloader/pkg/clif"
	"github.com/lazybark/go-img-downloader/pkg/getters"
	"github.com/lazybark/go-img-downloader/pkg/parsers"
)

var (
	menuActionsRecusive = clif.ActionList{
		Message:           "Please, provide a link for recursive link scan",
		WrongInputMessage: "unknown command",
		CanGoBack:         true,
		BackKey:           "B",
		CanExit:           true,
		ExitKey:           "E",
	}

	menuIgnoreExternal = clif.ActionList{
		Message:           "Ignore external sites?",
		WrongInputMessage: "unknown command",
		Actions:           []clif.Action{yes, no},
		CanGoBack:         false,
		CanExit:           false,
	}
	yes = clif.Action{Key: "Y", Text: "yes"}
	no  = clif.Action{Key: "N", Text: "no"}
)

// GetAllImgsRecursively runs sequence to get all links from the page and to download all images on each page.
func GetAllImgsRecursively(cfg config.Config) error {
	var downloadedList = make(map[string]bool)
	var ignoreExternal bool
	for {
		ignoreExternal = false
		downloadedList = map[string]bool{}

		menuActionsRecusive.Promt()
		userCommand, userLink := menuActionsRecusive.AwaitInput(parsers.ParseURL)
		if userCommand.Key == menuActionsRecusive.ExitKey {
			fmt.Println("Exiting")
			os.Exit(0)
		}
		if userCommand.Key == menuActionsRecusive.BackKey {
			return nil
		}
		if userLink == "" {
			fmt.Println(clf.Red("no link found"))
			continue
		}
		u, err := url.ParseRequestURI(userLink)
		if err != nil {
			fmt.Println(clf.Red("[BAD LINK]"), err)
			continue
		}

		menuIgnoreExternal.Promt()
		userCommand = menuIgnoreExternal.AwaitCommand()
		if userCommand.Key == yes.Key {
			ignoreExternal = true
		}

		pageText, err := getters.GetPageHTML(u, cfg.Debug)
		if err != nil {
			fmt.Println(clf.Red(err))
			continue
		}

		links := parsers.ParseHTMLForLinks(pageText)
		l := len(links)

		if l == 0 {
			pageText, err = getters.GetPageHTMLByChromeDP(u, cfg.Debug)
			if err != nil {
				fmt.Println(clf.Red(err))
			}
			links = parsers.ParseHTMLForLinks(pageText)
			l = len(links)
		}
		if l == 0 {
			fmt.Println("No images could be parsed on the page :(")
			continue
		}

		fmt.Printf("Found %d links\n", l)

		for _, l := range links {
			if l.Host == "" {
				l.Host = u.Host
			}

			if ignoreExternal && l.Host != u.Host {
				fmt.Println(clf.Red(fmt.Sprintf("%s is external = ignored", u)))
				continue
			}
			if _, ok := downloadedList[l.String()]; ok {
				fmt.Println(clf.Red(fmt.Sprintf("%s already downloaded = ignored", u)))
				continue
			}

			if err := PerformDownloadByURL(l.String(), cfg); err != nil {
				fmt.Println(clf.Red(err))
			} else {
				downloadedList[l.String()] = true
			}
			fmt.Println()
		}
		fmt.Println()
	}
}
