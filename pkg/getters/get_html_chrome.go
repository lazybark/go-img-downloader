package getters

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/lazybark/go-helpers/fsw"
)

// GetPageHTMLByChromeDP returns HTML code of desired URL. It uses ChromeDP package to get page data.
// Basically, it runs Chrome in Incognito
func GetPageHTMLByChromeDP(url *url.URL, writeDebug bool) (string, error) {
	var html string
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", true),
	)...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(
		ctx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	if err := chromedp.Run(ctx,
		chromedp.Navigate(url.String()),
		chromedp.OuterHTML("html", &html),
	); err != nil {
		return "", err
	}

	if writeDebug {
		f, err := fsw.MakePathToFile(filepath.Join("debug", fmt.Sprint(time.Now().Unix())), true)
		if err != nil {
			return "", fmt.Errorf("[GetPageHTMLByChromeDP] %w", err)
		}
		defer f.Close()
		f.WriteString(html)
	}

	return html, nil
}
