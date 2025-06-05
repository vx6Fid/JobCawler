package pkg

import (
	"testing"
	"time"
)

func TestUpsertJob(t *testing.T) {
	err := ConnectMongo()
	if err != nil {
		t.Fatalf("MongoDB connection failed: %v", err)
	}

	job := JobPosting{
		Title:       "DevOps Engineer",
		Company:     "Chipcolate",
		Location:    "Remote",
		Salary:      "$75,000 - $99,999 USD",
		PostedOn:    time.Now().AddDate(0, 0, -1),
		Description: "Work with Docker, Kubernetes, Terraform...",
		URL:         "https://weworkremotely.com/jobs/devops-engineer",
		Source:      "WeWorkRemotely",
		ApplyURL:    "https://apply-link.com",
		Skills:      []string{"Docker", "Kubernetes", "Terraform"},
		Experience:  "3+ years",
	}

	err = UpsertJob(job)
	if err != nil {
		t.Fatalf("Failed to upsert job: %v", err)
	}
}
