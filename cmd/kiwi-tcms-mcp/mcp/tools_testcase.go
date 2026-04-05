package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog"

	"github.com/skynet2/kiwi-tcms-mcp/pkg/kiwi"
)

func registerTestCaseTools(s *server.MCPServer, c KiwiClient, log zerolog.Logger) {
	s.AddTool(mcplib.NewTool("testcase_filter",
		mcplib.WithDescription("Filter Kiwi TCMS TestCases by query object (e.g. {\"pk\":1} or {\"summary__icontains\":\"login\"})."),
		mcplib.WithObject("query", mcplib.Required(), mcplib.Description("filter expression as a map")),
	), server.ToolHandlerFunc(testCaseFilterHandler(c, log)))

	s.AddTool(mcplib.NewTool("testcase_get",
		mcplib.WithDescription("Get a single TestCase by id."),
		mcplib.WithNumber("id", mcplib.Required()),
	), server.ToolHandlerFunc(testCaseGetHandler(c, log)))

	s.AddTool(mcplib.NewTool("testcase_create",
		mcplib.WithDescription("Create a TestCase."),
		mcplib.WithString("summary", mcplib.Required()),
		mcplib.WithNumber("category", mcplib.Required()),
		mcplib.WithNumber("priority", mcplib.Required()),
		mcplib.WithNumber("case_status", mcplib.Required()),
		mcplib.WithString("notes"),
		mcplib.WithString("text"),
		mcplib.WithBoolean("is_automated"),
	), server.ToolHandlerFunc(testCaseCreateHandler(c, log)))

	s.AddTool(mcplib.NewTool("testcase_update",
		mcplib.WithDescription("Update a TestCase by id with a patch object."),
		mcplib.WithNumber("id", mcplib.Required()),
		mcplib.WithObject("patch", mcplib.Required()),
	), server.ToolHandlerFunc(testCaseUpdateHandler(c, log)))

	s.AddTool(mcplib.NewTool("testcase_remove",
		mcplib.WithDescription("Remove TestCases matching query."),
		mcplib.WithObject("query", mcplib.Required()),
	), server.ToolHandlerFunc(testCaseRemoveHandler(c, log)))

	s.AddTool(mcplib.NewTool("testcase_add_tag",
		mcplib.WithDescription("Add a tag to a TestCase by name."),
		mcplib.WithNumber("case_id", mcplib.Required()),
		mcplib.WithString("tag", mcplib.Required()),
	), server.ToolHandlerFunc(testCaseAddTagHandler(c, log)))

	s.AddTool(mcplib.NewTool("testcase_remove_tag",
		mcplib.WithDescription("Remove a tag from a TestCase."),
		mcplib.WithNumber("case_id", mcplib.Required()),
		mcplib.WithString("tag", mcplib.Required()),
	), server.ToolHandlerFunc(testCaseRemoveTagHandler(c, log)))

	s.AddTool(mcplib.NewTool("testcase_add_component",
		mcplib.WithDescription("Add a component to a TestCase by name."),
		mcplib.WithNumber("case_id", mcplib.Required()),
		mcplib.WithString("component", mcplib.Required()),
	), server.ToolHandlerFunc(testCaseAddComponentHandler(c, log)))

	s.AddTool(mcplib.NewTool("testcase_remove_component",
		mcplib.WithDescription("Remove a component from a TestCase by component id."),
		mcplib.WithNumber("case_id", mcplib.Required()),
		mcplib.WithNumber("component_id", mcplib.Required()),
	), server.ToolHandlerFunc(testCaseRemoveComponentHandler(c, log)))

	s.AddTool(mcplib.NewTool("testcase_add_comment",
		mcplib.WithDescription("Add a comment to a TestCase."),
		mcplib.WithNumber("case_id", mcplib.Required()),
		mcplib.WithString("comment", mcplib.Required()),
	), server.ToolHandlerFunc(testCaseAddCommentHandler(c, log)))

	s.AddTool(mcplib.NewTool("testcase_add_property",
		mcplib.WithDescription("Add a property (name/value) to a TestCase."),
		mcplib.WithNumber("case_id", mcplib.Required()),
		mcplib.WithString("name", mcplib.Required()),
		mcplib.WithString("value", mcplib.Required()),
	), server.ToolHandlerFunc(testCaseAddPropertyHandler(c, log)))

	s.AddTool(mcplib.NewTool("testcase_remove_property",
		mcplib.WithDescription("Remove TestCase properties matching query."),
		mcplib.WithObject("query", mcplib.Required()),
	), server.ToolHandlerFunc(testCaseRemovePropertyHandler(c, log)))

	s.AddTool(mcplib.NewTool("testcase_add_link",
		mcplib.WithDescription("Add an external link (e.g. bug/requirement URL) to a TestCase."),
		mcplib.WithObject("link", mcplib.Required()),
	), server.ToolHandlerFunc(testCaseAddLinkHandler(c, log)))

	s.AddTool(mcplib.NewTool("testcase_remove_link",
		mcplib.WithDescription("Remove TestCase links matching query."),
		mcplib.WithObject("query", mcplib.Required()),
	), server.ToolHandlerFunc(testCaseRemoveLinkHandler(c, log)))
}

