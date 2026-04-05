package mcp

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/skynet2/kiwi-tcms-mcp/cmd/kiwi-tcms-mcp/mcp/mocks"
	"github.com/skynet2/kiwi-tcms-mcp/pkg/kiwi"
)

func TestTestCaseFilter_ToolSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockKiwiClient(ctrl)
	m.EXPECT().
		TestCaseFilter(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, q map[string]any) ([]kiwi.TestCase, error) {
			assert.Equal(t, float64(1), q["pk"])
			return []kiwi.TestCase{{ID: 1, Summary: "x"}}, nil
		})

	h := testCaseFilterHandler(m, zerolog.Nop())
	req := mcplib.CallToolRequest{}
	req.Params.Arguments = map[string]any{"query": map[string]any{"pk": 1.0}}
	res, err := h(context.Background(), req)
	require.NoError(t, err)
	assert.False(t, res.IsError)
}

func TestTestCaseCreate_ToolSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockKiwiClient(ctrl)
	m.EXPECT().
		TestCaseCreate(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, in kiwi.TestCaseCreateInput) (kiwi.TestCase, error) {
			assert.Equal(t, "new", in.Summary)
			assert.Equal(t, int64(1), in.Category)
			return kiwi.TestCase{ID: 42, Summary: "new"}, nil
		})

	h := testCaseCreateHandler(m, zerolog.Nop())
	req := mcplib.CallToolRequest{}
	req.Params.Arguments = map[string]any{
		"summary": "new", "category": 1.0, "priority": 2.0, "case_status": 3.0,
	}
	res, err := h(context.Background(), req)
	require.NoError(t, err)
	assert.False(t, res.IsError)
}

func TestTestCase_ToolFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockKiwiClient(ctrl)
	m.EXPECT().TestCaseFilter(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)

	h := testCaseFilterHandler(m, zerolog.Nop())
	req := mcplib.CallToolRequest{}
	req.Params.Arguments = map[string]any{"query": map[string]any{}}
	res, err := h(context.Background(), req)
	require.NoError(t, err)
	assert.True(t, res.IsError)
}

func TestTestCaseCreate_MissingArgs(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockKiwiClient(ctrl)
	h := testCaseCreateHandler(m, zerolog.Nop())
	req := mcplib.CallToolRequest{}
	req.Params.Arguments = map[string]any{}
	res, err := h(context.Background(), req)
	require.NoError(t, err)
	assert.True(t, res.IsError)
}

func TestTestCaseRelations_ToolSuccess(t *testing.T) {
	type tc struct {
		name    string
		args    map[string]any
		handler func(KiwiClient, zerolog.Logger) toolHandler
		expect  func(m *mocks.MockKiwiClient)
	}
	cases := []tc{
		{"add_tag", map[string]any{"case_id": 1.0, "tag": "smoke"},
			testCaseAddTagHandler,
			func(m *mocks.MockKiwiClient) {
				m.EXPECT().TestCaseAddTag(gomock.Any(), int64(1), "smoke").Return(nil)
			}},
		{"remove_tag", map[string]any{"case_id": 1.0, "tag": "smoke"},
			testCaseRemoveTagHandler,
			func(m *mocks.MockKiwiClient) {
				m.EXPECT().TestCaseRemoveTag(gomock.Any(), int64(1), "smoke").Return(nil)
			}},
		{"add_component", map[string]any{"case_id": 1.0, "component": "web"},
			testCaseAddComponentHandler,
			func(m *mocks.MockKiwiClient) {
				m.EXPECT().TestCaseAddComponent(gomock.Any(), int64(1), "web").Return(nil)
			}},
		{"remove_component", map[string]any{"case_id": 1.0, "component_id": 1.0},
			testCaseRemoveComponentHandler,
			func(m *mocks.MockKiwiClient) {
				m.EXPECT().TestCaseRemoveComponent(gomock.Any(), int64(1), int64(1)).Return(nil)
			}},
		{"add_comment", map[string]any{"case_id": 1.0, "comment": "hi"},
			testCaseAddCommentHandler,
			func(m *mocks.MockKiwiClient) {
				m.EXPECT().TestCaseAddComment(gomock.Any(), int64(1), "hi").Return(nil)
			}},
		{"add_property", map[string]any{"case_id": 1.0, "name": "env", "value": "prod"},
			testCaseAddPropertyHandler,
			func(m *mocks.MockKiwiClient) {
				m.EXPECT().TestCaseAddProperty(gomock.Any(), int64(1), "env", "prod").Return(nil)
			}},
		{"remove_property", map[string]any{"query": map[string]any{"case": 1.0}},
			testCaseRemovePropertyHandler,
			func(m *mocks.MockKiwiClient) {
				m.EXPECT().TestCaseRemoveProperty(gomock.Any(), gomock.Any()).Return(nil)
			}},
		{"add_link", map[string]any{"link": map[string]any{"case_id": 1.0, "url": "http://x"}},
			testCaseAddLinkHandler,
			func(m *mocks.MockKiwiClient) {
				m.EXPECT().TestCaseAddLink(gomock.Any(), gomock.Any()).Return(map[string]any{"id": 1.0}, nil)
			}},
		{"remove_link", map[string]any{"query": map[string]any{"case": 1.0}},
			testCaseRemoveLinkHandler,
			func(m *mocks.MockKiwiClient) {
				m.EXPECT().TestCaseRemoveLink(gomock.Any(), gomock.Any()).Return(nil)
			}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			m := mocks.NewMockKiwiClient(ctrl)
			c.expect(m)
			h := c.handler(m, zerolog.Nop())
			req := mcplib.CallToolRequest{}
			req.Params.Arguments = c.args
			res, err := h(context.Background(), req)
			require.NoError(t, err)
			assert.False(t, res.IsError)
		})
	}
}
