# kiwi-tcms-mcp

MCP stdio server for [Kiwi TCMS](https://kiwitcms.org/) test-case management.

## Install

Download a pre-built binary from [Releases](https://github.com/skynet2/kiwi-tcms-mcp/releases), or build from source:

    go install github.com/skynet2/kiwi-tcms-mcp/cmd/kiwi-tcms-mcp@latest

## Configure

Copy `configs/config.example.yaml` to `config.yaml` and fill in credentials,
or set env vars (prefix `KIWI_`):

| Variable | Description |
|----------|-------------|
| `KIWI_URL` | Kiwi TCMS instance URL |
| `KIWI_USERNAME` | Login username |
| `KIWI_PASSWORD` | Login password |
| `KIWI_HEADERS` | Optional extra headers, e.g. `CF-Access-Client-Id=xxx,CF-Access-Client-Secret=yyy` |
| `KIWI_TIMEOUT` | HTTP timeout (default `30s`) |
| `KIWI_LOG_LEVEL` | Log level: `debug`, `info`, `warn`, `error` (default `info`) |

## Usage

### Claude Code

Add to `~/.claude.json`:

```json
{
  "mcpServers": {
    "kiwi-tcms": {
      "type": "stdio",
      "command": "/path/to/kiwi-tcms-mcp",
      "args": ["serve"],
      "env": {
        "KIWI_URL": "https://tcms.example.com",
        "KIWI_USERNAME": "your-user",
        "KIWI_PASSWORD": "your-password"
      }
    }
  }
}
```

### Other MCP clients

Any MCP-compatible client can launch the binary as a stdio server:

    /path/to/kiwi-tcms-mcp serve

## Tools

### TestCase CRUD

| Tool | Description |
|------|-------------|
| `testcase_filter` | Search test cases by query (e.g. `{"summary__icontains":"login"}`) |
| `testcase_get` | Get a single test case by ID |
| `testcase_create` | Create a new test case |
| `testcase_update` | Update a test case by ID with a patch object |
| `testcase_remove` | Remove test cases matching a query |

### TestCase relationships

| Tool | Description |
|------|-------------|
| `testcase_add_tag` | Add a tag to a test case by name |
| `testcase_remove_tag` | Remove a tag from a test case |
| `testcase_add_component` | Add a component to a test case by name |
| `testcase_remove_component` | Remove a component from a test case by component ID |
| `testcase_add_comment` | Add a comment to a test case |
| `testcase_add_property` | Add a property (name/value pair) to a test case |
| `testcase_remove_property` | Remove test case properties matching a query |
| `testcase_add_link` | Add an external link (bug/requirement URL) to a test case |
| `testcase_remove_link` | Remove test case links matching a query |

### Lookups

| Tool | Description |
|------|-------------|
| `category_filter` | Filter categories |
| `priority_filter` | Filter priorities |
| `product_filter` | Filter products |
| `component_filter` | Filter components |
| `tag_filter` | Filter tags |
| `testcase_status_filter` | Filter test case statuses (e.g. CONFIRMED) |
| `user_filter` | Filter users |

## Dev

    make build     # compile binary
    make test      # run unit tests
    make lint      # golangci-lint
    make generate  # regenerate mocks

## License

[MIT](LICENSE)
