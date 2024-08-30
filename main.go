package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Taucoder-com/taucoder-go-client/pkg/job"
)

var (
	gitCommit string = "unknown"
	gitTag    string = "unknown"
)

func downloadFile(url string, filepath string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func directoryExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func main() {

	fmt.Printf("taucoder go client version %s %s\n", gitCommit, gitTag)

	apiKey := flag.String("apikey", "", "API key for authentication")
	output := flag.String("output", "", "Output file name")
	quality := flag.Int("quality", 50, "Quality of the output image")

	flag.Parse()
	inputs := flag.Args()

	if quality == nil || *quality < 25 || *quality > 95 {
		log.Fatal("Quality must be between 25 and 95")
	}

	if apiKey == nil || *apiKey == "" {
		log.Fatal("API key is required")
	}

	if isdir, err := directoryExists(*output); !isdir || err != nil {
		log.Fatal("Output directory does not exist")
	}

	if len(inputs) == 0 {
		log.Fatal("At least one input file is required")
	}

	filenames := []string{}
	contentTypes := []string{}

	for _, input := range inputs {
		contentType, err := job.GetFileMimeType(input)
		if err != nil {
			log.Fatalf("Failed to get file mime type: %v", err)
		}
		filenames = append(filenames, input)
		contentTypes = append(contentTypes, contentType)
	}

	client := job.NewClient(*apiKey)

	req, err := client.NewJobCreateRequest(*quality, filenames, contentTypes)
	if err != nil {
		log.Fatalf("Failed to create multipart request: %v", err)
	}

	log.Printf("Sending request...\n")
	status, err := job.DoRequest(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}

	log.Printf("Request status: %s\n", status.RequestStatus)
	log.Printf("%d jobs created\n", len(status.Jobs))
	for _, j := range status.Jobs {
		log.Printf("%s\n", j.JobID)
	}
	log.Printf("waiting for jobs to complete...\n")

	downloadedJobs := map[string]bool{}

	for len(downloadedJobs) < len(status.Jobs) {
		time.Sleep(5 * time.Second)

		var jobIDs []string
		for _, j := range status.Jobs {
			jobIDs = append(jobIDs, j.JobID)
		}

		req, err = client.NewJobStatusRequest(jobIDs)
		if err != nil {
			log.Fatalf("Failed to create status request: %v", err)
		}
		status, err := job.DoRequest(req)
		if err != nil {
			log.Fatalf("Failed to get job status: %v", err)
		}

		for _, j := range status.Jobs {
			if _, ok := downloadedJobs[j.JobID]; ok {
				continue
			}

			switch j.Status {
			case job.JobStatusError:
				log.Printf("Job %s failed: %s", j.JobID, *j.Message)
				downloadedJobs[j.JobID] = true
			case job.JobStatusDone:
				downloadedJobs[j.JobID] = true
				outputPath := filepath.Join(*output, fmt.Sprintf("%s-%s.jpg", *j.InputFilename, j.JobID))
				err = downloadFile(*j.OutputURL, outputPath)
				if err != nil {
					log.Printf("Job: %s failed to download file: %v", j.JobID, err)
				}
				log.Printf("Job %s done => %s", j.JobID, outputPath)
			default:
			}
		}
	}
}
