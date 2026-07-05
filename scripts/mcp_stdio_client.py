"""
Minimal MCP JSON-RPC 2.0 stdio client.

Drives any compiled MCP server binary over its real stdin/stdout using the
actual protocol lifecycle: initialize -> notifications/initialized ->
tools/list -> tools/call. Shared by scripts/e2e_mcp_harness.py (manual,
live-account testing) and tests/e2e/ (automated, mock-server-backed pytest
suite) so the stdio JSON-RPC plumbing exists in exactly one place.
"""

import json
import subprocess
import sys
import threading
import time

# Bing responses can contain arbitrary Unicode (real search queries, non-English content).
# Reconfigure our own stdout/stderr so printing that data never raises UnicodeEncodeError,
# regardless of the host console's default codepage (notably an issue on Windows).
for _stream in (sys.stdout, sys.stderr):
    if hasattr(_stream, "reconfigure"):
        _stream.reconfigure(encoding="utf-8", errors="replace")


class McpClient:
    def __init__(self, command, env):
        self.proc = subprocess.Popen(
            command,
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            env=env,
            text=True,
            encoding="utf-8",
            errors="replace",
            bufsize=1,
        )
        self._next_id = 1
        self._stderr_lines = []
        self._stderr_thread = threading.Thread(target=self._drain_stderr, daemon=True)
        self._stderr_thread.start()

    def _drain_stderr(self):
        for line in self.proc.stderr:
            self._stderr_lines.append(line.rstrip("\n"))

    def _send(self, obj):
        line = json.dumps(obj)
        self.proc.stdin.write(line + "\n")
        self.proc.stdin.flush()

    def _read_response(self, expected_id, timeout=15):
        deadline = time.time() + timeout
        while time.time() < deadline:
            line = self.proc.stdout.readline()
            if not line:
                if self.proc.poll() is not None:
                    raise RuntimeError(
                        f"Server exited (code={self.proc.returncode}) before responding. "
                        f"stderr: {self._stderr_lines}"
                    )
                continue
            line = line.strip()
            if not line:
                continue
            try:
                msg = json.loads(line)
            except json.JSONDecodeError:
                continue
            if msg.get("id") == expected_id:
                return msg
        raise TimeoutError(f"Timed out waiting for response id={expected_id}. stderr: {self._stderr_lines}")

    def request(self, method, params=None):
        req_id = self._next_id
        self._next_id += 1
        self._send({"jsonrpc": "2.0", "id": req_id, "method": method, "params": params or {}})
        return self._read_response(req_id)

    def notify(self, method, params=None):
        self._send({"jsonrpc": "2.0", "method": method, "params": params or {}})

    def initialize(self):
        result = self.request(
            "initialize",
            {
                "protocolVersion": "2024-11-05",
                "capabilities": {},
                "clientInfo": {"name": "e2e-test-harness", "version": "0.1"},
            },
        )
        self.notify("notifications/initialized")
        return result

    def list_tools(self):
        return self.request("tools/list")

    def call_tool(self, name, arguments):
        return self.request("tools/call", {"name": name, "arguments": arguments})

    def close(self):
        try:
            self.proc.stdin.close()
        except Exception:
            pass
        try:
            self.proc.terminate()
            self.proc.wait(timeout=5)
        except Exception:
            try:
                self.proc.kill()
            except Exception:
                pass
