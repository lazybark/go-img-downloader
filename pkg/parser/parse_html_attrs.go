package parser

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// ParseURL returns error if URL in 's' is incorrect
func ParseURL(s string) error {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return err
	}
	return nil
}

// ParseHTMLAttrsForImgs parses all attributes to find links to images (jpg, png, webp)
func ParseHTMLAttrsForImgs(attrs []html.Attribute) []*url.URL {
	var strs []*url.URL
	var splitted []string
	var u *url.URL
	var err error

	//We must parse all attrs as some CMS put links to images into unexpected fields
	for _, at := range attrs {
		if at.Val != "" {
			//At first, we simply check that the link has any image formats. It's faster than using
			//Regexp
			if strings.Contains(at.Val, ".jpg") || strings.Contains(at.Val, ".jpeg") || strings.Contains(at.Val, ".png") || strings.Contains(at.Val, ".webp") || strings.Contains(at.Val, ".JPG") || strings.Contains(at.Val, ".JPEG") || strings.Contains(at.Val, ".PNG") {
				//We try to split in case there are several img links in one attr (for some CMS)
				splitted = strings.Split(at.Val, ",")
				/*if len(splitted) < 2 {
					//Try different keys
				}*/
				for _, spl := range splitted {
					//fmt.Println(at.Key, ":", spl)
					//at.Val = strings.TrimPrefix(at.Val, "//")
					if strings.HasPrefix(spl, "//") {
						spl = "https:" + spl
					}

					//Then we convert strings into native Go urls - to check for some basic problems
					u, err = url.ParseRequestURI(spl)
					if err == nil && u.Path != "" {
						//fmt.Println(u)
						//fmt.Println()
						if u.Scheme == "" {
							u.Scheme = "https"
						}
						strs = append(strs, u)
					}
				}

			}
		}
	}
	return strs
}

// GetHTMLAttr returns value of specified HTML attribute
func GetHTMLAttr(attr string, attrs []html.Attribute) string {
	for _, at := range attrs {
		if at.Key == attr {
			return at.Val
		}
	}
	return ""
}
