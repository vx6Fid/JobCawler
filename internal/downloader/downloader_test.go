package downloader

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gocolly/colly/v2"
)

func TestDownloader_FetchWithParser(t *testing.T) {
	// Create a fake HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
		<html>
		  <body>
		    <h1 class="title">Go Developer</h1>
		    <div class="company">Google</div>
		  </body>
		</html>
		`))
	}))
	defer server.Close()

	d := NewDownloader()

	var jobTitle, company string

	err := d.FetchWithParser(server.URL, func(e *colly.HTMLElement) {
		jobTitle = e.DOM.Find(".title").Text()
		company = e.DOM.Find(".company").Text()
	})

	if err != nil {
		t.Fatalf("FetchWithParser failed: %v", err)
	}

	if jobTitle != "Go Developer" {
		t.Errorf("Expected job title 'Go Developer', got '%s'", jobTitle)
	}
	if company != "Google" {
		t.Errorf("Expected company 'Google', got '%s'", company)
	}
}
