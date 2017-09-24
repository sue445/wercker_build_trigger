package main

import (
	"encoding/json"
	"github.com/bouk/monkey"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
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
