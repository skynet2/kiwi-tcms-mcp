package cli

import (
	"context"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/skynet2/kiwi-tcms-mcp/cmd/kiwi-tcms-mcp/mcp"
	"github.com/skynet2/kiwi-tcms-mcp/internal/config"
	"github.com/skynet2/kiwi-tcms-mcp/pkg/kiwi"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the MCP stdio server",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

		cfg, err := config.Load(configPath)
		if err != nil {
			logger.Error().Err(err).Msg("config load failed")
			return err
		}
		logger = logger.Level(parseLevel(cfg.LogLevel))

		client, err := kiwi.New(kiwi.Config{
			BaseURL:      cfg.URL,
			Username:     cfg.Username,
			Password:     cfg.Password,
			ExtraHeaders: cfg.Headers,
			HTTPClient:   &http.Client{Timeout: cfg.Timeout},
		})
		if err != nil {
			return err
		}
		if err := client.Login(context.Background()); err != nil {
			logger.Error().Err(err).Msg("kiwi login failed")
			return err
		}
		logger.Info().Str("url", cfg.URL).Msg("kiwi login ok")

		srv := mcp.NewServer(client, logger, "0.1.0")
		logger.Info().Int("tools", 21).Msg("starting MCP stdio server")
		return server.ServeStdio(srv)
	},
}

func parseLevel(s string) zerolog.Level {
	lvl, err := zerolog.ParseLevel(s)
	if err != nil {
		return zerolog.InfoLevel
	}
	return lvl
}
