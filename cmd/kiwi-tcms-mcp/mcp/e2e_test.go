//go:build e2e

package mcp

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/skynet2/kiwi-tcms-mcp/pkg/kiwi"
)

func e2eClient(t *testing.T) KiwiClient {
	t.Helper()
	if testing.Short() {
		t.Skip("skip E2E in -short mode")
	}
	url := os.Getenv("KIWI_E2E_URL")
	user := os.Getenv("KIWI_E2E_USERNAME")
	pass := os.Getenv("KIWI_E2E_PASSWORD")
	if url == "" || user == "" || pass == "" {
		t.Skip("E2E disabled: set KIWI_E2E_URL, KIWI_E2E_USERNAME, KIWI_E2E_PASSWORD")
	}
	c, err := kiwi.New(kiwi.Config{BaseURL: url, Username: user, Password: pass})
	require.NoError(t, err)
	require.NoError(t, c.Login(context.Background()))
	return c
}

func TestE2E_MCP_LookupTools(t *testing.T) {
	c := e2eClient(t)
	log := zerolog.Nop()

	type tc struct {
		name    string
		handler func(KiwiClient, zerolog.Logger) toolHandler
	}
	cases := []tc{
		{"priority_filter", priorityFilterHandler},
		{"testcase_status_filter", testCaseStatusFilterHandler},
		{"user_filter", userFilterHandler},
		{"product_filter", productFilterHandler},
		{"category_filter", categoryFilterHandler},
		{"component_filter", componentFilterHandler},
		{"tag_filter", tagFilterHandler},
	}
	for _, cse := range cases {
		t.Run(cse.name, func(t *testing.T) {
			h := cse.handler(c, log)
			req := mcplib.CallToolRequest{}
			req.Params.Arguments = map[string]any{"query": map[string]any{}}
			res, err := h(context.Background(), req)
			require.NoError(t, err)
			assert.False(t, res.IsError)
			assert.NotEmpty(t, res.Content)
		})
	}
}

func TestE2E_MCP_TestCaseCreate(t *testing.T) {
	c := e2eClient(t)
	log := zerolog.Nop()

	createH := testCaseCreateHandler(c, log)
	req := mcplib.CallToolRequest{}
	req.Params.Arguments = map[string]any{
		"summary":     "MCP-E2E via tool " + time.Now().Format(time.RFC3339Nano),
		"category":    1.0,
		"priority":    1.0,
		"case_status": 2.0,
	}
	res, err := createH(context.Background(), req)
	require.NoError(t, err)
	require.False(t, res.IsError)

	text := firstTextContent(t, res)
	var created struct {
		ID int64 `json:"id"`
	}
	require.NoError(t, json.Unmarshal([]byte(text), &created))
	require.NotZero(t, created.ID)

	t.Cleanup(func() {
		_ = c.TestCaseRemove(context.Background(), map[string]any{"pk": created.ID})
	})

	t.Run("testcase_get", func(t *testing.T) {
		h := testCaseGetHandler(c, log)
		r := mcplib.CallToolRequest{}
		r.Params.Arguments = map[string]any{"id": float64(created.ID)}
		res, err := h(context.Background(), r)
		require.NoError(t, err)
		assert.False(t, res.IsError)
	})

	t.Run("testcase_update", func(t *testing.T) {
		h := testCaseUpdateHandler(c, log)
		r := mcplib.CallToolRequest{}
		r.Params.Arguments = map[string]any{
			"id":    float64(created.ID),
			"patch": map[string]any{"summary": "updated via mcp handler"},
		}
		res, err := h(context.Background(), r)
		require.NoError(t, err)
		assert.False(t, res.IsError)
	})

	t.Run("testcase_add_tag", func(t *testing.T) {
		h := testCaseAddTagHandler(c, log)
		r := mcplib.CallToolRequest{}
		r.Params.Arguments = map[string]any{"case_id": float64(created.ID), "tag": "mcp-e2e"}
		res, err := h(context.Background(), r)
		require.NoError(t, err)
		assert.False(t, res.IsError)
	})

	t.Run("testcase_remove_tag", func(t *testing.T) {
		h := testCaseRemoveTagHandler(c, log)
		r := mcplib.CallToolRequest{}
		r.Params.Arguments = map[string]any{"case_id": float64(created.ID), "tag": "mcp-e2e"}
		res, err := h(context.Background(), r)
		require.NoError(t, err)
		assert.False(t, res.IsError)
	})

	t.Run("testcase_add_comment", func(t *testing.T) {
		h := testCaseAddCommentHandler(c, log)
		r := mcplib.CallToolRequest{}
		r.Params.Arguments = map[string]any{"case_id": float64(created.ID), "comment": "mcp e2e comment"}
		res, err := h(context.Background(), r)
		require.NoError(t, err)
		assert.False(t, res.IsError)
	})

	t.Run("testcase_add_component", func(t *testing.T) {
		h := testCaseAddComponentHandler(c, log)
		r := mcplib.CallToolRequest{}
		r.Params.Arguments = map[string]any{"case_id": float64(created.ID), "component": "web"}
		res, err := h(context.Background(), r)
		require.NoError(t, err)
		assert.False(t, res.IsError)
	})

	t.Run("testcase_remove_component", func(t *testing.T) {
		h := testCaseRemoveComponentHandler(c, log)
		r := mcplib.CallToolRequest{}
		r.Params.Arguments = map[string]any{"case_id": float64(created.ID), "component_id": 1.0}
		res, err := h(context.Background(), r)
		require.NoError(t, err)
		assert.False(t, res.IsError)
	})

	t.Run("testcase_add_property", func(t *testing.T) {
		h := testCaseAddPropertyHandler(c, log)
		r := mcplib.CallToolRequest{}
		r.Params.Arguments = map[string]any{"case_id": float64(created.ID), "name": "env", "value": "mcp"}
		res, err := h(context.Background(), r)
		require.NoError(t, err)
		assert.False(t, res.IsError)
	})

	t.Run("testcase_remove_property", func(t *testing.T) {
		h := testCaseRemovePropertyHandler(c, log)
		r := mcplib.CallToolRequest{}
		r.Params.Arguments = map[string]any{"query": map[string]any{"case": float64(created.ID), "name": "env"}}
		res, err := h(context.Background(), r)
		require.NoError(t, err)
		assert.False(t, res.IsError)
	})

	t.Run("testcase_filter", func(t *testing.T) {
		h := testCaseFilterHandler(c, log)
		r := mcplib.CallToolRequest{}
		r.Params.Arguments = map[string]any{"query": map[string]any{"pk": float64(created.ID)}}
		res, err := h(context.Background(), r)
		require.NoError(t, err)
		assert.False(t, res.IsError)
	})
}

func TestE2E_MCP_AddLink_KnownBroken(t *testing.T) {
	c := e2eClient(t)
	log := zerolog.Nop()

	h := testCaseAddLinkHandler(c, log)
	req := mcplib.CallToolRequest{}
	req.Params.Arguments = map[string]any{"link": map[string]any{"case_id": 1, "url": "http://x"}}
	res, err := h(context.Background(), req)
	require.NoError(t, err)
	assert.True(t, res.IsError, "expected IsError: TestCase.add_link is not available on this Kiwi TCMS version")
}

func firstTextContent(t *testing.T, res *mcplib.CallToolResult) string {
	t.Helper()
	require.NotEmpty(t, res.Content)
	tc, ok := res.Content[0].(mcplib.TextContent)
	require.True(t, ok, "expected TextContent, got %T", res.Content[0])
	return tc.Text
}
