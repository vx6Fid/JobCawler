package downloader

import (
	"github.com/gocolly/colly/v2"
)

type Downloader struct {
	collector *colly.Collector
}

func NewDownloader() *Downloader {
	c := colly.NewCollector(
		colly.Async(true),
		// colly.AllowedDomains("weworkremotely.com", "amazon.jobs", "linkedin.com"),
	)
	return &Downloader{collector: c}
}

func (d *Downloader) FetchWithParser(url string, parseFunc func(e *colly.HTMLElement)) error {
	d.collector.OnHTML("body", parseFunc)

	err := d.collector.Visit(url)
	if err != nil {
		return err
	}

	d.collector.Wait() // wait for async
	return nil
}
