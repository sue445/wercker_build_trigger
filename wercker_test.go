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
	pipelineId := "123456789012345678901234"
	branch := "develop"

	// stub Time.Now()
	currentTime := time.Date(2017, time.September, 12, 1, 2, 3, 4, time.UTC)
	patch := monkey.Patch(time.Now, func() time.Time { return currentTime })
	defer patch.Unpatch()

	// mock http POST
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	runUrl := "https://app.wercker.com/api/v3/runs/577a36013828f9673b035730"
	message := "wercker_build_trigger: 2017-09-12 01:02:03"

	httpmock.RegisterResponder(
		"POST", "https://app.wercker.com/api/v3/runs/",
		func(req *http.Request) (*http.Response, error) {
			param := WerckerTriggerNewRunParam{}
			if err := json.NewDecoder(req.Body).Decode(&param); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			assert.Equal(t, pipelineId, param.PipelineId)
			assert.Equal(t, branch, param.Branch)
			assert.Equal(t, message, param.Message)

			response := WerckerRun{Url: runUrl, Message: param.Message}

			resp, err := httpmock.NewJsonResponse(200, response)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	wercker := NewWercker(token)

	ret, err := wercker.TriggerNewRun(pipelineId, branch)

	assert.NoError(t, err)

	assert.Equal(t, runUrl, ret.Url)
	assert.Equal(t, message, ret.Message)
}

func TestWercker_GetApplication(t *testing.T) {
	token := "api_token"
	appPath := "wercker/docs"
	appId := "54c9168980c7075225004157"

	// mock http GET
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	url := fmt.Sprintf("https://app.wercker.com/api/v3/applications/%s", appPath)

	httpmock.RegisterResponder(
		"GET", url,
		httpmock.NewStringResponder(200, readFile("test/GetApplication.json")),
	)

	wercker := NewWercker(token)

	ret, err := wercker.GetApplication(appPath)

	assert.NoError(t, err)
	assert.Equal(t, appId, ret.Id)
}

func TestWercker_GetRuns(t *testing.T) {
	token := "api_token"
	appId := "54c9168980c7075225004157"
	skip := 0

	// mock http GET
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	url := fmt.Sprintf("https://app.wercker.com/api/v3/runs?applicationId=%s&skip=%d", appId, skip)

	httpmock.RegisterResponder(
		"GET", url,
		httpmock.NewStringResponder(200, readFile("test/GetRuns.json")),
	)

	wercker := NewWercker(token)

	ret, err := wercker.GetRuns(appId, skip)

	assert.NoError(t, err)
	assert.Equal(t, "588a61d30a002301003b44d5", ret[0].Id)
	assert.Equal(t, "54c9168980c7075225004157", ret[0].Pipeline.Id)
	assert.Equal(t, "build", ret[0].Pipeline.Name)
}

func readFile(fileName string) string {
	buf, err := ioutil.ReadFile(fileName)

	if err != nil {
		panic(err)
	}

	return string(buf)
}
