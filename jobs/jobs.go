package jobs

import (
	"time"
)

type (
	// Represents a job, but without the specific details of the work to be processed.
	// This is more for recording job metadata such as sending an email where
	// details such as an S3 bucket or file key are not necessary.
	BaseJob struct {
		CreatedAt   time.Time `json:"createdAt"`
		ID          string    `json:"id"`
		UserID      string    `json:"userId"`
		Email       string    `json:"email"`
		Type        string    `json:"type"`
		Description string    `json:"description,omitempty"`
	}

	// Job represents a unit of work to be processed.
	Job struct {
		BaseJob
		Bucket      string `json:"bucket"`
		Key         string `json:"key"`
		Filename    string `json:"filename"`
		ContentType string `json:"contentType"`
		Size        int64  `json:"size"`
	}
)