type toolHandler func(context.Context, mcplib.CallToolRequest) (*mcplib.CallToolResult, error)

func toolErr(log zerolog.Logger, tool string, err error) (*mcplib.CallToolResult, error) {
	log.Error().Err(err).Str("tool", tool).Msg("tool call failed")
	return mcplib.NewToolResultError(err.Error()), nil
}

func toolJSON(v any) (*mcplib.CallToolResult, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return mcplib.NewToolResultError(err.Error()), nil
	}
	return mcplib.NewToolResultText(string(b)), nil
}

func reqString(req mcplib.CallToolRequest, key string) (string, error) { return req.RequireString(key) }

func reqInt64(req mcplib.CallToolRequest, key string) (int64, error) {
	v, err := req.RequireInt(key)
	if err != nil {
		return 0, err
	}
	return int64(v), nil
}

func reqMap(req mcplib.CallToolRequest, key string) (map[string]any, error) {
	args := req.GetArguments()
	v, ok := args[key]
	if !ok {
		return nil, fmt.Errorf("required argument %q not found", key)
	}
	m, ok := v.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("argument %q is not an object", key)
	}
	return m, nil
}

func optString(req mcplib.CallToolRequest, key string) string { return req.GetString(key, "") }
func optBool(req mcplib.CallToolRequest, key string) bool     { return req.GetBool(key, false) }

func testCaseFilterHandler(c KiwiClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		q, err := reqMap(req, "query")
		if err != nil {
			return toolErr(log, "testcase_filter", err)
		}
		out, err := c.TestCaseFilter(ctx, q)
		if err != nil {
			return toolErr(log, "testcase_filter", err)
		}
		return toolJSON(out)
	}
}

func testCaseGetHandler(c KiwiClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		id, err := reqInt64(req, "id")
		if err != nil {
			return toolErr(log, "testcase_get", err)
		}
		out, err := c.TestCaseFilter(ctx, map[string]any{"pk": id})
		if err != nil {
			return toolErr(log, "testcase_get", err)
		}
		if len(out) == 0 {
			return mcplib.NewToolResultError("testcase not found"), nil
		}
		return toolJSON(out[0])
	}
}

func testCaseCreateHandler(c KiwiClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		summary, err := reqString(req, "summary")
		if err != nil {
			return toolErr(log, "testcase_create", err)
		}
		category, err := reqInt64(req, "category")
		if err != nil {
			return toolErr(log, "testcase_create", err)
		}
		priority, err := reqInt64(req, "priority")
		if err != nil {
			return toolErr(log, "testcase_create", err)
		}
		caseStatus, err := reqInt64(req, "case_status")
		if err != nil {
			return toolErr(log, "testcase_create", err)
		}
		in := kiwi.TestCaseCreateInput{
			Summary:     summary,
			Category:    category,
			Priority:    priority,
			CaseStatus:  caseStatus,
			Notes:       optString(req, "notes"),
			Text:        optString(req, "text"),
			IsAutomated: optBool(req, "is_automated"),
		}
		out, err := c.TestCaseCreate(ctx, in)
		if err != nil {
			return toolErr(log, "testcase_create", err)
		}
		return toolJSON(out)
	}
}

func testCaseUpdateHandler(c KiwiClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		id, err := reqInt64(req, "id")
		if err != nil {
			return toolErr(log, "testcase_update", err)
		}
		patch, err := reqMap(req, "patch")
		if err != nil {
			return toolErr(log, "testcase_update", err)
		}
		out, err := c.TestCaseUpdate(ctx, id, patch)
		if err != nil {
			return toolErr(log, "testcase_update", err)
		}
		return toolJSON(out)
	}
}

func testCaseRemoveHandler(c KiwiClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		q, err := reqMap(req, "query")
		if err != nil {
			return toolErr(log, "testcase_remove", err)
		}
		if err := c.TestCaseRemove(ctx, q); err != nil {
			return toolErr(log, "testcase_remove", err)
		}
		return mcplib.NewToolResultText("ok"), nil
	}
}

