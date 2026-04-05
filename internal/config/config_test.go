package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Success(t *testing.T) {
	type tc struct {
		name   string
		yaml   string
		env    map[string]string
		assert func(t *testing.T, c Config)
	}
	cases := []tc{
		{
			name: "yaml only",
			yaml: "url: https://tcms.example.com\nusername: alice\npassword: s3cr3t\n",
			assert: func(t *testing.T, c Config) {
				assert.Equal(t, "https://tcms.example.com", c.URL)
				assert.Equal(t, "alice", c.Username)
			},
		},
		{
			name: "env overrides yaml",
			yaml: "url: https://old\nusername: u\npassword: p\n",
			env:  map[string]string{"KIWI_URL": "https://new"},
			assert: func(t *testing.T, c Config) {
				assert.Equal(t, "https://new", c.URL)
			},
		},
		{
			name: "headers from env",
			yaml: "url: https://x\nusername: u\npassword: p\n",
			env:  map[string]string{"KIWI_HEADERS": "CF-Access-Client-Id=aaa,CF-Access-Client-Secret=bbb"},
			assert: func(t *testing.T, c Config) {
				assert.Equal(t, "aaa", c.Headers["CF-Access-Client-Id"])
				assert.Equal(t, "bbb", c.Headers["CF-Access-Client-Secret"])
			},
		},
		{
			name: "no yaml file env only",
			yaml: "",
			env: map[string]string{
				"KIWI_URL": "https://env-only", "KIWI_USERNAME": "envuser", "KIWI_PASSWORD": "envpass",
			},
			assert: func(t *testing.T, c Config) {
				assert.Equal(t, "https://env-only", c.URL)
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "config.yaml")
			if c.yaml != "" {
				require.NoError(t, os.WriteFile(path, []byte(c.yaml), 0o600))
			}
			for k, v := range c.env {
				t.Setenv(k, v)
			}
			got, err := Load(path)
			require.NoError(t, err)
			c.assert(t, got)
		})
	}
}

func TestLoad_Failure(t *testing.T) {
	type tc struct {
		name string
		yaml string
	}
	cases := []tc{
		{"missing url", "username: u\npassword: p\n"},
		{"missing username", "url: https://x\npassword: p\n"},
		{"missing password", "url: https://x\nusername: u\n"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "config.yaml")
			require.NoError(t, os.WriteFile(path, []byte(c.yaml), 0o600))
			_, err := Load(path)
			require.Error(t, err)
		})
	}
}
