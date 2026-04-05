package kiwi

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type rpcRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
	ID      int    `json:"id"`
}

type rpcResponse struct {
	Jsonrpc string    `json:"jsonrpc"`
	Result  any       `json:"result,omitempty"`
	Error   *RPCError `json:"error,omitempty"`
	ID      int       `json:"id"`
}

func newTestServer(t *testing.T, h http.HandlerFunc) *httptest.Server {
	t.Helper()
	return httptest.NewServer(h)
}

func TestClient_Do_Success(t *testing.T) {
	type tc struct {
		name   string
		method string
		params any
		reply  any
		want   map[string]any
	}
	cases := []tc{
		{
			name:   "simple object result",
			method: "TestCase.filter",
			params: []any{map[string]any{"pk": 1}},
			reply:  map[string]any{"id": float64(1), "summary": "case"},
			want:   map[string]any{"id": float64(1), "summary": "case"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/json-rpc/", r.URL.Path)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				body, _ := io.ReadAll(r.Body)
				var req rpcRequest
				require.NoError(t, json.Unmarshal(body, &req))
				assert.Equal(t, c.method, req.Method)
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(rpcResponse{Jsonrpc: "2.0", Result: c.reply, ID: req.ID})
			})
			defer srv.Close()

			client, err := New(Config{
				BaseURL: srv.URL, Username: "u", Password: "p",
				HTTPClient: &http.Client{Timeout: 5 * time.Second},
			})
			require.NoError(t, err)

			var got map[string]any
			err = client.Do(context.Background(), c.method, c.params, &got)
			require.NoError(t, err)
			assert.Equal(t, c.want, got)
		})
	}
}

func TestClient_Login_Success(t *testing.T) {
	var calls int
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		body, _ := io.ReadAll(r.Body)
		var req rpcRequest
		require.NoError(t, json.Unmarshal(body, &req))
		assert.Equal(t, "Auth.login", req.Method)
		params, _ := req.Params.([]any)
		require.Len(t, params, 2)
		assert.Equal(t, "user", params[0])
		assert.Equal(t, "pass", params[1])
		http.SetCookie(w, &http.Cookie{Name: "sessionid", Value: "abc123"})
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(rpcResponse{Jsonrpc: "2.0", Result: "ok", ID: req.ID})
	})
	defer srv.Close()

	client, err := New(Config{BaseURL: srv.URL, Username: "user", Password: "pass"})
	require.NoError(t, err)

	require.NoError(t, client.Login(context.Background()))
	assert.Equal(t, 1, calls)
}

type scriptedHandler struct {
	steps []http.HandlerFunc
	idx   atomic.Int32
	calls []string
}

func (s *scriptedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var req rpcRequest
	_ = json.Unmarshal(body, &req)
	s.calls = append(s.calls, req.Method)
	r.Body = io.NopCloser(bytes.NewReader(body))
	i := int(s.idx.Add(1)) - 1
	s.steps[i](w, r)
}

func unauthorized() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusUnauthorized) }
}

func loginOK() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req rpcRequest
		_ = json.Unmarshal(body, &req)
		http.SetCookie(w, &http.Cookie{Name: "sessionid", Value: "fresh"})
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(rpcResponse{Jsonrpc: "2.0", Result: "ok", ID: req.ID})
	}
}

func rpcOK(reply any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req rpcRequest
		_ = json.Unmarshal(body, &req)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(rpcResponse{Jsonrpc: "2.0", Result: reply, ID: req.ID})
	}
}

func TestClient_Do_RetriesOn401(t *testing.T) {
	s := &scriptedHandler{steps: []http.HandlerFunc{unauthorized(), loginOK(), rpcOK([]any{})}}
	srv := httptest.NewServer(s)
	defer srv.Close()

	client, err := New(Config{BaseURL: srv.URL, Username: "u", Password: "p"})
	require.NoError(t, err)

	var result []any
	require.NoError(t, client.Do(context.Background(), "TestCase.filter", []any{map[string]any{}}, &result))
	assert.Equal(t, []string{"TestCase.filter", "Auth.login", "TestCase.filter"}, s.calls)
}

func TestClient_Do_Failure_Unauthorized_Persistent(t *testing.T) {
	s := &scriptedHandler{steps: []http.HandlerFunc{unauthorized(), loginOK(), unauthorized()}}
	srv := httptest.NewServer(s)
	defer srv.Close()

	client, err := New(Config{BaseURL: srv.URL, Username: "u", Password: "p"})
	require.NoError(t, err)

	var result []any
	err = client.Do(context.Background(), "TestCase.filter", []any{}, &result)
	require.ErrorIs(t, err, ErrUnauthorized)
}
