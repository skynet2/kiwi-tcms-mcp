package kiwi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLookupFilters_Success(t *testing.T) {
	type tc struct {
		name       string
		wantMethod string
		reply      any
		call       func(c *Client) (any, error)
	}
	cases := []tc{
		{"category", "Category.filter",
			[]map[string]any{{"id": float64(1), "name": "UI", "product": float64(2)}},
			func(c *Client) (any, error) {
				return c.CategoryFilter(context.Background(), map[string]any{"product": 2})
			}},
		{"priority", "Priority.filter",
			[]map[string]any{{"id": float64(1), "value": "P1"}},
			func(c *Client) (any, error) { return c.PriorityFilter(context.Background(), map[string]any{}) }},
		{"product", "Product.filter",
			[]map[string]any{{"id": float64(1), "name": "X", "classification": float64(1)}},
			func(c *Client) (any, error) { return c.ProductFilter(context.Background(), map[string]any{}) }},
		{"component", "Component.filter",
			[]map[string]any{{"id": float64(1), "name": "web", "product": float64(1)}},
			func(c *Client) (any, error) { return c.ComponentFilter(context.Background(), map[string]any{}) }},
		{"tag", "Tag.filter",
			[]map[string]any{{"id": float64(1), "name": "smoke"}},
			func(c *Client) (any, error) { return c.TagFilter(context.Background(), map[string]any{}) }},
		{"testcase_status", "TestCaseStatus.filter",
			[]map[string]any{{"id": float64(1), "name": "CONFIRMED"}},
			func(c *Client) (any, error) { return c.TestCaseStatusFilter(context.Background(), map[string]any{}) }},
		{"user", "User.filter",
			[]map[string]any{{"id": float64(1), "username": "alice", "email": "a@x"}},
			func(c *Client) (any, error) { return c.UserFilter(context.Background(), map[string]any{}) }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			srv := newTestServer(t, rpcHandler(t, c.wantMethod, c.reply))
			defer srv.Close()
			client, _ := New(Config{BaseURL: srv.URL, Username: "u", Password: "p"})
			got, err := c.call(client)
			require.NoError(t, err)
			assert.NotNil(t, got)
		})
	}
}
