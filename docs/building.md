---
description: Build the Bing Webmaster Tools MCP server from source -- Go and C# build commands, test commands, and Native AOT publish instructions for all platforms.
---

# Building from Source

Both implementations can be built locally from the repository. Pre-built binaries are available
on the [Releases page](https://github.com/ncosentino/bing-webmaster-mcp/releases/latest) if you
don't need to build from source.

---

## Go

**Requirements:** Go 1.26+

```bash
cd go
go mod tidy
go build -ldflags="-s -w" -trimpath -o bwt-mcp-go .
```

Run tests:

```bash
go test ./...
```

Run linter (requires [golangci-lint](https://golangci-lint.run/)):

```bash
golangci-lint run
```

---

## C# (.NET 10)

**Requirements:** .NET 10 SDK

Build for development (no AOT):

```bash
cd csharp
dotnet restore BingWebmasterMcp.slnx
dotnet build BingWebmasterMcp.slnx -c Release --no-restore
```

Run tests:

```bash
dotnet test BingWebmasterMcp.slnx -c Release --no-build
```

Publish as a Native AOT self-contained binary:

=== "Linux x64"
    ```bash
    dotnet publish src/BingWebmasterMcp/BingWebmasterMcp.csproj \
      -r linux-x64 -c Release --self-contained true
    ```

=== "macOS arm64"
    ```bash
    dotnet publish src/BingWebmasterMcp/BingWebmasterMcp.csproj \
      -r osx-arm64 -c Release --self-contained true
    ```

=== "Windows x64"
    ```bash
    dotnet publish src/BingWebmasterMcp/BingWebmasterMcp.csproj ^
      -r win-x64 -c Release --self-contained true
    ```

!!! note "Native AOT requirements"
    Native AOT compilation on Linux requires `clang` and `zlib1g-dev`. Install with:
    ```bash
    sudo apt-get install -y clang zlib1g-dev
    ```

---

## End-to-End Testing

Unit tests (`go test ./...`, `dotnet test`) mock the HTTP layer and are genuinely thorough about
request-building and response-parsing -- but they call client methods directly and never go
through the MCP protocol itself. That gap matters: real bugs have hidden specifically in MCP's
argument-binding and tool-dispatch layer, invisible to 100+ passing client-layer unit tests, and
only surfaced by driving the actual compiled binary over real MCP stdio. Two layers of E2E testing
close that gap.

### Automated suite (`tests/e2e/`, runs in CI)

A [pytest](https://pytest.org) suite that drives both compiled binaries over the real MCP stdio
JSON-RPC protocol (`initialize` -> `notifications/initialized` -> `tools/list` -> `tools/call`)
against a local, stdlib-only mock Bing API server (`tests/e2e/mock_bing_server.py`). No live
credentials, network access, or real Bing account are required -- this is what runs on every push
and pull request.

```bash
pip install pytest
python -m pytest tests/e2e -v
```

Both `go/bwt-mcp-go(.exe)` and the C# build output at
`csharp/src/BingWebmasterMcp/bin/Release/net10.0/bwt-mcp-csharp(.exe)` must already be built (see
above); the suite skips with a clear message if a binary is missing rather than building it itself.

For every one of the 43 tools, against both languages, the suite calls the tool twice: once with
only its required arguments (omitting every optional parameter) and once with every optional
argument explicitly supplied. The first case is exactly the scenario that catches a missing or
mismatched default value at the MCP layer -- unit tests structurally cannot reach this, since they
bypass MCP argument binding entirely by calling the client directly.

`tests/e2e/test_regression.py` permanently guards three real bugs found this way, each of which
100+ passing unit tests missed:

1. Go's `add_site_role` once called the wrong Bing endpoint (`AddSiteRole` instead of the real
   `AddSiteRoles`) -- unit tests asserted the request body but not the URL path.
2. C#'s `get_site_roles` / `add_site_role` were missing default values on optional bool parameters,
   which only breaks when the MCP SDK itself tries to bind a missing argument -- invisible to
   tests that call the client directly.
3. Go's `add_site_role` defaulted `is_read_only` to `false` (bool's zero value) when the argument
   was omitted, while C# defaulted it to `true` -- a silent, security-relevant cross-language
   divergence for the same tool and the same omitted argument, found by this suite and not by any
   prior test or manual round.

Each of these was confirmed to actually fail the corresponding regression test when deliberately
reintroduced, then re-fixed -- this isn't just passing trivially.

### Manual live-account harness (`scripts/e2e_mcp_harness.py`)

The lower-level MCP stdio client the automated suite is itself built on
(`scripts/mcp_stdio_client.py`) is also usable standalone for occasional manual testing against the
**real** Bing API and a real account -- useful before a release, or when validating a change the
mock server can't represent. Requires Python 3 and a real or dummy API key.

```bash
python scripts/e2e_mcp_harness.py \
  --label "Go E2E" \
  --command go/bwt-mcp-go.exe \
  --api-key YOUR_REAL_OR_DUMMY_KEY
```

With a real API key this exercises live Bing responses end-to-end. With a dummy key it still
validates the full protocol lifecycle, tool schema registration, and that the server correctly
surfaces Bing's real error responses (`{"ErrorCode":..., "Message":...}` at HTTP 400) as clean
MCP tool errors -- without needing a real account. Pass `--read-only` to restrict calls to
non-mutating tools when testing against an account you don't want changed, `--site-url` to
substitute a real site for the `https://example.test/` placeholder, or `--call TOOL_NAME
'{"json":"args"}'` to make one explicit, ad-hoc tool call. Run it against both `go/bwt-mcp-go.exe`
and the published `bwt-mcp-csharp` binary and diff the `tools/list` output to confirm Go/C# parity.

---

## Contributing

1. Open an issue describing the bug or feature before submitting a PR
2. Run `golangci-lint run` (Go) or `dotnet build` with zero warnings (C#) before submitting
3. Keep both implementations in sync -- a feature added to Go should also be added to C#
