"""
Shared pytest fixtures for the automated E2E MCP test suite.

Every test in this package drives a *compiled* server binary (Go or C#)
through the real MCP stdio JSON-RPC protocol, pointed at a local mock Bing
API server (mock_bing_server.py) via the BING_WEBMASTER_API_BASE_URL /
BING_INDEXNOW_API_BASE_URL override hooks. No live credentials or network
access are required or used, so this suite is safe to run in CI on every
push.

Binaries are expected to already be built (see docs/building.md); this suite
does not build them itself so it composes cleanly with a CI pipeline that
builds once (per language) and tests many times.
"""

import os
import platform
import sys
from pathlib import Path

import pytest

REPO_ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(REPO_ROOT / "scripts"))
sys.path.insert(0, str(Path(__file__).resolve().parent))

from mcp_stdio_client import McpClient  # noqa: E402
from mock_bing_server import MockBingServer  # noqa: E402


def _exe(name: str) -> str:
    return f"{name}.exe" if platform.system() == "Windows" else name


IMPLEMENTATIONS = {
    "go": REPO_ROOT / "go" / _exe("bwt-mcp-go"),
    "csharp": REPO_ROOT
    / "csharp"
    / "src"
    / "BingWebmasterMcp"
    / "bin"
    / "Release"
    / "net10.0"
    / _exe("bwt-mcp-csharp"),
}

BUILD_HINTS = {
    "go": f"cd go && go build -o {_exe('bwt-mcp-go')} .",
    "csharp": "cd csharp && dotnet build BingWebmasterMcp.slnx -c Release",
}


@pytest.fixture(scope="session")
def mock_server():
    server = MockBingServer()
    server.start()
    yield server
    server.stop()


@pytest.fixture(autouse=True)
def _clear_mock_requests(mock_server):
    """Every test starts with a clean request log so "last request for X"
    assertions are unambiguous, without needing a fresh server per test."""
    mock_server.clear_requests()
    yield


@pytest.fixture(params=sorted(IMPLEMENTATIONS.keys()))
def implementation(request):
    name = request.param
    path = IMPLEMENTATIONS[name]
    if not path.exists():
        pytest.skip(f"{name} binary not built at {path} -- build it first: {BUILD_HINTS[name]}")
    return name, path


@pytest.fixture
def mcp_client(implementation, mock_server):
    _, binary_path = implementation
    env = dict(os.environ)
    env["BING_WEBMASTER_API_KEY"] = "dummy-e2e-key"
    env["BING_INDEXNOW_KEY"] = "dummy-indexnow-key"
    env["BING_WEBMASTER_API_BASE_URL"] = mock_server.webmaster_base_url
    env["BING_INDEXNOW_API_BASE_URL"] = mock_server.indexnow_url

    client = McpClient([str(binary_path)], env)
    client.initialize()
    yield client
    client.close()