func testCaseAddTagHandler(c KiwiClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		caseID, err := reqInt64(req, "case_id")
		if err != nil {
			return toolErr(log, "testcase_add_tag", err)
		}
		tag, err := reqString(req, "tag")
		if err != nil {
			return toolErr(log, "testcase_add_tag", err)
		}
		if err := c.TestCaseAddTag(ctx, caseID, tag); err != nil {
			return toolErr(log, "testcase_add_tag", err)
		}
		return mcplib.NewToolResultText("ok"), nil
	}
}

func testCaseRemoveTagHandler(c KiwiClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		caseID, err := reqInt64(req, "case_id")
		if err != nil {
			return toolErr(log, "testcase_remove_tag", err)
		}
		tag, err := reqString(req, "tag")
		if err != nil {
			return toolErr(log, "testcase_remove_tag", err)
		}
		if err := c.TestCaseRemoveTag(ctx, caseID, tag); err != nil {
			return toolErr(log, "testcase_remove_tag", err)
		}
		return mcplib.NewToolResultText("ok"), nil
	}
}

func testCaseAddComponentHandler(c KiwiClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		caseID, err := reqInt64(req, "case_id")
		if err != nil {
			return toolErr(log, "testcase_add_component", err)
		}
		comp, err := reqString(req, "component")
		if err != nil {
			return toolErr(log, "testcase_add_component", err)
		}
		if err := c.TestCaseAddComponent(ctx, caseID, comp); err != nil {
			return toolErr(log, "testcase_add_component", err)
		}
		return mcplib.NewToolResultText("ok"), nil
	}
}

func testCaseRemoveComponentHandler(c KiwiClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		caseID, err := reqInt64(req, "case_id")
		if err != nil {
			return toolErr(log, "testcase_remove_component", err)
		}
		compID, err := reqInt64(req, "component_id")
		if err != nil {
			return toolErr(log, "testcase_remove_component", err)
		}
		if err := c.TestCaseRemoveComponent(ctx, caseID, compID); err != nil {
			return toolErr(log, "testcase_remove_component", err)
		}
		return mcplib.NewToolResultText("ok"), nil
	}
}

func testCaseAddCommentHandler(c KiwiClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		caseID, err := reqInt64(req, "case_id")
		if err != nil {
			return toolErr(log, "testcase_add_comment", err)
		}
		comment, err := reqString(req, "comment")
		if err != nil {
			return toolErr(log, "testcase_add_comment", err)
		}
		if err := c.TestCaseAddComment(ctx, caseID, comment); err != nil {
			return toolErr(log, "testcase_add_comment", err)
		}
		return mcplib.NewToolResultText("ok"), nil
	}
}

func testCaseAddPropertyHandler(c KiwiClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		caseID, err := reqInt64(req, "case_id")
		if err != nil {
			return toolErr(log, "testcase_add_property", err)
		}
		name, err := reqString(req, "name")
		if err != nil {
			return toolErr(log, "testcase_add_property", err)
		}
		value, err := reqString(req, "value")
		if err != nil {
			return toolErr(log, "testcase_add_property", err)
		}
		if err := c.TestCaseAddProperty(ctx, caseID, name, value); err != nil {
			return toolErr(log, "testcase_add_property", err)
		}
		return mcplib.NewToolResultText("ok"), nil
	}
}

func testCaseRemovePropertyHandler(c KiwiClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		q, err := reqMap(req, "query")
		if err != nil {
			return toolErr(log, "testcase_remove_property", err)
		}
		if err := c.TestCaseRemoveProperty(ctx, q); err != nil {
			return toolErr(log, "testcase_remove_property", err)
		}
		return mcplib.NewToolResultText("ok"), nil
	}
}

func testCaseAddLinkHandler(c KiwiClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		link, err := reqMap(req, "link")
		if err != nil {
			return toolErr(log, "testcase_add_link", err)
		}
		out, err := c.TestCaseAddLink(ctx, link)
		if err != nil {
			return toolErr(log, "testcase_add_link", err)
		}
		return toolJSON(out)
	}
}

func testCaseRemoveLinkHandler(c KiwiClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		q, err := reqMap(req, "query")
		if err != nil {
			return toolErr(log, "testcase_remove_link", err)
		}
		if err := c.TestCaseRemoveLink(ctx, q); err != nil {
			return toolErr(log, "testcase_remove_link", err)
		}
		return mcplib.NewToolResultText("ok"), nil
	}
}
