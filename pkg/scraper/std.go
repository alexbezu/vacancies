package scraper

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

type Std struct {
	log *logrus.Logger
}

func NewSTD(log *logrus.Logger) Std {
	return Std{log: log}
}

func (s Std) UrlsFromSite(ctx context.Context, link, filter string) ([]string, error) {
	var ret []string

	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel() // Ensure the cancel function is called to release resources

	// Create a new HTTP GET request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return nil, err
	}

	// Create an HTTP client
	client := &http.Client{}

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error performing request: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					// Check if we have relative link. If so, add a host
					u, err := url.Parse(attr.Val)
					if err != nil {
						s.log.Error(err)
					}
					if u.Host == "" {
						l, err := url.Parse(link)
						if err != nil {
							s.log.Error(err)
						}
						u.Host = l.Host
						u.Scheme = l.Scheme
					}

					if filter != "" {
						r, err := regexp.Compile(filter)
						if err != nil {
							s.log.Error(err)
						} else {
							str := u.String()
							if r.MatchString(str) {
								ret = append(ret, str)
							}
						}
					} else {
						ret = append(ret, u.String())
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return ret, nil
}
