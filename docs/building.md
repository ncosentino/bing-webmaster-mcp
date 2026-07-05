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

`scripts/e2e_mcp_harness.py` drives a compiled binary over the real MCP stdio JSON-RPC protocol
(`initialize` -> `notifications/initialized` -> `tools/list` -> `tools/call` for every tool) --
useful for validating the actual compiled server, not just unit tests of internal classes.
Requires Python 3.

```bash
python scripts/e2e_mcp_harness.py \
  --label "Go E2E" \
  --command go/bwt-mcp-go.exe \
  --api-key YOUR_REAL_OR_DUMMY_KEY
```

With a real API key this exercises live Bing responses end-to-end. With a dummy key it still
validates the full protocol lifecycle, tool schema registration, and that the server correctly
surfaces Bing's real error responses (`{"ErrorCode":..., "Message":...}` at HTTP 400) as clean
MCP tool errors -- without needing a real account. Run it against both `go/bwt-mcp-go.exe` and the
published `bwt-mcp-csharp` binary and diff the `tools/list` output to confirm Go/C# parity.

---

## Contributing

1. Open an issue describing the bug or feature before submitting a PR
2. Run `golangci-lint run` (Go) or `dotnet build` with zero warnings (C#) before submitting
3. Keep both implementations in sync -- a feature added to Go should also be added to C#
