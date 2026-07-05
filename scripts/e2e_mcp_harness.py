"""
End-to-end MCP protocol test harness for bing-webmaster-mcp.

Drives a compiled server binary (Go or C#) over real stdio using the actual
MCP JSON-RPC 2.0 lifecycle: initialize -> notifications/initialized ->
tools/list -> tools/call (for every tool). Used to validate the real,
compiled binaries end-to-end -- not just unit tests of internal classes.
"""

import json
import subprocess
import sys
import threading
import time


class McpClient:
    def __init__(self, command, env):
        self.proc = subprocess.Popen(
            command,
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            env=env,
            text=True,
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


TOOL_CALLS = [
    ("list_sites", {}),
    ("add_site", {"site_url": "https://example.test/"}),
    ("verify_site", {"site_url": "https://example.test/"}),
    ("list_sitemaps", {"site_url": "https://example.test/"}),
    ("get_sitemap_details", {"site_url": "https://example.test/", "feed_url": "https://example.test/sitemap.xml"}),
    ("submit_sitemap", {"site_url": "https://example.test/", "feed_url": "https://example.test/sitemap.xml"}),
    ("submit_url", {"site_url": "https://example.test/", "url": "https://example.test/page"}),
    ("submit_url_batch", {"site_url": "https://example.test/", "url_list": ["https://example.test/a"]}),
    ("submit_url_indexnow", {"host": "example.test", "url_list": ["https://example.test/a"], "key": "deadbeef"}),
    ("get_url_submission_quota", {"site_url": "https://example.test/"}),
    ("get_crawl_issues", {"site_url": "https://example.test/"}),
    ("get_crawl_stats", {"site_url": "https://example.test/"}),
    ("get_url_info", {"site_url": "https://example.test/", "url": "https://example.test/page"}),
    ("get_url_traffic_info", {"site_url": "https://example.test/", "url": "https://example.test/page"}),
    ("get_url_links", {"site_url": "https://example.test/", "link": "https://example.test/page", "page": 0}),
    ("get_link_counts", {"site_url": "https://example.test/", "page": 0}),
    ("get_rank_and_traffic_stats", {"site_url": "https://example.test/"}),
    ("get_query_stats", {"site_url": "https://example.test/"}),
    ("get_page_stats", {"site_url": "https://example.test/"}),
    ("get_page_query_stats", {"site_url": "https://example.test/", "page": "https://example.test/page"}),
    ("get_query_page_stats", {"site_url": "https://example.test/", "query": "bing webmaster api"}),
    ("get_keyword_stats", {"query": "bing webmaster api", "country": "us", "language": "en-US"}),
]


def run(label, command, env):
    print(f"\n{'=' * 70}\n{label}\n{'=' * 70}")
    client = McpClient(command, env)
    results = {}
    try:
        init = client.initialize()
        server_info = init.get("result", {}).get("serverInfo", {})
        print(f"initialize OK -- serverInfo={server_info}")

        tools_resp = client.list_tools()
        tools = tools_resp.get("result", {}).get("tools", [])
        tool_names = sorted(t["name"] for t in tools)
        print(f"tools/list OK -- {len(tools)} tools: {tool_names}")
        results["tool_count"] = len(tools)
        results["tool_names"] = tool_names

        call_results = {}
        for name, args in TOOL_CALLS:
            try:
                resp = client.call_tool(name, args)
                if "error" in resp:
                    call_results[name] = {"protocol_error": resp["error"]}
                    print(f"  [PROTOCOL ERROR] {name}: {resp['error']}")
                    continue
                content = resp.get("result", {}).get("content", [])
                text = content[0]["text"] if content else "<no content>"
                call_results[name] = {"raw_text": text}
                snippet = text if len(text) < 200 else text[:200] + "..."
                print(f"  [OK] {name}: {snippet}")
            except Exception as e:
                call_results[name] = {"exception": str(e)}
                print(f"  [EXCEPTION] {name}: {e}")
        results["calls"] = call_results
    finally:
        client.close()
    return results


if __name__ == "__main__":
    import argparse

    parser = argparse.ArgumentParser()
    parser.add_argument("--label", required=True)
    parser.add_argument("--command", required=True, nargs="+")
    parser.add_argument("--api-key", default="test-invalid-key-e2e")
    parser.add_argument("--indexnow-key", default="")
    parser.add_argument("--output", default=None)
    args = parser.parse_args()

    import os

    env = dict(os.environ)
    env["BING_WEBMASTER_API_KEY"] = args.api_key
    if args.indexnow_key:
        env["BING_INDEXNOW_KEY"] = args.indexnow_key

    results = run(args.label, args.command, env)

    if args.output:
        with open(args.output, "w", encoding="utf-8") as f:
            json.dump(results, f, indent=2)
        print(f"\nWrote results to {args.output}")
