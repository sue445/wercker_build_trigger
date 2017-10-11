package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	// MaxLimit represents paging per 1 page
	MaxLimit = 20
)

// Wercker API client
type Wercker struct {
	token string
}

// WerckerApplication represents a `application` resource
type WerckerApplication struct {
	ID string `json:"id"`
}

// WerckerRun represents a `run` resource
type WerckerRun struct {
	ID         string          `json:"id"`
	URL        string          `json:"url"`
	CreatedAt  string          `json:"createdAt"`
	CommitHash string          `json:"commitHash"`
	Message    string          `json:"message"`
	Result     string          `json:"result"`
	Status     string          `json:"status"`
	Pipeline   WerckerPipeline `json:"pipeline"`
}

// WerckerPipeline represents a `pipeline` resource
type WerckerPipeline struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// WerckerTriggerNewRunParam represents request parameter of TriggerNewRun API
type WerckerTriggerNewRunParam struct {
	PipelineID string `json:"pipelineId"`
	Branch     string `json:"branch"`
	Message    string `json:"message"`
}

// WerckerError represents error response body
type WerckerError struct {
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
	Message    string `json:"message"`
}

// WerckerTrigger represents interface of API client (for stubbing from test)
type WerckerTrigger interface {
	FindPipeline(appPath string, pipelineName string) (pipeline *WerckerPipeline, err error)
	TriggerNewRun(pipelineID string, branch string) (run *WerckerRun, err error)
}

// NewWercker returns new Wercker object
func NewWercker(token string) *Wercker {
	w := new(Wercker)
	w.token = token
	return w
}

// GetApplication gets `application` with specified application path
func (w *Wercker) GetApplication(applicationPath string) (app *WerckerApplication, err error) {
	req, err := http.NewRequest(
		"GET",
		"https://app.wercker.com/api/v3/applications/"+applicationPath,
		nil,
	)
	if err != nil {
		return nil, err
	}

	body, err := w.execute(req)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &app)
	return app, err
}

// GetRuns gets `runs` with specified application
func (w *Wercker) GetRuns(applicationID string, skip int) (runs []WerckerRun, err error) {
	req, err := http.NewRequest(
		"GET",
		"https://app.wercker.com/api/v3/runs",
		nil,
	)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Add("applicationId", applicationID)
	values.Add("skip", strconv.Itoa(skip))
	values.Add("limit", strconv.Itoa(MaxLimit))
	req.URL.RawQuery = values.Encode()

	body, err := w.execute(req)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &runs)
	return runs, err
}

// FindPipeline find `pipeline` with specified application path and pipeline name
func (w *Wercker) FindPipeline(applicationPath string, pipelineName string) (pipeline *WerckerPipeline, err error) {
	application, err := w.GetApplication(applicationPath)
	if err != nil {
		return nil, err
	}

	skip := 0

	for {
		runs, err := w.GetRuns(application.ID, skip)
		if err != nil {
			return nil, err
		}

		for _, run := range runs {
			if run.Pipeline.Name == pipelineName {
				return &run.Pipeline, nil
			}
		}

		if len(runs) < MaxLimit {
			return nil, fmt.Errorf("Not Found pipeline: %s", pipelineName)
		}

		skip += MaxLimit
	}
}

// TriggerNewRun trigger new run with specified pipeline id and branch
func (w *Wercker) TriggerNewRun(pipelineId string, branch string) (run *WerckerRun, err error) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	params := WerckerTriggerNewRunParam{
		PipelineID: pipelineId,
		Branch:     branch,
		Message:    "wercker_build_trigger: " + currentTime,
	}

	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		"https://app.wercker.com/api/v3/runs/",
		bytes.NewBuffer(jsonBytes),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	body, err := w.execute(req)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &run)
	return run, err
}

func (w *Wercker) execute(req *http.Request) ([]byte, error) {
	client := &http.Client{}

	req.Header.Add("Authorization", "Bearer "+w.token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		// parse error response
		e := WerckerError{}
		json.Unmarshal(body, &e)
		return nil, fmt.Errorf("statusCode: %d, error: %s, message: %s", e.StatusCode, e.Error, e.Message)
	}

	return body, nil
}
