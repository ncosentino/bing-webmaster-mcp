"""
Permanent regressions for real bugs found during this project's development
that unit tests (client-layer, real HTTP fixtures) did not and structurally
could not catch, because all three live specifically in the MCP
argument-binding / tool-dispatch layer sitting on top of the client:

1. Go's add_site_role once called the wrong Bing endpoint ("AddSiteRole"
   instead of the real "AddSiteRoles") -- every client-layer unit test for it
   asserted the request body but not the URL path, so 100+ passing unit tests
   did not catch a call that would have 404'd against the real API.

2. C#'s get_site_roles / add_site_role were missing default values on
   optional bool parameters (include_all_subdomains, is_administrator,
   is_read_only). Client-layer unit tests call the client directly and never
   go through MCP's own argument binding, so a required parameter with no
   default is invisible to them -- it only breaks when the MCP SDK itself
   tries (and fails) to bind a missing argument.

3. Go's add_site_role defaulted is_read_only to false (bool's zero value)
   when the argument was omitted, while C# defaulted it to true -- a
   security-relevant, silent cross-language behavioral divergence for the
   exact same tool and the exact same omitted argument. Found by this suite,
   not by any prior unit test or manual test round.

Each test below calls the real compiled binary through the real MCP stdio
protocol with the argument omitted/wrong-in-spirit and asserts the correct,
now-fixed behavior -- so a regression on any of these three would fail here
immediately rather than requiring a live account to notice.
"""

import json


def _call_and_parse(mcp_client, tool, args):
    resp = mcp_client.call_tool(tool, args)
    assert "error" not in resp, f"{tool}: protocol-level error: {resp.get('error')}"
    content = resp["result"]["content"]
    return json.loads(content[0]["text"])


def test_add_site_role_hits_add_site_roles_not_singular(mcp_client, mock_server):
    """Regression for the wrong-endpoint bug: the real Bing method is plural
    "AddSiteRoles". A regression to the singular "AddSiteRole" must be caught
    here even though it would also 404 against the real API -- this test
    doesn't require a live account to catch it."""
    _call_and_parse(
        mcp_client,
        "add_site_role",
        {
            "site_url": "https://example.test/",
            "delegated_url": "https://example.test/",
            "user_email": "someone@example.test",
            "authentication_code": "abc123",
        },
    )

    req = mock_server.last_request_for("AddSiteRoles")
    assert req is not None, (
        "no request recorded against the real Bing method 'AddSiteRoles' -- "
        "either the tool didn't fire, or it regressed to the wrong endpoint name"
    )
    request_path = req["path"].split("?")[0]
    assert request_path.endswith("/AddSiteRoles"), f"expected endpoint 'AddSiteRoles', hit {request_path!r}"
    assert not request_path.endswith("/AddSiteRole"), (
        f"regression: request hit the singular 'AddSiteRole' endpoint ({request_path!r}), "
        "which is not a real Bing method and would 404 live"
    )


def test_optional_bool_params_bind_correctly_when_omitted(mcp_client):
    """Regression for the missing-MCP-default-value bug: calling with every
    optional bool omitted must succeed (not fail at argument binding) on both
    languages, for both tools that were affected."""
    payload = _call_and_parse(mcp_client, "get_site_roles", {"site_url": "https://example.test/"})
    assert payload["includeAllSubdomains"] is False, payload

    payload = _call_and_parse(
        mcp_client,
        "add_site_role",
        {
            "site_url": "https://example.test/",
            "delegated_url": "https://example.test/",
            "user_email": "someone@example.test",
            "authentication_code": "abc123",
        },
    )
    assert payload["isAdministrator"] is False, payload


def test_add_site_role_defaults_is_read_only_to_true_on_both_languages(mcp_client):
    """Regression for the cross-language default-value divergence: omitting
    is_read_only must default to true on both Go and C#, matching the tool's
    documented behavior ("Grant read-only access when true"). Go previously
    silently defaulted to false because bool's zero value is false and
    nothing distinguished "omitted" from "explicitly false"."""
    payload = _call_and_parse(
        mcp_client,
        "add_site_role",
        {
            "site_url": "https://example.test/",
            "delegated_url": "https://example.test/",
            "user_email": "someone@example.test",
            "authentication_code": "abc123",
        },
    )
    assert payload["isReadOnly"] is True, (
        f"is_read_only must default to true when omitted, got {payload.get('isReadOnly')!r} -- "
        "this is exactly the silent Go/C# default-value divergence found and fixed this session"
    )


def test_add_site_role_still_honors_explicit_is_read_only_false(mcp_client):
    """Explicitly passing is_read_only=false must still be honored -- the
    default-value fix must not make the parameter impossible to turn off."""
    payload = _call_and_parse(
        mcp_client,
        "add_site_role",
        {
            "site_url": "https://example.test/",
            "delegated_url": "https://example.test/",
            "user_email": "someone@example.test",
            "authentication_code": "abc123",
            "is_read_only": False,
        },
    )
    assert payload["isReadOnly"] is False, payload
