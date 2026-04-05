package kiwi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const jsonRPCPath = "/json-rpc/"

type Config struct {
	BaseURL      string
	Username     string
	Password     string
	ExtraHeaders map[string]string
	HTTPClient   *http.Client
}

type Client struct {
	baseURL      string
	username     string
	password     string
	extraHeaders map[string]string
	http         *http.Client
	reqID        atomic.Int64
	loginMu      sync.Mutex
	loginGen     atomic.Int64
}

func New(cfg Config) (*Client, error) {
	if cfg.BaseURL == "" {
		return nil, errors.New("kiwi: BaseURL required")
	}
	if cfg.Username == "" || cfg.Password == "" {
		return nil, errors.New("kiwi: Username and Password required")
	}
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		jar, _ := cookiejar.New(nil)
		httpClient = &http.Client{Timeout: 30 * time.Second, Jar: jar}
	}
	if httpClient.Jar == nil {
		jar, _ := cookiejar.New(nil)
		httpClient.Jar = jar
	}
	return &Client{
		baseURL:      strings.TrimRight(cfg.BaseURL, "/"),
		username:     cfg.Username,
		password:     cfg.Password,
		extraHeaders: cfg.ExtraHeaders,
		http:         httpClient,
	}, nil
}

type rpcEnvelope struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
	ID      int64  `json:"id"`
}

type rpcResult struct {
	Jsonrpc string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
	ID      int64           `json:"id"`
}

func (c *Client) Login(ctx context.Context) error {
	if err := c.doOnce(ctx, "Auth.login", []any{c.username, c.password}, nil); err != nil {
		return err
	}
	c.loginGen.Add(1)
	return nil
}

func (c *Client) Do(ctx context.Context, method string, params any, result any) error {
	gen := c.loginGen.Load()
	err := c.doOnce(ctx, method, params, result)
	if !errors.Is(err, ErrUnauthorized) {
		return err
	}
	if err := c.reloginIfStale(ctx, gen); err != nil {
		return fmt.Errorf("kiwi: relogin failed: %w", err)
	}
	return c.doOnce(ctx, method, params, result)
}

func (c *Client) reloginIfStale(ctx context.Context, observed int64) error {
	c.loginMu.Lock()
	defer c.loginMu.Unlock()
	if c.loginGen.Load() != observed {
		return nil
	}
	if err := c.doOnce(ctx, "Auth.login", []any{c.username, c.password}, nil); err != nil {
		return err
	}
	c.loginGen.Add(1)
	return nil
}

func (c *Client) doOnce(ctx context.Context, method string, params any, result any) error {
	id := c.reqID.Add(1)
	env := rpcEnvelope{Jsonrpc: "2.0", Method: method, Params: params, ID: id}
	body, err := json.Marshal(env)
	if err != nil {
		return fmt.Errorf("kiwi: marshal request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+jsonRPCPath, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("kiwi: new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	for k, v := range c.extraHeaders {
		req.Header.Set(k, v)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("kiwi: http: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("kiwi: http status %d", resp.StatusCode)
	}
	var out rpcResult
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return fmt.Errorf("kiwi: decode response: %w", err)
	}
	if out.Error != nil {
		return out.Error
	}
	if result != nil && len(out.Result) > 0 {
		if err := json.Unmarshal(out.Result, result); err != nil {
			return fmt.Errorf("kiwi: unmarshal result: %w", err)
		}
	}
	return nil
}
