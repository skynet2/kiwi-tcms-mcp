package mcp

import (
	"context"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog"
)

func registerLookupTools(s *server.MCPServer, c KiwiClient, log zerolog.Logger) {
	s.AddTool(mcplib.NewTool("category_filter",
		mcplib.WithDescription("Filter Kiwi TCMS Categories."),
		mcplib.WithObject("query", mcplib.Required()),
	), server.ToolHandlerFunc(categoryFilterHandler(c, log)))

	s.AddTool(mcplib.NewTool("priority_filter",
		mcplib.WithDescription("Filter Priorities."),
		mcplib.WithObject("query", mcplib.Required()),
	), server.ToolHandlerFunc(priorityFilterHandler(c, log)))

	s.AddTool(mcplib.NewTool("product_filter",
		mcplib.WithDescription("Filter Products."),
		mcplib.WithObject("query", mcplib.Required()),
	), server.ToolHandlerFunc(productFilterHandler(c, log)))

	s.AddTool(mcplib.NewTool("component_filter",
		mcplib.WithDescription("Filter Components."),
		mcplib.WithObject("query", mcplib.Required()),
	), server.ToolHandlerFunc(componentFilterHandler(c, log)))

	s.AddTool(mcplib.NewTool("tag_filter",
		mcplib.WithDescription("Filter Tags."),
		mcplib.WithObject("query", mcplib.Required()),
	), server.ToolHandlerFunc(tagFilterHandler(c, log)))

	s.AddTool(mcplib.NewTool("testcase_status_filter",
		mcplib.WithDescription("Filter TestCaseStatus values (e.g. CONFIRMED)."),
		mcplib.WithObject("query", mcplib.Required()),
	), server.ToolHandlerFunc(testCaseStatusFilterHandler(c, log)))

	s.AddTool(mcplib.NewTool("user_filter",
		mcplib.WithDescription("Filter Users."),
		mcplib.WithObject("query", mcplib.Required()),
	), server.ToolHandlerFunc(userFilterHandler(c, log)))
}

func categoryFilterHandler(c KiwiClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		q, err := reqMap(req, "query")
		if err != nil {
			return toolErr(log, "category_filter", err)
		}
		out, err := c.CategoryFilter(ctx, q)
		if err != nil {
			return toolErr(log, "category_filter", err)
		}
		return toolJSON(out)
	}
}

func priorityFilterHandler(c KiwiClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		q, err := reqMap(req, "query")
		if err != nil {
			return toolErr(log, "priority_filter", err)
		}
		out, err := c.PriorityFilter(ctx, q)
		if err != nil {
			return toolErr(log, "priority_filter", err)
		}
		return toolJSON(out)
	}
}

func productFilterHandler(c KiwiClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		q, err := reqMap(req, "query")
		if err != nil {
			return toolErr(log, "product_filter", err)
		}
		out, err := c.ProductFilter(ctx, q)
		if err != nil {
			return toolErr(log, "product_filter", err)
		}
		return toolJSON(out)
	}
}

func componentFilterHandler(c KiwiClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		q, err := reqMap(req, "query")
		if err != nil {
			return toolErr(log, "component_filter", err)
		}
		out, err := c.ComponentFilter(ctx, q)
		if err != nil {
			return toolErr(log, "component_filter", err)
		}
		return toolJSON(out)
	}
}

func tagFilterHandler(c KiwiClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		q, err := reqMap(req, "query")
		if err != nil {
			return toolErr(log, "tag_filter", err)
		}
		out, err := c.TagFilter(ctx, q)
		if err != nil {
			return toolErr(log, "tag_filter", err)
		}
		return toolJSON(out)
	}
}

func testCaseStatusFilterHandler(c KiwiClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		q, err := reqMap(req, "query")
		if err != nil {
			return toolErr(log, "testcase_status_filter", err)
		}
		out, err := c.TestCaseStatusFilter(ctx, q)
		if err != nil {
			return toolErr(log, "testcase_status_filter", err)
		}
		return toolJSON(out)
	}
}

func userFilterHandler(c KiwiClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		q, err := reqMap(req, "query")
		if err != nil {
			return toolErr(log, "user_filter", err)
		}
		out, err := c.UserFilter(ctx, q)
		if err != nil {
			return toolErr(log, "user_filter", err)
		}
		return toolJSON(out)
	}
}
