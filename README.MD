# Remote Job Scraper

A powerful and extensible web scraper built in Go using the Colly framework to collect remote job listings from [WeWorkRemotely](https://weworkremotely.com/).

## 🛠 Tech Stack
- **Go (Golang)**
- **Colly** – Fast and elegant scraping framework for Go
- **HTML Parsing** – Custom selectors to handle dynamic job page structures

## 📦 Features
- Extracts structured job data: `Title`, `Company`, `Location`, `Salary`, `Skills`, `Description`, `Posted Date`, `Apply URL`
- Skill detection via keyword matching (future-ready for NLP-based tech stack detection)
- Graceful error handling for missing or malformed data
- Modular design for easy expansion to other platforms or data pipelines

## 🧱 Data Model
```go
type JobPosting struct {
  Title       string
  Company     string
  Location    string
  Salary      string
  PostedOn    string
  Description string
  URL         string
  Source      string
  ApplyURL    string
  Skills      []string
  Experience  string
}

# Future Improvements
- Integrate NLP for more accurate tech stack detection
- Add support for multiple job boards
- Export options: JSON, CSV, database storage

