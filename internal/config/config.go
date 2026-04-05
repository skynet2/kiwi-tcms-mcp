package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	URL      string            `mapstructure:"url"`
	Username string            `mapstructure:"username"`
	Password string            `mapstructure:"password"`
	Headers  map[string]string `mapstructure:"headers"`
	Timeout  time.Duration     `mapstructure:"timeout"`
	LogLevel string            `mapstructure:"log_level"`
}

func Load(path string) (Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetEnvPrefix("KIWI")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("timeout", "30s")
	v.SetDefault("log_level", "info")

	for _, k := range []string{"url", "username", "password", "timeout", "log_level"} {
		_ = v.BindEnv(k)
	}

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) && !os.IsNotExist(err) {
			return Config{}, fmt.Errorf("config: read: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("config: unmarshal: %w", err)
	}

	if raw := os.Getenv("KIWI_HEADERS"); raw != "" {
		cfg.Headers = parseHeadersString(raw)
	}

	if cfg.URL == "" {
		return Config{}, errors.New("config: url is required")
	}
	if cfg.Username == "" {
		return Config{}, errors.New("config: username is required")
	}
	if cfg.Password == "" {
		return Config{}, errors.New("config: password is required")
	}
	return cfg, nil
}

func parseHeadersString(raw string) map[string]string {
	out := map[string]string{}
	for _, pair := range strings.Split(raw, ",") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			continue
		}
		out[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}
	return out
}
