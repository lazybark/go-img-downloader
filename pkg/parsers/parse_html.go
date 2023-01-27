package parsers

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// ParseHTMLForImgs searches for tokens (tags) that may contain links to images
func ParseHTMLForImgs(text string) []*url.URL {
	var links []*url.URL
	tkn := html.NewTokenizer(strings.NewReader(text))
	var a []*url.URL
	for {
		a = []*url.URL{}
		tt := tkn.Next()
		if tt == html.ErrorToken {
			break
		}

		t := tkn.Token()
		if t.Data == "img" || t.Data == "li" || t.Data == "a" || t.Data == "div" {
			a = ParseHTMLAttrsForImgs(t.Attr)
		}

		if len(a) > 0 {
			links = append(links, a...)

		}
	}
	return links
}

// ParseHTMLForLinks searces for links in HTML code
func ParseHTMLForLinks(text string) []*url.URL {
	var links []*url.URL
	tkn := html.NewTokenizer(strings.NewReader(text))
	var a *url.URL
	var lnk string
	var err error
	for {
		tt := tkn.Next()
		if tt == html.ErrorToken {
			break
		}

		t := tkn.Token()
		if t.Data == "a" {

			lnk = GetHTMLAttr("href", t.Attr)
			if lnk != "" {
				a, err = url.ParseRequestURI(lnk)
				if err != nil {
					continue
				}
				if a.Scheme == "" {
					a.Scheme = "https"
				}
				links = append(links, a)
			}
		}

	}
	return links
}
