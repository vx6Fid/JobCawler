```
/crawler-backend
  /cmd
    /crawler
      main.go                # Application entrypoint; sets up and runs crawler
  /internal
    /frontend               # URL Frontier (queue, scheduler)
    /downloader             # HTTP fetching logic
    /parser                 # HTML parsing and content extraction
    /storage
      /content              # Content storage (DB, files)
      /url                  # URL metadata storage (seen URLs, status)
    /dedup                  # Deduplication logic for content and URLs
    /extractor              # Link extraction from HTML
    /filter                 # URL filtering rules and logic
    /config                 # Configuration loading/parsing
    /logger                 # Logging setup and utilities
    /worker                 # Coordinating crawling workers, concurrency management
  /pkg
    /model                  # Data models and structs shared across packages
    /utils                  # Utility helpers, e.g. string, time, error handling
  /scripts                  # Automation scripts (DB migration, setup)
  /docs                     # Documentation about your crawler design & usage
  go.mod
  go.sum
  README.md

```