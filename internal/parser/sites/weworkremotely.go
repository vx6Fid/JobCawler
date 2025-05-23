package sites

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/gocolly/colly/v2"
)

type WeWorkRemotelyParser struct{}

func init() {
	Register(&WeWorkRemotelyParser{})
}

func (p *WeWorkRemotelyParser) Matches(url string) bool {
	return strings.Contains(url, "weworkremotely.com")
}

func (p *WeWorkRemotelyParser) Parse(e *colly.HTMLElement) ([]JobPosting, error) {
	var jobs []JobPosting

	e.ForEach("li.new-listing-container.feature > a[href^='/listings/']", func(_ int, el *colly.HTMLElement) {
		job := JobPosting{
			Title:    strings.TrimSpace(el.ChildText("h4.new-listing__header__title")),
			Company:  strings.TrimSpace(el.ChildText("p.new-listing__company-name")),
			Location: strings.TrimSpace(el.ChildText("p.new-listing__company-headquarters")),
			ApplyURL: "https://weworkremotely.com" + el.Attr("href"),
		}

		fmt.Println("[weworkremotely] -> Parsed:", job.Title, "-", job.Company)
		jobs = append(jobs, job)
	})

	return jobs, nil
}

func (p *WeWorkRemotelyParser) ParseJobDescription(e *colly.HTMLElement) (JobPosting, error) {
	job := JobPosting{
		Title:       strings.TrimSpace(e.ChildText("h2.lis-container__header__hero__company-info__title")),
		Company:     strings.TrimSpace(e.ChildText("div.lis-container__header__hero__company-info__description strong")),
		URL:         e.Request.URL.String(),
		Source:      "weworkremotely.com",
		Description: strings.TrimSpace(e.ChildText("div.lis-container__job__content__description")),
		ApplyURL:    e.ChildAttr("a#job-cta-alt", "href"),
	}

	// Required fields check
	if job.Title == "" || job.Company == "" {
		return JobPosting{}, fmt.Errorf("failed to parse job description: title or company missing")
	}

	// PostedOn
	job.PostedOn = strings.TrimSpace(e.ChildText("li.lis-container__job__sidebar__job-about__list__item:contains('Posted on') span"))

	// Salary
	job.Salary = strings.TrimSpace(e.ChildText("li.lis-container__job__sidebar__job-about__list__item:contains('Salary') span"))

	// Region/Location (handle absence gracefully)
	location := ""
	e.ForEach("li.lis-container__job__sidebar__job-about__list__item--full:contains('Region') span.box", func(_ int, el *colly.HTMLElement) {
		region := strings.TrimSpace(el.Text)
		if region != "" {
			location += region + ", "
		}
	})
	job.Location = strings.TrimSuffix(location, ", ")

	// Skills from sidebar
	skillsSet := make(map[string]struct{})
	e.ForEach("li.lis-container__job__sidebar__job-about__list__item--full:contains('Skills') span.box", func(_ int, el *colly.HTMLElement) {
		skill := strings.TrimSpace(el.Text)
		if skill != "" {
			skillsSet[skill] = struct{}{}
		}
	})

	// Skills from description (keyword match)
	lowerDesc := strings.ToLower(job.Description)
	for _, keyword := range knownTechnologies {
		if strings.Contains(lowerDesc, strings.ToLower(keyword)) {
			skillsSet[keyword] = struct{}{}
		}
	}

	// Deduplicate and sort skills
	for skill := range skillsSet {
		job.Skills = append(job.Skills, skill)
	}
	sort.Strings(job.Skills)

	// Extract experience from description
	job.Experience = extractExperience(job.Description)

	return job, nil
}

func extractExperience(desc string) string {
	re := regexp.MustCompile(`\b(\d{1,2})\+?\s*(years|yrs?)\b`)
	match := re.FindStringSubmatch(strings.ToLower(desc))
	if len(match) > 1 {
		return fmt.Sprintf("%s years", match[1])
	}
	return ""
}
