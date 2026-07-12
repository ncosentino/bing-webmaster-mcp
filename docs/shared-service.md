---
description: Run one Bing Webmaster MCP service shared by every agent session.
---

# Running One Shared Service

Streamable HTTP lets every local agent connect to one long-lived Bing
Webmaster process instead of launching one STDIO child per session.

## Prepare credentials

Place the credentials in a protected `.env` file beside the binary:

```env
BING_WEBMASTER_API_KEY=your-api-key
BING_INDEXNOW_KEY=your-optional-indexnow-key
```

## Start the service

```bash
./bwt-mcp-go-linux-amd64 \
  --transport http \
  --listen-address 127.0.0.1 \
  --port 8083
```

The C# Native AOT binary uses the same arguments.

## Lifecycle management

Use a platform service supervisor, or reuse the generic
[`manage-mcp-service.ps1`](https://github.com/ncosentino/google-psi-mcp/blob/main/scripts/manage-mcp-service.ps1)
maintained for the native NexusLabs MCP servers:

```powershell
.\manage-mcp-service.ps1 Start `
  -ServiceName bing-webmaster-mcp `
  -BinaryPath C:\path\to\bwt-mcp-go.exe `
  -Port 8083
```

The manager health-checks before starting, serializes concurrent start
attempts, records process identity, and uses a per-run authenticated loopback
shutdown before falling back to terminating the verified process.

## Configure Copilot CLI

```json
{
  "mcpServers": {
    "bing-webmaster": {
      "type": "http",
      "url": "http://127.0.0.1:8083/mcp",
      "tools": ["*"]
    }
  }
}
```

Remove `command`, `args`, and `env` from the HTTP entry because Copilot no
longer launches the process. Existing sessions retain their STDIO child until
they restart.

## Network deployment

The generic manager is deliberately intended for loopback services.
Non-loopback hosting requires a platform supervisor, TLS, authentication and
authorization on every request, trusted proxy handling, and ingress limits.
