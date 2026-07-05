"""
End-to-end MCP protocol test harness for bing-webmaster-mcp.

Drives a compiled server binary (Go or C#) over real stdio using the actual
MCP JSON-RPC 2.0 lifecycle: initialize -> notifications/initialized ->
tools/list -> tools/call (for every tool). This is the manual, live-account
tool for occasional pre-release smoke testing against the real Bing API --
see tests/e2e/ for the automated, mock-server-backed pytest suite that runs
in CI on every push with zero live credentials.
"""

import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from mcp_stdio_client import McpClient  # noqa: E402 -- path setup must precede this import


# Tools that only read data -- safe to call against a real account.
READ_ONLY_TOOLS = frozenset({
    "list_sites",
    "list_sitemaps",
    "get_sitemap_details",
    "get_url_submission_quota",
    "get_crawl_issues",
    "get_crawl_stats",
    "get_url_info",
    "get_url_traffic_info",
    "get_url_links",
    "get_link_counts",
    "get_rank_and_traffic_stats",
    "get_query_stats",
    "get_page_stats",
    "get_page_query_stats",
    "get_query_page_stats",
    "get_keyword_stats",
})

# Tools that add, submit, or otherwise mutate account/site state.
# NEVER call these against a real account without explicit, informed consent.
WRITE_TOOLS = frozenset({
    "add_site",
    "verify_site",
    "submit_sitemap",
    "submit_url",
    "submit_url_batch",
    "submit_url_indexnow",
})

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

_all_tool_names = {n for n, _ in TOOL_CALLS}
assert READ_ONLY_TOOLS | WRITE_TOOLS == _all_tool_names, (
    "READ_ONLY_TOOLS/WRITE_TOOLS classification is out of sync with TOOL_CALLS -- "
    f"unclassified: {_all_tool_names - (READ_ONLY_TOOLS | WRITE_TOOLS)}, "
    f"unknown: {(READ_ONLY_TOOLS | WRITE_TOOLS) - _all_tool_names}"
)
assert not (READ_ONLY_TOOLS & WRITE_TOOLS), "a tool cannot be both read-only and a write tool"


def _substitute_placeholder(value, site_url, host):
    """Recursively replace the https://example.test/ placeholder with a real site."""
    if isinstance(value, str):
        return (
            value.replace("https://example.test/", site_url)
            .replace("example.test", host)
        )
    if isinstance(value, list):
        return [_substitute_placeholder(v, site_url, host) for v in value]
    return value


def run(label, command, env, read_only=False, site_url=None):
    print(f"\n{'=' * 70}\n{label}{'  [READ-ONLY MODE]' if read_only else ''}\n{'=' * 70}")
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

        host = None
        if site_url:
            host = site_url.split("://", 1)[-1].strip("/").split("/", 1)[0]

        calls = TOOL_CALLS
        if read_only:
            calls = [(n, a) for n, a in TOOL_CALLS if n in READ_ONLY_TOOLS]
            skipped = [n for n, _ in TOOL_CALLS if n not in READ_ONLY_TOOLS]
            print(f"Read-only mode: skipping {len(skipped)} write tools entirely: {skipped}")
            # Safety net: never allow a write tool through even if the filter above has a bug.
            assert all(n not in WRITE_TOOLS for n, _ in calls), "refusing to proceed: a write tool leaked into read-only call list"

        call_results = {}
        for name, args in calls:
            if site_url:
                args = {k: _substitute_placeholder(v, site_url, host) for k, v in args.items()}
            try:
                resp = client.call_tool(name, args)
                if "error" in resp:
                    call_results[name] = {"protocol_error": resp["error"]}
                    print(f"  [PROTOCOL ERROR] {name}: {resp['error']}")
                    continue
                content = resp.get("result", {}).get("content", [])
                text = content[0]["text"] if content else "<no content>"
                call_results[name] = {"raw_text": text}
                snippet = text if len(text) < 300 else text[:300] + "..."
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
    parser.add_argument(
        "--read-only",
        action="store_true",
        help="Only call read-only tools -- skips add_site/verify_site/submit_* entirely. "
        "Required when testing against a real account you don't want mutated.",
    )
    parser.add_argument(
        "--site-url",
        default=None,
        help="Substitute this for the https://example.test/ placeholder in every tool call.",
    )
    parser.add_argument(
        "--call",
        action="append",
        nargs=2,
        metavar=("TOOL_NAME", "JSON_ARGS"),
        help="Make exactly this tool call instead of running the full TOOL_CALLS battery. "
        "Repeatable. Example: --call add_site '{\"site_url\": \"https://example.com/\"}'. "
        "Bypasses --read-only filtering -- you are explicitly naming the call, so be sure "
        "you have authorization for whatever tool you name here.",
    )
    args = parser.parse_args()

    import os

    env = dict(os.environ)
    env["BING_WEBMASTER_API_KEY"] = args.api_key
    if args.indexnow_key:
        env["BING_INDEXNOW_KEY"] = args.indexnow_key

    if args.call:
        explicit_calls = [(name, json.loads(raw_args)) for name, raw_args in args.call]
        print(f"\n{'=' * 70}\n{args.label}  [EXPLICIT CALLS: {[n for n, _ in explicit_calls]}]\n{'=' * 70}")
        client = McpClient(args.command, env)
        results = {"calls": {}}
        try:
            init = client.initialize()
            print(f"initialize OK -- serverInfo={init.get('result', {}).get('serverInfo', {})}")
            for name, call_args in explicit_calls:
                resp = client.call_tool(name, call_args)
                if "error" in resp:
                    results["calls"][name] = {"protocol_error": resp["error"]}
                    print(f"  [PROTOCOL ERROR] {name}: {resp['error']}")
                    continue
                content = resp.get("result", {}).get("content", [])
                text = content[0]["text"] if content else "<no content>"
                results["calls"][name] = {"raw_text": text}
                print(f"  [OK] {name}: {text}")
        finally:
            client.close()
    else:
        results = run(args.label, args.command, env, read_only=args.read_only, site_url=args.site_url)

    if args.output:
        with open(args.output, "w", encoding="utf-8") as f:
            json.dump(results, f, indent=2)
        print(f"\nWrote results to {args.output}")
