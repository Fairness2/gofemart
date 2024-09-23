package accrual

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"gofemart/internal/models"
	"gofemart/internal/payloads"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAccrual(t *testing.T) {
	cases := []struct {
		name       string
		status     int
		response   *payloads.Accrual
		retryAfter time.Duration
		hasRetry   bool
		err        error
	}{
		{
			name:   "no_content",
			status: http.StatusNoContent,
			err:    ErrorOrderNotRegistered,
		},
		{
			name:   "internal_server_error",
			status: http.StatusInternalServerError,
			err:    ErrorInternalAccrual,
		},
		{
			name:       "too_many_requests_without_retry_after",
			status:     http.StatusTooManyRequests,
			retryAfter: time.Minute,
			err:        &TooManyRequestError{InternalError: ErrorTooManyRequests, PauseDuration: time.Minute},
		},
		{
			name:       "too_many_requests_with_retry_after",
			status:     http.StatusTooManyRequests,
			retryAfter: 5 * time.Minute,
			hasRetry:   true,
			err:        &TooManyRequestError{InternalError: ErrorTooManyRequests, PauseDuration: 5 * time.Minute},
		},
		{
			name:   "unknown_error",
			status: http.StatusNotFound,
			err:    ErrorUnknownStatusRequests,
		},
		{
			name:   "ok",
			status: http.StatusOK,
			response: &payloads.Accrual{
				Order: "1",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			router := chi.NewRouter()
			router.Get(getOrderURL+"1", func(writer http.ResponseWriter, request *http.Request) {
				if tc.hasRetry {
					writer.Header().Set("Retry-After", tc.retryAfter.String())
				}
				writer.WriteHeader(tc.status)
				if tc.response != nil {
					b, err := json.Marshal(tc.response)
					if err != nil {
						t.Fatal(err)
					}
					if _, err = writer.Write(b); err != nil {
						t.Fatal(err)
					}
				}
			})
			server := httptest.NewServer(router)
			defer server.Close()

			client := resty.New().SetBaseURL(server.URL)
			proxy := &Proxy{
				pauseDuration: time.Minute,
				client:        client,
			}
			order := &models.Order{
				Number: "1",
			}
			res, err := proxy.Accrual(order)
			if tc.err != nil {
				var tooManyRequestError *TooManyRequestError
				if errors.As(err, &tooManyRequestError) && tc.retryAfter != 0 {
					if tooManyRequestError.PauseDuration != tc.retryAfter {
						t.Errorf("expected pause duration %v, got %v", tc.retryAfter, tooManyRequestError.PauseDuration)
					}
				} else if !errors.Is(err, tc.err) {
					t.Errorf("expected error %v, got %v", tc.err, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error %v", err)
				} else {
					if res == nil || res.Order != tc.response.Order {
						t.Errorf("expected response %v, got %v", tc.response, res)
					}
				}
			}
		})
	}
}
