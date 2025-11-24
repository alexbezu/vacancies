package scraper

import (
	"context"
	"net/url"
	"regexp"

	"github.com/gocolly/colly/v2"
	"github.com/sirupsen/logrus"
)

type Colly struct {
	log *logrus.Logger
}

func NewColly(log *logrus.Logger) Colly {
	return Colly{log: log}
}

func (c Colly) UrlsFromSite(ctx context.Context, jobsite, filter string) ([]string, error) {
	var ret []string

	js, err := url.Parse(jobsite)
	if err != nil {
		c.log.Error(err)
	}

	collector := colly.NewCollector(
		colly.AllowedDomains(js.Host),
	)
	collector.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

	// On every a element which has href attribute call callback
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		u, err := url.Parse(link)
		if err != nil {
			c.log.Error(err)
		}
		if u.Host == "" {
			u.Host = js.Host
			u.Scheme = js.Scheme
		}

		if filter != "" {
			r, err := regexp.Compile(filter)
			if err != nil {
				c.log.Error(err)
			} else {
				str := u.String()
				if r.MatchString(str) {
					ret = append(ret, str)
				}
			}
		} else {
			ret = append(ret, u.String())
		}
	})

	err = collector.Visit(jobsite)
	if err != nil {
		c.log.Error(err)
		return nil, err
	}

	return ret, nil
}
