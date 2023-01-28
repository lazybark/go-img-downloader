package getter

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/lazybark/go-helpers/fsw"
)

// GetPageHTML returns HTML code of desired URL
func GetPageHTML(url *url.URL, writeDebug bool) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	resp.Body.Close()

	html := string(body)

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
