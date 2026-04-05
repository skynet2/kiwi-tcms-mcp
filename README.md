# kiwi-tcms-mcp

MCP stdio server for [Kiwi TCMS](https://kiwitcms.org/) test-case management.

## Build

    make build

## Configure

Copy `configs/config.example.yaml` to `config.yaml` and fill in credentials,
or set env vars (prefix `KIWI_`):

    KIWI_URL=https://tcms.example.com
    KIWI_USERNAME=alice
    KIWI_PASSWORD=...
    KIWI_HEADERS="CF-Access-Client-Id=xxx,CF-Access-Client-Secret=yyy"

## Run

    ./bin/kiwi-tcms-mcp serve

Or register as an MCP server in your client config pointing to the binary.

## Tools (21 total)

**TestCase CRUD:** `testcase_filter`, `testcase_get`, `testcase_create`,
`testcase_update`, `testcase_remove`.

**TestCase relationships:** `testcase_add_tag`, `testcase_remove_tag`,
`testcase_add_component`, `testcase_remove_component`, `testcase_add_comment`,
`testcase_add_property`, `testcase_remove_property`, `testcase_add_link`,
`testcase_remove_link`.

**Lookups:** `category_filter`, `priority_filter`, `product_filter`,
`component_filter`, `tag_filter`, `testcase_status_filter`, `user_filter`.

## Dev

    make test      # run all tests
    make lint      # golangci-lint
    make generate  # regenerate mocks
