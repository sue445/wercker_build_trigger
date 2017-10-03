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
	MAX_LIMIT = 20
)

type Wercker struct {
	token string
}

type WerckerApplication struct {
	Id string `json:"id"`
}

type WerckerRun struct {
	Id         string          `json:"id"`
	Url        string          `json:"url"`
	CreatedAt  string          `json:"createdAt"`
	CommitHash string          `json:"commitHash"`
	Message    string          `json:"message"`
	Result     string          `json:"result"`
	Status     string          `json:"status"`
	Pipeline   WerckerPipeline `json:"pipeline"`
}

type WerckerPipeline struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type WerckerTriggerNewRunParam struct {
	PipelineId string `json:"pipelineId"`
	Branch     string `json:"branch"`
	Message    string `json:"message"`
}

type WerckerError struct {
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
	Message    string `json:"message"`
}

type WerckerTrigger interface {
	FindPipeline(appPath string, pipelineName string) (pipeline *WerckerPipeline, err error)
	TriggerNewRun(pipelineId string, branch string) (run *WerckerRun, err error)
}

func NewWercker(token string) *Wercker {
	w := new(Wercker)
	w.token = token
	return w
}

func (w *Wercker) GetApplication(appPath string) (app *WerckerApplication, err error) {
	req, err := http.NewRequest(
		"GET",
		"https://app.wercker.com/api/v3/applications/"+appPath,
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

func (w *Wercker) GetRuns(applicationId string, skip int) (runs []WerckerRun, err error) {
	req, err := http.NewRequest(
		"GET",
		"https://app.wercker.com/api/v3/runs",
		nil,
	)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Add("applicationId", applicationId)
	values.Add("skip", strconv.Itoa(skip))
	req.URL.RawQuery = values.Encode()

	body, err := w.execute(req)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &runs)
	return runs, err
}

func (w *Wercker) FindPipeline(appPath string, pipelineName string) (pipeline *WerckerPipeline, err error) {
	application, err := w.GetApplication(appPath)
	if err != nil {
		return nil, err
	}

	skip := 0

	for {
		runs, err := w.GetRuns(application.Id, skip)
		if err != nil {
			return nil, err
		}

		for _, run := range runs {
			if run.Pipeline.Name == pipelineName {
				return &run.Pipeline, nil
			}
		}

		if len(runs) < MAX_LIMIT {
			return nil, fmt.Errorf("Not Found pipeline: %s", pipelineName)
		}

		skip += MAX_LIMIT
	}
}

func (w *Wercker) TriggerNewRun(pipelineId string, branch string) (run *WerckerRun, err error) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	params := WerckerTriggerNewRunParam{
		PipelineId: pipelineId,
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
