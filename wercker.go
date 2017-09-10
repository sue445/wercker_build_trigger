package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	DEFAULT_BRANCH = "master"
)

type Wercker struct {
	token string
}

type WerckerApplication struct {
	Id string `json:"id"`
}

type WerckerRun struct {
	Url        string `json:"url"`
	CreatedAt  string `json:"createdAt"`
	CommitHash string `json:"commitHash"`
	Message    string `json:"message"`
	Result     string `json:"result"`
	Status     string `json:"status"`
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

func NewWercker(token string) *Wercker {
	w := new(Wercker)
	w.token = token
	return w
}

func (w *Wercker) TriggerNewRun(pipelineId string, branch string) (run *WerckerRun, err error) {
	if len(branch) == 0 {
		branch = DEFAULT_BRANCH
	}

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
