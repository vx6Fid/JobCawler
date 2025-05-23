package main

import (
	"log"

	"github.com/gocolly/colly/v2"
	"github.com/vx6fid/job-crawler/internal/downloader"
	"github.com/vx6fid/job-crawler/internal/parser/sites"
	"github.com/vx6fid/job-crawler/internal/urlfrontier"
)

func main() {
	frontier := urlfrontier.NewFrontier(100)
	d := downloader.NewDownloader()

	// Add the search results listing page as a task
	frontier.Add(urlfrontier.CrawlTask{
		URL:  "https://weworkremotely.com/remote-jobs/search?term=DevOps",
		Type: "listing",
	})

	for frontier.QueueSize() > 0 {
		task := frontier.GetNext()

		log.Printf("Crawling: %s [%s]", task.URL, task.Type)

		parser := sites.GetParser(task.URL)
		if parser == nil {
			log.Printf("--- X --- No parser for URL: %s", task.URL)
			continue
		}

		switch task.Type {
		case "listing":
			err := d.FetchWithParser(task.URL, func(e *colly.HTMLElement) {
				jobs, err := parser.Parse(e)
				if err != nil {
					log.Printf("--- X --- Parser error: %v", err)
					return
				}
				for _, job := range jobs {
					frontier.Add(urlfrontier.CrawlTask{
						URL:  job.ApplyURL,
						Type: "job",
						Meta: map[string]string{
							"title":   job.Title,
							"company": job.Company,
						},
					})
				}
			})
			if err != nil {
				log.Printf("--- X --- Fetch error: %v", err)
			}

		case "job":
			err := d.FetchWithParser(task.URL, func(e *colly.HTMLElement) {
				job, err := parser.ParseJobDescription(e)
				if err != nil {
					log.Printf("--- X --- Job parser error: %v", err)
					return
				}
				log.Printf("--- :) Full Job Parsed: %+v", job)
			})
			if err != nil {
				log.Printf("--- X --- Fetch error: %v", err)
			}
		}
	}
}
