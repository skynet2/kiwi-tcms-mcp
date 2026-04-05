package mcp

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog"
)

func NewServer(client KiwiClient, log zerolog.Logger, version string) *server.MCPServer {
	s := server.NewMCPServer("kiwi-tcms-mcp", version)
	registerTestCaseTools(s, client, log)
	registerLookupTools(s, client, log)
	return s
}
