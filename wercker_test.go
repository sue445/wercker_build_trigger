package main

import (
	"encoding/json"
	"fmt"
	"github.com/bouk/monkey"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestWercker_TriggerNewRun(t *testing.T) {
	token := "api_token"
	pipelineID := "123456789012345678901234"
	branch := "develop"

	// stub Time.Now()
	currentTime := time.Date(2017, time.September, 12, 1, 2, 3, 4, time.UTC)
	patch := monkey.Patch(time.Now, func() time.Time { return currentTime })
	defer patch.Unpatch()

	// mock http POST
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	runURL := "https://app.wercker.com/api/v3/runs/577a36013828f9673b035730"
	message := "wercker_build_trigger: 2017-09-12 01:02:03"

	httpmock.RegisterResponder(
		"POST", "https://app.wercker.com/api/v3/runs/",
		func(req *http.Request) (*http.Response, error) {
			param := WerckerTriggerNewRunParam{}
			if err := json.NewDecoder(req.Body).Decode(&param); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			assert.Equal(t, pipelineID, param.PipelineID)
			assert.Equal(t, branch, param.Branch)
			assert.Equal(t, message, param.Message)

			response := WerckerRun{URL: runURL, Message: param.Message}

			resp, err := httpmock.NewJsonResponse(200, response)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	wercker := NewWercker(token)

	ret, err := wercker.TriggerNewRun(pipelineID, branch)

	assert.NoError(t, err)

	assert.Equal(t, runURL, ret.URL)
	assert.Equal(t, message, ret.Message)
}

func TestWercker_GetApplication(t *testing.T) {
	token := "api_token"
	applicationPath := "wercker/docs"
	applicationID := "54c9168980c7075225004157"

	// mock http GET
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	url := fmt.Sprintf("https://app.wercker.com/api/v3/applications/%s", applicationPath)

	httpmock.RegisterResponder(
		"GET", url,
		httpmock.NewStringResponder(200, readFile("test/GetApplication.json")),
	)

	wercker := NewWercker(token)

	ret, err := wercker.GetApplication(applicationPath)

	assert.NoError(t, err)
	assert.Equal(t, applicationID, ret.ID)
}

func TestWercker_GetRuns(t *testing.T) {
	token := "api_token"
	applicationId := "54c9168980c7075225004157"
	skip := 0

	// mock http GET
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	url := fmt.Sprintf("https://app.wercker.com/api/v3/runs?applicationId=%s&limit=%d&skip=%d", applicationId, MaxLimit, skip)

	httpmock.RegisterResponder(
		"GET", url,
		httpmock.NewStringResponder(200, readFile("test/GetRuns.json")),
	)

	wercker := NewWercker(token)

	ret, err := wercker.GetRuns(applicationId, skip)

	assert.NoError(t, err)
	assert.Equal(t, "588a61d30a002301003b44d5", ret[0].ID)
	assert.Equal(t, "54c9168980c7075225004157", ret[0].Pipeline.ID)
	assert.Equal(t, "build", ret[0].Pipeline.Name)
}

func TestWercker_FindPipeline(t *testing.T) {
	token := "api_token"
	applicationPath := "wercker/docs"
	pipelineName := "build"
	applicationId := "54c9168980c7075225004157"
	skip := 0

	// mock http GET
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET", fmt.Sprintf("https://app.wercker.com/api/v3/applications/%s", applicationPath),
		httpmock.NewStringResponder(200, readFile("test/GetApplication.json")),
	)

	httpmock.RegisterResponder(
		"GET", fmt.Sprintf("https://app.wercker.com/api/v3/runs?applicationId=%s&limit=%d&skip=%d", applicationId, MaxLimit, skip),
		httpmock.NewStringResponder(200, readFile("test/GetRuns.json")),
	)

	wercker := NewWercker(token)

	pipeline, err := wercker.FindPipeline(applicationPath, pipelineName)

	assert.NoError(t, err)
	assert.Equal(t, "54c9168980c7075225004157", pipeline.ID)
	assert.Equal(t, "build", pipeline.Name)
}

func readFile(fileName string) string {
	buf, err := ioutil.ReadFile(fileName)

	if err != nil {
		panic(err)
	}

	return string(buf)
}
