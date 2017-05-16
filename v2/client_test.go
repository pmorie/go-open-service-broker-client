package v2

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
)

const okCatalogBytes = `{
  "services": [{
    "name": "fake-service",
    "id": "acb56d7c-XXXX-XXXX-XXXX-feb140a59a66",
    "description": "fake service",
    "tags": ["tag1", "tag2"],
    "requires": ["route_forwarding"],
    "bindable": true,
    "metadata": {
    	"a": "b",
    	"c": "d"
    },
    "dashboard_client": {
      "id": "398e2f8e-XXXX-XXXX-XXXX-19a71ecbcf64",
      "secret": "277cabb0-XXXX-XXXX-XXXX-7822c0a90e5d",
      "redirect_uri": "http://localhost:1234"
    },
    "plan_updateable": true,
    "plans": [{
      "name": "fake-plan-1",
      "id": "d3031751-XXXX-XXXX-XXXX-a42377d3320e",
      "description": "description1",
      "metadata": {
      	"b": "c",
      	"d": "e"
      }
    }]
  }]
}`

func okCatalogResponse() *CatalogResponse {
	return &CatalogResponse{
		Services: []Service{
			{
				ID:          "acb56d7c-XXXX-XXXX-XXXX-feb140a59a66",
				Name:        "fake-service",
				Description: "fake service",
				Tags: []string{
					"tag1",
					"tag2",
				},
				Requires: []string{
					"route_forwarding",
				},
				Bindable:      true,
				PlanUpdatable: truePtr(),
				Plans: []Plan{
					{
						ID:          "d3031751-XXXX-XXXX-XXXX-a42377d3320e",
						Name:        "fake-plan-1",
						Description: "description1",
						Metadata: map[string]interface{}{
							"b": "c",
							"d": "e",
						},
					},
				},
				DashboardClient: &DashboardClient{
					ID:          "398e2f8e-XXXX-XXXX-XXXX-19a71ecbcf64",
					Secret:      "277cabb0-XXXX-XXXX-XXXX-7822c0a90e5d",
					RedirectURI: "http://localhost:1234",
				},
				Metadata: map[string]interface{}{
					"a": "b",
					"c": "d",
				},
			},
		},
	}
}

const okCatalog2Bytes = `{
  "services": [{
    "name": "fake-service-2",
    "id": "fake-service-2-id",
    "description": "service-description-2",
    "bindable": false,
    "plans": [{
      "name": "fake-plan-2",
      "id": "fake-plan-2-id",
      "description": "description-2",
      "bindable": true
    }]
  }]
}`

func okCatalog2Response() *CatalogResponse {
	return &CatalogResponse{
		Services: []Service{
			{
				ID:          "fake-service-2-id",
				Name:        "fake-service-2",
				Description: "service-description-2",
				Bindable:    false,
				Plans: []Plan{
					{
						ID:          "fake-plan-2-id",
						Name:        "fake-plan-2",
						Description: "description-2",
						Bindable:    truePtr(),
					},
				},
			},
		},
	}
}

func truePtr() *bool {
	b := true
	return &b
}

func falsePtr() *bool {
	b := false
	return &b
}

func TestGetCatalog(t *testing.T) {
	cases := []struct {
		name             string
		responseBody     string
		expectedResponse *CatalogResponse
		expectedErr      error
	}{
		{
			name:             "success 1",
			responseBody:     okCatalogBytes,
			expectedResponse: okCatalogResponse(),
		},
		{
			name:             "success 2",
			responseBody:     okCatalog2Bytes,
			expectedResponse: okCatalog2Response(),
		},
	}

	for _, tc := range cases {
		doGetCatalogTest(t, tc.name, tc.responseBody, tc.expectedResponse, tc.expectedErr)
	}
}

func doGetCatalogTest(
	t *testing.T,
	name, responseBody string,
	expectedResponse *CatalogResponse,
	expectedErr error,
) {
	router := mux.NewRouter()
	router.HandleFunc("/v2/catalog", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		bodyBytes := []byte(responseBody)
		_, err := w.Write(bodyBytes)
		if err != nil {
			t.Errorf("%v: error writing response bytes: %v", name, err)
		}
	})

	server := httptest.NewServer(router)
	URL := server.URL
	defer server.Close()

	config := DefaultClientConfiguration()
	config.URL = URL

	client, err := NewClient(config)
	if err != nil {
		t.Errorf("%v: error creating client: %v", name, err)
		return
	}

	response, err := client.GetCatalog()
	if err != nil {
		t.Errorf("%v: error getting catalog: %v", name, err)
		return
	}

	if e, a := expectedResponse, response; !reflect.DeepEqual(e, a) {
		t.Errorf("%v: unexpected diff in catalog response; expected %+v, got %+v", name, e, a)
		return
	}
}
