package serviceaclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_ForwardCEP(t *testing.T) {
	tests := []struct {
		name           string
		serverStatus   int
		serverBody     string
		wantStatusCode int
	}{
		{name: "success", serverStatus: http.StatusOK, serverBody: `{"city":"São Paulo","temp_C":28.5,"temp_F":83.3,"temp_K":301.5}`, wantStatusCode: http.StatusOK},
		{name: "invalid zipcode", serverStatus: http.StatusUnprocessableEntity, serverBody: `{"message":"invalid zipcode"}`, wantStatusCode: http.StatusUnprocessableEntity},
		{name: "not found", serverStatus: http.StatusNotFound, serverBody: `{"message":"can not find zipcode"}`, wantStatusCode: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.serverStatus)
				w.Write([]byte(tt.serverBody))
			}))
			defer server.Close()

			client := Client{BaseURL: server.URL, HTTPClient: server.Client()}
			resp, err := client.ForwardCEP(context.Background(), "29902555")

			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatusCode, resp.StatusCode)
			assert.JSONEq(t, tt.serverBody, string(resp.Body))
		})
	}
}
