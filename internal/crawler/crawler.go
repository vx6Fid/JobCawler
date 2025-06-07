package crawler

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/vx6fid/job-crawler/internal/crawler/sites"
	"github.com/vx6fid/job-crawler/internal/downloader"
	"github.com/vx6fid/job-crawler/internal/urlfrontier"
	"github.com/vx6fid/job-crawler/pkg"
)

func StartCrawling(roles []string, maxJobs int, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := pkg.ConnectMongo(); err != nil {
		return fmt.Errorf("MongoDB connection failed: %v", err)
	}

	start := time.Now()
	jobCounter := 0

	frontier := urlfrontier.NewFrontier(100)
	d := downloader.NewDownloader()

	for _, role := range roles {
		if !IsRoleAllowed(role) {
			log.Printf("--- [ERROR] --- Role not allowed: %s", role)
			continue // skip invalid
		}
		searchURL := fmt.Sprintf("https://weworkremotely.com/remote-jobs/search?term=%s", url.QueryEscape(role))
		frontier.Add(urlfrontier.CrawlTask{
			URL:  searchURL,
			Type: "listing",
		})
	}

	tick := time.NewTicker(5 * time.Second)
	defer tick.Stop()

	for frontier.QueueSize() > 0 {
		select {
		case <-ctx.Done():
			log.Println("--- [STOP] --- Context timeout reached.")
			log.Printf("Crawler finished early. Total jobs saved: %d | Duration: %.2fs", jobCounter, time.Since(start).Seconds())
			return nil
		case <-tick.C:
			log.Printf("[DEBUG] Queue Size: %d, Job Count: %d", frontier.QueueSize(), jobCounter)
		default:
			// continue loop
		}

		if jobCounter >= maxJobs {
			log.Printf("--- :| --- Reached job limit (%d). Stopping crawl.", maxJobs)
			break
		}

		task := frontier.GetNext()
		log.Printf("Crawling: %s [%s]", task.URL, task.Type)

		parser := sites.GetParser(task.URL)
		if parser == nil {
			log.Printf("--- [ERROR] --- No parser for URL: %s", task.URL)
			continue
		}

		log.Printf("[STEP] Starting fetch: %s", task.URL)
		switch task.Type {
		case "listing":
			err := d.FetchWithParser(ctx, task.URL, func(e *colly.HTMLElement) {
				log.Printf("[STEP] Entered listing handler for: %s", task.URL)
				jobs, err := parser.Parse(ctx, e)
				if err != nil {
					log.Printf("--- [ERROR] --- Parser error: %v", err)
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
				log.Printf("--- [ERROR] --- Fetch error: %v", err)
			} else {
				log.Printf("[STEP] List fetch completed: %s", task.URL)
			}

		case "job":
			err := d.FetchWithParser(ctx, task.URL, func(e *colly.HTMLElement) {
				log.Printf("[STEP] Entered job handler for: %s", task.URL)
				job, err := parser.ParseJobDescription(e)
				if err != nil {
					log.Printf("--- [ERROR] --- Job parser error: %v", err)
					return
				}
				if err := pkg.UpsertJob(job); err != nil {
					log.Printf("--- [ERROR] --- Failed to save job: %v", err)
				} else {
					log.Printf("--- :) --- Saved job: %s @ %s", job.Title, job.Company)
					jobCounter++
				}
			})
			if err != nil {
				log.Printf("--- [ERROR] --- Fetch error: %v", err)
			} else {
				log.Printf("[STEP] Job fetch completed: %s", task.URL)
			}
		}
	}

	log.Printf("Crawler finished. Total jobs saved: %d | Duration: %.2fs", jobCounter, time.Since(start).Seconds())
	return nil
}
