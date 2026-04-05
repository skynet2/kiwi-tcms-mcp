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

func TestLookupTools_Success(t *testing.T) {
	type tc struct {
		name    string
		handler func(KiwiClient, zerolog.Logger) toolHandler
		expect  func(m *mocks.MockKiwiClient)
	}
	cases := []tc{
		{"category", categoryFilterHandler, func(m *mocks.MockKiwiClient) {
			m.EXPECT().CategoryFilter(gomock.Any(), gomock.Any()).Return([]kiwi.Category{{ID: 1}}, nil)
		}},
		{"priority", priorityFilterHandler, func(m *mocks.MockKiwiClient) {
			m.EXPECT().PriorityFilter(gomock.Any(), gomock.Any()).Return([]kiwi.Priority{{ID: 1}}, nil)
		}},
		{"product", productFilterHandler, func(m *mocks.MockKiwiClient) {
			m.EXPECT().ProductFilter(gomock.Any(), gomock.Any()).Return([]kiwi.Product{{ID: 1}}, nil)
		}},
		{"component", componentFilterHandler, func(m *mocks.MockKiwiClient) {
			m.EXPECT().ComponentFilter(gomock.Any(), gomock.Any()).Return([]kiwi.Component{{ID: 1}}, nil)
		}},
		{"tag", tagFilterHandler, func(m *mocks.MockKiwiClient) {
			m.EXPECT().TagFilter(gomock.Any(), gomock.Any()).Return([]kiwi.Tag{{ID: 1}}, nil)
		}},
		{"testcase_status", testCaseStatusFilterHandler, func(m *mocks.MockKiwiClient) {
			m.EXPECT().TestCaseStatusFilter(gomock.Any(), gomock.Any()).Return([]kiwi.TestCaseStatus{{ID: 1}}, nil)
		}},
		{"user", userFilterHandler, func(m *mocks.MockKiwiClient) {
			m.EXPECT().UserFilter(gomock.Any(), gomock.Any()).Return([]kiwi.User{{ID: 1}}, nil)
		}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			m := mocks.NewMockKiwiClient(ctrl)
			c.expect(m)
			h := c.handler(m, zerolog.Nop())
			req := mcplib.CallToolRequest{}
			req.Params.Arguments = map[string]any{"query": map[string]any{}}
			res, err := h(context.Background(), req)
			require.NoError(t, err)
			assert.False(t, res.IsError)
		})
	}
}
