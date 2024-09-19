package helpers

import (
	"gofemart/internal/logger"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetHTTPResponse(t *testing.T) {
	cases := []struct {
		name             string
		mockStatus       int
		mockMessage      string
		expectedResponse string
	}{
		{
			name:             "OK_status",
			mockStatus:       http.StatusOK,
			mockMessage:      "OK",
			expectedResponse: "OK",
		},
		{
			name:             "not_found_status",
			mockStatus:       http.StatusNotFound,
			mockMessage:      "Not Found",
			expectedResponse: "Not Found",
		},
		{
			name:             "internal_server_error_status",
			mockStatus:       http.StatusInternalServerError,
			mockMessage:      "Internal Server Error",
			expectedResponse: "Internal Server Error",
		},
		{
			name:             "empty_message",
			mockStatus:       http.StatusInternalServerError,
			mockMessage:      "",
			expectedResponse: "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			err := SetHTTPResponse(response, c.mockStatus, []byte(c.mockMessage))
			if err != nil {
				t.Errorf(
					"SetHTTPResponse error = %v, wantErr %v",
					err,
					false,
				)
			}
			result := response.Result()
			defer func() {
				if cErr := result.Body.Close(); cErr != nil {
					logger.Log.Warn(cErr)
				}
			}()

			if result.StatusCode != c.mockStatus {
				t.Errorf("Status is not expected. result.StatusCode = %v, want %v", result.StatusCode, c.mockStatus)
			}
			body, err := io.ReadAll(result.Body)
			if err != nil {
				t.Errorf("Read boay error is not expected. error = %v, wantErr %v", err, false)
			}
			if c.expectedResponse != string(body) {
				t.Errorf("Body is not expected. body = %v, want %v", string(body), c.expectedResponse)
			}
		})
	}
}
