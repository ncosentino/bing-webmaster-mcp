"""
Automated, mock-server-backed E2E tests for every MCP tool.

For each tool this asserts, against BOTH compiled binaries, that:
  - the MCP tool call succeeds with only required arguments supplied. This is
    exactly the layer client-level unit tests cannot reach -- it is what
    would have caught a missing default value on an optional parameter,
    since unit tests call client methods directly and never go through MCP's
    own argument binding.
  - the MCP tool call succeeds with every optional argument supplied too.
  - the response is valid JSON with no "error" key.
  - the request actually reached the Bing method name the tool is documented
    to call, catching the class of bug where a tool silently hits the wrong
    endpoint (e.g. "AddSiteRole" instead of "AddSiteRoles").
"""

import json

import pytest

from tool_manifest import TOOL_CASES, assert_manifest_matches_tool_list


def test_tools_list_matches_manifest(mcp_client):
    tools = mcp_client.list_tools()["result"]["tools"]
    tool_names = [t["name"] for t in tools]
    assert len(tool_names) == 43, f"expected 43 tools, got {len(tool_names)}: {sorted(tool_names)}"
    assert_manifest_matches_tool_list(tool_names)


def _call_and_parse(mcp_client, tool, args):
    resp = mcp_client.call_tool(tool, args)
    assert "error" not in resp, f"{tool}: protocol-level error: {resp.get('error')}"

    content = resp.get("result", {}).get("content", [])
    assert content, f"{tool}: no content returned in response {resp}"

    payload = json.loads(content[0]["text"])
    assert "error" not in payload, f"{tool}: tool-level error with args {args}: {payload}"
    return payload


@pytest.mark.parametrize("case", TOOL_CASES, ids=lambda c: c.tool)
def test_tool_minimal_args(mcp_client, mock_server, case):
    """Calling with only the required arguments must succeed -- every optional
    parameter is omitted, exercising MCP's own default-value binding."""
    _call_and_parse(mcp_client, case.tool, case.minimal_args)

    if case.bing_method:
        req = mock_server.last_request_for(case.bing_method)
        assert req is not None, f"{case.tool}: expected a request to Bing method {case.bing_method!r}, saw none"
        request_path = req["path"].split("?")[0]
        assert request_path.endswith(f"/{case.bing_method}"), (
            f"{case.tool}: request hit {request_path!r}, expected it to end with /{case.bing_method}"
        )


@pytest.mark.parametrize("case", [c for c in TOOL_CASES if c.full_args is not None], ids=lambda c: c.tool)
def test_tool_full_args(mcp_client, case):
    """Calling with every optional argument explicitly supplied must also succeed."""
    _call_and_parse(mcp_client, case.tool, case.full_args)


def test_list_sites_round_trips_real_data(mcp_client):
    payload = _call_and_parse(mcp_client, "list_sites", {})
    assert payload["sites"][0]["siteUrl"] == "https://example.test/"
    assert payload["sites"][0]["isVerified"] is True
    assert payload["sites"][0]["dnsVerificationCode"] == "mock-dns-code"


def test_get_site_roles_decodes_role_enum(mcp_client):
    """Role is an int on the wire (0=Administrator, 1=ReadOnly, 2=ReadWrite) --
    both languages must decode it into the matching readable string."""
    payload = _call_and_parse(mcp_client, "get_site_roles", {"site_url": "https://example.test/"})
    assert payload["rows"][0]["role"] == "ReadOnly"


def test_get_crawl_issues_decodes_bitflags(mcp_client):
    """Issues=5 (binary 101) exercises crawl-issue bitflag decoding into a
    readable string array; both languages had a real decoding bug here once."""
    payload = _call_and_parse(mcp_client, "get_crawl_issues", {"site_url": "https://example.test/"})
    issues = payload["issues"][0]["issues"]
    assert isinstance(issues, list)
    assert len(issues) == 2, f"expected 2 decoded flags for bitmask 5, got {issues}"


def test_get_sitemap_details_handles_array_shaped_response(mcp_client):
    """Bing's real GetFeedDetails endpoint sometimes returns a single-element
    array instead of a bare object -- this previously crashed Go."""
    payload = _call_and_parse(
        mcp_client,
        "get_sitemap_details",
        {"site_url": "https://example.test/", "feed_url": "https://example.test/sitemap.xml"},
    )
    assert payload["feedUrl"] == "https://example.test/sitemap.xml"
    assert payload["sitemap"]["url"] == "https://example.test/sitemap.xml"


def test_get_keyword_reports_found_true_for_nonempty_query(mcp_client):
    payload = _call_and_parse(
        mcp_client,
        "get_keyword",
        {
            "query": "mock keyword",
            "country": "US",
            "language": "en-US",
            "start_date": "2024-01-01",
            "end_date": "2024-01-31",
        },
    )
    assert payload["found"] is True


def test_fire_and_forget_commands_ignore_null_payload(mcp_client):
    """Fire-and-forget Bing commands (AddSite, SubmitUrl, etc.) return "d":null
    unreliably even on success -- both clients must treat HTTP 2xx as success
    rather than trying to parse "d" as a meaningful value."""
    payload = _call_and_parse(mcp_client, "add_site", {"site_url": "https://example.test/"})
    assert payload["success"] is True


def test_bing_error_shape_surfaces_as_clean_tool_error(mcp_client):
    """Bing's real error shape ({"ErrorCode":N,"Message":"..."} at HTTP 400)
    must surface as a clean MCP tool-level error, not a crash or protocol
    failure."""
    from mock_bing_server import ERROR_TRIGGER

    resp = mcp_client.call_tool("get_crawl_issues", {"site_url": ERROR_TRIGGER})
    assert "error" not in resp, f"expected a clean tool-level error, got a protocol error: {resp.get('error')}"

    content = resp["result"]["content"]
    payload = json.loads(content[0]["text"])
    assert "error" in payload, f"expected a tool-level error for the Bing error trigger, got: {payload}"
