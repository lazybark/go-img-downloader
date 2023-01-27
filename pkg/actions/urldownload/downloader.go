package urldownload

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/lazybark/go-helpers/cli/clf"
	"github.com/lazybark/go-img-downloader/config"
	"github.com/lazybark/go-img-downloader/pkg/getters"
	"github.com/lazybark/go-img-downloader/pkg/imgworks"
	"github.com/lazybark/go-img-downloader/pkg/parsers"
	"github.com/pterm/pterm"
)

// DownloadImagesOnPage parses gets HTML code, parses for images and downloads each image.
// It also calls to resizer and/or converter if config demands.
func DownloadImagesOnPage(path *url.URL, downloadInto string, cfg config.Config) error {
	var links []*url.URL

	spinnerLiveText, _ := pterm.DefaultSpinner.Start("Opening page...")
	spinnerLiveText.WithRemoveWhenDone()
	defer spinnerLiveText.Stop()

	var text string
	var err error
	if cfg.ForceChrome {
		text, err = getters.GetPageHTMLByChromeDP(path, cfg.Debug)
	} else {
		text, err = getters.GetPageHTML(path, cfg.Debug)
	}
	if err != nil {
		return err
	}

	//Parse code & find all images
	spinnerLiveText.UpdateText("Parsing code")
	links = parsers.ParseHTMLForImgs(text)
	l := len(links)

	//Retry using ChromeDP if there are no images (this case often happens with SPAs like React -
	//we need to use JS)
	if l == 0 {
		text, err = getters.GetPageHTMLByChromeDP(path, cfg.Debug)
		if err != nil {
			return err
		}
		links = parsers.ParseHTMLForImgs(text)
		l = len(links)
	}
	if l == 0 {
		spinnerLiveText.Fail("No images could be parsed on the page :(")
		return nil
	}
	spinnerLiveText.Success(fmt.Sprintf("Found %d images", l))

	p, _ := pterm.DefaultProgressbar.WithTotal(l).WithTitle("Downloaded images").Start()

	for n, u := range links {
		//It's better to have host checking here then passing it to img URL parser and making code
		//more complex.
		if u.Host == "" {
			u.Host = path.Host
		}

		imageBytes, err := DownloadImage(u)
		if err != nil {
			pterm.Error.Println(fmt.Sprintf("%s: %s", u.String(), err))
			continue
		}

		//We add img number, because some sites have same names for deifferent sized of the same img.
		//And we cut out any queries to get proper save name.
		if u.RawQuery != "" {
			u.RawQuery = ""
		}
		imgName := filepath.Base(u.String())
		imgPath := filepath.Join(downloadInto, fmt.Sprint(n)+"_"+imgName)

		decoded, format, err := imgworks.DecodeImage(imageBytes)
		if err != nil {
			pterm.Error.Println(fmt.Sprintf("[ERROR DECODING] %s: %s", u.String(), err))
			continue
		}

		//Ignore by size
		var ignored bool
		size := decoded.Bounds().Size()
		if cfg.MinImgHeight > 0 && cfg.MinImgWidth > 0 {
			if size.X < cfg.MinImgWidth && size.Y < cfg.MinImgHeight {
				ignored = true
			}
		} else if cfg.MinImgWidth > 0 && size.X < cfg.MinImgWidth {
			ignored = true
		} else if cfg.MinImgHeight > 0 && size.Y < cfg.MinImgHeight {
			ignored = true
		}
		if ignored {
			pterm.Error.Println(fmt.Sprintf("%s %dx%d = ignored\n", imgName, size.X, size.Y))
			continue
		}

		//Resize if needed
		if cfg.MaxImgHeight > 0 && cfg.MaxImgWidth > 0 && size.X > cfg.MaxImgWidth && size.Y > cfg.MaxImgHeight {
			if size.Y > size.X {
				decoded = imgworks.ResizeImage(decoded, cfg.MaxImgHeight, 0, imaging.Lanczos)
			} else {
				decoded = imgworks.ResizeImage(decoded, 0, cfg.MaxImgWidth, imaging.Lanczos)
			}
		} else if cfg.MaxImgWidth > 0 && size.X > cfg.MaxImgWidth {
			decoded = imgworks.ResizeImage(decoded, 0, cfg.MaxImgWidth, imaging.Lanczos)
		} else if cfg.MaxImgHeight > 0 && size.Y > cfg.MaxImgHeight {
			decoded = imgworks.ResizeImage(decoded, cfg.MaxImgHeight, 0, imaging.Lanczos)
		}

		//Save ending image
		err = imaging.Save(decoded, imgPath)
		if err != nil {
			pterm.Error.Println(fmt.Sprintf("[SAVE ERROR] %s: %s", imgName, err))
			continue
		}

		//We convert image after, because there is a chance that convertion will spoil something. So
		//we need both copies of the file.
		//We also save initial names after convertion, just change format. So user will totally
		//understand that the image was converted. Just in case.
		var converted bool
		var newName string
		var newPath string
		if cfg.ConvertAllToJPG && format != "jpeg" {
			decoded, err = imgworks.ConvertToJPG(decoded, imgworks.ImgWriter{}, format)
			converted = true
			newName = imgName + ".jpg"
			newPath = imgPath + ".jpg"
		}
		if cfg.ConvertAllToPNG && format != "png" {
			decoded, err = imgworks.ConvertToPNG(decoded, imgworks.ImgWriter{})
			converted = true
			newName = imgName + ".png"
			newPath = imgPath + ".png"
		}
		if err != nil {
			pterm.Error.Println(fmt.Sprintf("[CONVERT ERROR] %s: %s", newName, err))
			continue
		}

		if converted {
			err = imaging.Save(decoded, newPath)
			if err != nil {
				pterm.Error.Println(fmt.Sprintf("[SAVE ERROR] %s: %s", newName, err))
				continue
			}
		}

		pterm.Success.Println(u.String())
		p.Increment()
	}
	p.Stop()

	return nil
}

// DownloadImage makes GET call and returns request body
func DownloadImage(u *url.URL) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	return body, nil
}

// PerformDownloadByURL creates subdir for current download and runs DownloadImagesOnPage
func PerformDownloadByURL(lnk string, cfg config.Config) error {
	var subdir string
	downloadInto := cfg.OutputPath
	//Parse URL
	link, err := url.ParseRequestURI(lnk)
	if err != nil {
		return fmt.Errorf(clf.Red("[BAD URL]")+": %s\n", err)
	}

	//Make subdir
	if !cfg.DoNotUseSubfolders {
		path := strings.Split(link.Path, "/")
		pl := len(path)
		for i := 1; i <= pl; i++ {
			subdir = path[pl-i]
			if subdir != "" {
				break
			}
		}
		if subdir == "" {
			subdir = fmt.Sprint(time.Now().Unix())
		}
		downloadInto = filepath.Join(cfg.OutputPath, subdir)
		if err := os.MkdirAll(downloadInto, os.ModePerm); err != nil {
			return fmt.Errorf(clf.Red("[ERROR CREATING DIR]")+": %s\n", err)
		}
	}

	fmt.Printf("Saving into %s\n\n", downloadInto)

	err = DownloadImagesOnPage(link, downloadInto, cfg)
	if err != nil {
		return fmt.Errorf(clf.Red("[ERROR DOWNLOADING]")+": %s\n", err)
	}

	return nil
}
