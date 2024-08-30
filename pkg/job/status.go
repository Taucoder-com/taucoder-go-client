package job

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type JobStatusRequest struct {
	JobIDs []string `json:"job_ids"`
}

type JobStatus string

const (
	JobStatusError      JobStatus = "error"
	JobStatusDone       JobStatus = "done"
	JobStatusInProgress JobStatus = "in-progress"
)

type JobStatusResponseItem struct {
	JobID             string    `json:"job_id"`
	Status            JobStatus `json:"status"`
	InputFilename     *string   `json:"input_filename,omitempty"`
	Message           *string   `json:"message,omitempty"`
	OutputURL         *string   `json:"output_url,omitempty"`
	OutputContentType *string   `json:"output_content_type,omitempty"`
	OutputSize        *int64    `json:"output_size,omitempty"`
}

type JobStatusResponse struct {
	RequestStatus string                  `json:"request_status"`
	Jobs          []JobStatusResponseItem `json:"jobs"`
}

func (c *client) NewJobStatusRequest(jobIDs []string) (*http.Request, error) {
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(JobStatusRequest{
		JobIDs: jobIDs,
	})

	req, err := http.NewRequest("POST", c.BaseURL+"/job-status", &b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+c.APIKeyBase64)

	return req, nil
}
