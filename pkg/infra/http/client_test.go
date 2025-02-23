package http_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	httpclient "order-system/pkg/infra/http"
	"order-system/pkg/infra/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	cfg := &config.Config{}
	cfg.HTTP.RequestTimeout = 30 * time.Second

	client := httpclient.NewClient(cfg, "http://example.com")
	require.NotNil(t, client)
}

func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "test", r.Header.Get("X-Test"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))
	defer server.Close()

	cfg := &config.Config{}
	cfg.HTTP.RequestTimeout = 30 * time.Second
	cfg.HTTP.MaxRequestSize = 1024

	client := httpclient.NewClient(cfg, server.URL)

	ctx := context.Background()
	resp, err := client.Get(ctx, "/test", &httpclient.RequestOption{
		Timeout:       5 * time.Second,
		RetryCount:    0,
		RetryInterval: time.Millisecond,
		MaxBodySize:   1024,
		Headers:       map[string]string{"X-Test": "test"},
	})

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, []byte("success"), resp.Body)
}

func TestPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		body, _ := io.ReadAll(r.Body)
		assert.Equal(t, []byte(`{"test":"data"}`), body)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("created"))
	}))
	defer server.Close()

	cfg := &config.Config{}
	cfg.HTTP.RequestTimeout = 30 * time.Second
	cfg.HTTP.MaxRequestSize = 2048

	client := httpclient.NewClient(cfg, server.URL)

	ctx := context.Background()
	resp, err := client.Post(ctx, "/test", []byte(`{"test":"data"}`), &httpclient.RequestOption{
		Headers:     map[string]string{"Content-Type": "application/json"},
		MaxBodySize: 2048,
	})

	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, []byte("created"), resp.Body)
}

func TestPut(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		body, _ := io.ReadAll(r.Body)
		assert.Equal(t, []byte(`{"test":"data"}`), body)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := &config.Config{}
	cfg.HTTP.RequestTimeout = 30 * time.Second
	cfg.HTTP.MaxRequestSize = 1024

	client := httpclient.NewClient(cfg, server.URL)

	ctx := context.Background()
	resp, err := client.Put(ctx, "/test", []byte(`{"test":"data"}`), nil)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	cfg := &config.Config{}
	cfg.HTTP.RequestTimeout = 30 * time.Second
	cfg.HTTP.MaxRequestSize = 1024

	client := httpclient.NewClient(cfg, server.URL)

	ctx := context.Background()
	resp, err := client.Delete(ctx, "/test", nil)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestRequestWithLargeResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "2048")
		w.WriteHeader(http.StatusOK)
		w.Write(make([]byte, 2048))
	}))
	defer server.Close()

	cfg := &config.Config{}
	cfg.HTTP.RequestTimeout = 30 * time.Second
	cfg.HTTP.MaxRequestSize = 1024

	client := httpclient.NewClient(cfg, server.URL)

	ctx := context.Background()
	_, err := client.Get(ctx, "/test", &httpclient.RequestOption{
		MaxBodySize: 1024,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "response body too large")
}

func TestRequestWithNetworkError(t *testing.T) {
	cfg := &config.Config{}
	cfg.HTTP.RequestTimeout = 30 * time.Second
	cfg.HTTP.MaxRequestSize = 1024

	client := httpclient.NewClient(cfg, "http://invalid-host")

	ctx := context.Background()
	_, err := client.Get(ctx, "/test", &httpclient.RequestOption{
		RetryCount:    2,
		RetryInterval: time.Millisecond,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "request failed")
}

func TestRequestWithInvalidURL(t *testing.T) {
	cfg := &config.Config{}
	cfg.HTTP.RequestTimeout = 30 * time.Second
	cfg.HTTP.MaxRequestSize = 1024

	client := httpclient.NewClient(cfg, "http://example.com")

	ctx := context.Background()
	_, err := client.Get(ctx, string([]byte{0x7f}), nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create request")
}

func TestDoRequestHeadersAndTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "test-value", r.Header.Get("X-Test-Header"))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := &config.Config{}
	cfg.HTTP.RequestTimeout = 30 * time.Second
	cfg.HTTP.MaxRequestSize = 1024

	client := httpclient.NewClient(cfg, server.URL)

	ctx := context.Background()
	_, err := client.Get(ctx, "/test", &httpclient.RequestOption{
		Headers:       map[string]string{"X-Test-Header": "test-value"},
		Timeout:       time.Millisecond,
	})

	require.NoError(t, err)
}
