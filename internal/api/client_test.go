package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testOrderNum                = "123"
	testOrderWrongNum           = "456"
	testOrderNumTooManyRequests = "789"
	testTimeout                 = time.Duration(1) * time.Second
	testRetryAfterValue         = time.Duration(60) * time.Second
)

func TestOrderUpdateClient_GetOrderFromAccrual(t *testing.T) {
	testOrder := &models.OrderAccrual{
		Number:  testOrderNum,
		Status:  "some status",
		Accrual: 10,
	}
	testOrderJSON, err := json.Marshal(testOrder)
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == fmt.Sprintf("/api/orders/%s", testOrderNum) {
			w.WriteHeader(http.StatusOK)
			w.Write(testOrderJSON)
		}
		if r.URL.Path == fmt.Sprintf("/api/orders/%s", testOrderWrongNum) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.URL.Path == fmt.Sprintf("/api/orders/%s", testOrderNumTooManyRequests) {
			w.Header().Set("Retry-After", strings.TrimSuffix(testRetryAfterValue.String(), "s"))
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	}))

	defer server.Close()

	restyClient := resty.New()
	client := NewOrderUpdateClient(restyClient, server.URL, testTimeout)

	t.Run("valid test", func(t *testing.T) {
		order, err := client.GetOrderFromAccrual(context.Background(), testOrderNum)
		assert.NoError(t, err)
		assert.Equal(t, testOrder, order)
	})

	t.Run("no content", func(t *testing.T) {
		_, err := client.GetOrderFromAccrual(context.Background(), testOrderWrongNum)
		assert.ErrorIs(t, err, ErrNoContent)
	})

	t.Run("too many requests", func(t *testing.T) {
		_, err := client.GetOrderFromAccrual(context.Background(), testOrderNumTooManyRequests)
		assert.ErrorIs(t, err, ErrTooManyRequests)
		e, ok := err.(*errRetryAfter)
		assert.True(t, ok)
		dur := e.GetRetryAfterDuration()
		assert.Equal(t, testRetryAfterValue, dur)
	})

	t.Run("unexpected status code", func(t *testing.T) {
		_, err := client.GetOrderFromAccrual(context.Background(), "")
		assert.Error(t, err)
	})
}
