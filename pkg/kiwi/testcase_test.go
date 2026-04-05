package kiwi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func rpcHandler(t *testing.T, wantMethod string, reply any) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req rpcRequest
		require.NoError(t, json.Unmarshal(body, &req))
		assert.Equal(t, wantMethod, req.Method)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(rpcResponse{Jsonrpc: "2.0", Result: reply, ID: req.ID})
	}
}

func TestTestCaseFilter_Success(t *testing.T) {
	srv := newTestServer(t, rpcHandler(t, "TestCase.filter",
		[]map[string]any{{"id": float64(1), "summary": "a"}}))
	defer srv.Close()
	client, _ := New(Config{BaseURL: srv.URL, Username: "u", Password: "p"})
	cases, err := client.TestCaseFilter(context.Background(), map[string]any{"pk": 1})
	require.NoError(t, err)
	assert.Len(t, cases, 1)
	assert.Equal(t, int64(1), cases[0].ID)
	assert.Equal(t, "a", cases[0].Summary)
}

func TestTestCaseCreate_Success(t *testing.T) {
	srv := newTestServer(t, rpcHandler(t, "TestCase.create",
		map[string]any{"id": float64(42), "summary": "new"}))
	defer srv.Close()
	client, _ := New(Config{BaseURL: srv.URL, Username: "u", Password: "p"})
	tc, err := client.TestCaseCreate(context.Background(), TestCaseCreateInput{
		Summary: "new", Category: 1, Priority: 2, CaseStatus: 3,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(42), tc.ID)
}

func TestTestCaseUpdate_Success(t *testing.T) {
	srv := newTestServer(t, rpcHandler(t, "TestCase.update",
		map[string]any{"id": float64(1), "summary": "updated"}))
	defer srv.Close()
	client, _ := New(Config{BaseURL: srv.URL, Username: "u", Password: "p"})
	tc, err := client.TestCaseUpdate(context.Background(), 1, map[string]any{"summary": "updated"})
	require.NoError(t, err)
	assert.Equal(t, "updated", tc.Summary)
}

func TestTestCaseRemove_Success(t *testing.T) {
	srv := newTestServer(t, rpcHandler(t, "TestCase.remove", nil))
	defer srv.Close()
	client, _ := New(Config{BaseURL: srv.URL, Username: "u", Password: "p"})
	require.NoError(t, client.TestCaseRemove(context.Background(), map[string]any{"pk": 1}))
}

func TestTestCaseRelations_Success(t *testing.T) {
	type tc struct {
		name   string
		method string
		call   func(c *Client) error
	}
	cases := []tc{
		{"add_tag", "TestCase.add_tag", func(c *Client) error {
			return c.TestCaseAddTag(context.Background(), 1, "smoke")
		}},
		{"remove_tag", "TestCase.remove_tag", func(c *Client) error {
			return c.TestCaseRemoveTag(context.Background(), 1, "smoke")
		}},
		{"add_component", "TestCase.add_component", func(c *Client) error {
			return c.TestCaseAddComponent(context.Background(), 1, "web")
		}},
		{"remove_component", "TestCase.remove_component", func(c *Client) error {
			return c.TestCaseRemoveComponent(context.Background(), 1, 1)
		}},
		{"add_comment", "TestCase.add_comment", func(c *Client) error {
			return c.TestCaseAddComment(context.Background(), 1, "hello")
		}},
		{"add_property", "TestCase.add_property", func(c *Client) error {
			return c.TestCaseAddProperty(context.Background(), 1, "env", "prod")
		}},
		{"remove_property", "TestCase.remove_property", func(c *Client) error {
			return c.TestCaseRemoveProperty(context.Background(), map[string]any{"case": 1})
		}},
		{"add_link", "TestCase.add_link", func(c *Client) error {
			_, err := c.TestCaseAddLink(context.Background(), map[string]any{"case_id": 1, "url": "http://x"})
			return err
		}},
		{"remove_link", "TestCase.remove_link", func(c *Client) error {
			return c.TestCaseRemoveLink(context.Background(), map[string]any{"case": 1})
		}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			srv := newTestServer(t, rpcHandler(t, c.method, map[string]any{"id": float64(1)}))
			defer srv.Close()
			client, _ := New(Config{BaseURL: srv.URL, Username: "u", Password: "p"})
			require.NoError(t, c.call(client))
		})
	}
}
