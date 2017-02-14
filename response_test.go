package leanplum

import (
	"net/http"
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestResponse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://www.leanplum.com/api",
		httpmock.NewStringResponder(200, `{
            "response": [
                {
                    "success": true
                }
            ]
        }`))

	url := Leanplum_api_url

	response, err := http.Get(url)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	commonResponser := new(CommonResponse)
	slice, _ := commonResponser.processResponse(response)

	if !slice[0].Success {
		t.Errorf("Success must be true")
	}

}
