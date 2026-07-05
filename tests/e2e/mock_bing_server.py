"""
Stdlib-only mock Bing Webmaster Tools API + IndexNow server.

Simulates the real ssl.bing.com/webmaster/api.svc/json/{MethodName} wire
protocol (the {"d": ...} envelope, PascalCase fields) and the IndexNow POST
endpoint well enough for the *compiled* MCP server binaries (Go and C#) to be
driven through the real MCP stdio JSON-RPC protocol without touching the live
Bing API. This exists specifically to close a gap unit tests cannot: verifying
the MCP tool-argument-binding and dispatch layer sits correctly on top of the
(separately, thoroughly unit-tested) HTTP client layer.

This is not trying to be a byte-perfect Bing simulator -- wire-format-to-model
correctness is already covered by client-layer unit tests with real fixtures.
This mock only needs to: (a) return a parseable success shape per method so
the full MCP -> client -> HTTP -> response -> MCP pipeline round-trips end to
end, (b) record exactly what request each tool sent so tests can assert the
right endpoint/verb/params were used, and (c) optionally return Bing's real
error shape to verify error surfacing end to end.
"""

from __future__ import annotations

import json
import threading
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from urllib.parse import parse_qs, urlparse

# Any query value or request body containing this exact token makes the mock
# respond with Bing's real error envelope instead of a canned success payload.
ERROR_TRIGGER = "TRIGGER_BING_ERROR"

_MOCK_DATE = "/Date(1700000000000-0800)/"

# Canned "d" payloads per Bing Webmaster JSON API method name, in the exact
# wire shape (PascalCase field names, matching the [JsonPropertyName] contract
# both clients deserialize against) keyed by the literal method name that
# appears as the last path segment of the request, e.g.
# https://ssl.bing.com/webmaster/api.svc/json/GetUserSites -> "GetUserSites".
WEBMASTER_RESPONSES: dict[str, object] = {
    "GetUserSites": [
        {
            "Url": "https://example.test/",
            "IsVerified": True,
            "DnsVerificationCode": "mock-dns-code",
            "AuthenticationCode": "mock-auth-code",
        }
    ],
    "AddSite": None,
    "VerifySite": True,
    "RemoveSite": None,
    "GetSiteRoles": [
        {
            "Date": _MOCK_DATE,
            "DelegatedCode": "mock-delegated-code",
            "DelegatorEmail": "owner@example.test",
            "DelegatedCodeOwnerEmail": "owner@example.test",
            "Email": "delegate@example.test",
            "Expired": False,
            "Role": 1,  # int enum on the wire: 0=Administrator, 1=ReadOnly, 2=ReadWrite
            "Site": "https://example.test/",
            "VerificationSite": "https://example.test/",
        }
    ],
    "AddSiteRoles": None,
    "RemoveSiteRole": None,
    "GetFeeds": [
        {
            "Url": "https://example.test/sitemap.xml",
            "Type": "Sitemap",
            "Compressed": False,
            "FileSize": 1024,
            "LastCrawled": _MOCK_DATE,
            "Submitted": _MOCK_DATE,
            "Status": "Ok",
            "UrlCount": 10,
        }
    ],
    # Real Bing returns this endpoint's payload as a single-element array in
    # some circumstances -- used deliberately here since that exact shape
    # caused a real parsing bug in Go earlier this session (only C# had
    # defensive handling until it was fixed).
    "GetFeedDetails": [
        {
            "Url": "https://example.test/sitemap.xml",
            "Type": "Sitemap",
            "Compressed": False,
            "FileSize": 1024,
            "LastCrawled": _MOCK_DATE,
            "Submitted": _MOCK_DATE,
            "Status": "Ok",
            "UrlCount": 10,
        }
    ],
    "SubmitFeed": None,
    "RemoveFeed": None,
    "SubmitUrl": None,
    "SubmitUrlBatch": None,
    "GetUrlSubmissionQuota": {"DailyQuota": 100, "MonthlyQuota": 3000},
    "SubmitContent": None,
    "GetContentSubmissionQuota": {"DailyQuota": 50, "MonthlyQuota": 1500},
    # Issues=5 (binary 101) exercises crawl-issue bitflag decoding (bit 0 + bit 2).
    "GetCrawlIssues": [
        {"Url": "https://example.test/page", "HttpCode": 200, "Issues": 5, "InLinks": 3}
    ],
    "GetCrawlStats": [
        {
            "Date": _MOCK_DATE,
            "CrawledPages": 42,
            "CrawlErrors": 1,
            "InIndex": 40,
            "InLinks": 100,
            "Code2xx": 38,
            "Code301": 1,
            "Code302": 1,
            "Code4xx": 2,
            "Code5xx": 0,
            "AllOtherCodes": 0,
            "BlockedByRobotsTxt": 0,
            "ContainsMalware": 0,
        }
    ],
    "GetUrlInfo": {
        "Url": "https://example.test/page",
        "IsPage": True,
        "HttpStatus": 200,
        "DocumentSize": 2048,
        "AnchorCount": 5,
        "DiscoveryDate": _MOCK_DATE,
        "LastCrawledDate": _MOCK_DATE,
        "TotalChildUrlCount": 0,
    },
    "GetUrlTrafficInfo": {
        "Url": "https://example.test/page",
        "IsPage": True,
        "Clicks": 12,
        "Impressions": 340,
    },
    "GetUrlLinks": {
        "Details": [{"AnchorText": "mock anchor", "Url": "https://example.test/linker"}],
        "TotalPages": 1,
    },
    "GetLinkCounts": {
        "Links": [{"Count": 7, "Url": "https://example.test/page"}],
        "TotalPages": 1,
    },
    "GetChildrenUrlInfo": [
        {
            "Url": "https://example.test/child",
            "IsPage": True,
            "HttpStatus": 200,
            "DocumentSize": 512,
            "AnchorCount": 1,
            "DiscoveryDate": _MOCK_DATE,
            "LastCrawledDate": _MOCK_DATE,
            "TotalChildUrlCount": 0,
        }
    ],
    "GetChildrenUrlTrafficInfo": [
        {"Url": "https://example.test/child", "IsPage": True, "Clicks": 3, "Impressions": 90}
    ],
    "GetBlockedUrls": [
        {"Date": _MOCK_DATE, "EntityType": 0, "RequestType": 0, "Url": "https://example.test/blocked"}
    ],
    "AddBlockedUrl": None,
    "RemoveBlockedUrl": None,
    "FetchUrl": None,
    "GetFetchedUrls": [
        {"Date": _MOCK_DATE, "Expired": False, "Fetched": True, "Url": "https://example.test/fetch-me"}
    ],
    "GetFetchedUrlDetails": {
        "Date": _MOCK_DATE,
        "Document": "mock-document-body",
        "Headers": "Content-Type: text/html",
        "Status": "200",
        "Url": "https://example.test/fetch-me",
    },
    "GetSiteMoves": [
        {
            "Date": _MOCK_DATE,
            "MoveScope": 0,
            "MoveType": 0,
            "SourceUrl": "https://example.test/",
            "TargetUrl": "https://example-new.test/",
        }
    ],
    "SubmitSiteMove": None,
    "GetRankAndTrafficStats": [{"Date": _MOCK_DATE, "Clicks": 20, "Impressions": 500}],
    # Bing reuses one wire shape (Query/Date/Clicks/Impressions/...) across four
    # distinct endpoints; for the two "page" variants the "Query" field holds a
    # page URL rather than search-query text, matching real Bing behavior.
    "GetQueryStats": [
        {
            "Query": "mock query",
            "Date": _MOCK_DATE,
            "Clicks": 5,
            "Impressions": 120,
            "AvgClickPosition": 2,
            "AvgImpressionPosition": 4,
        }
    ],
    "GetPageStats": [
        {
            "Query": "https://example.test/page",
            "Date": _MOCK_DATE,
            "Clicks": 8,
            "Impressions": 200,
            "AvgClickPosition": 1,
            "AvgImpressionPosition": 3,
        }
    ],
    "GetPageQueryStats": [
        {
            "Query": "mock query",
            "Date": _MOCK_DATE,
            "Clicks": 4,
            "Impressions": 88,
            "AvgClickPosition": 2,
            "AvgImpressionPosition": 3,
        }
    ],
    "GetQueryPageStats": [
        {
            "Query": "https://example.test/page",
            "Date": _MOCK_DATE,
            "Clicks": 6,
            "Impressions": 150,
            "AvgClickPosition": 1,
            "AvgImpressionPosition": 2,
        }
    ],
    "GetQueryPageDetailStats": [{"Date": _MOCK_DATE, "Clicks": 3, "Impressions": 60, "Position": 2}],
    "GetQueryTrafficStats": [{"Date": _MOCK_DATE, "Clicks": 9, "Impressions": 210}],
    "GetKeywordStats": [
        {"Query": "mock keyword", "Date": _MOCK_DATE, "Impressions": 1000, "BroadImpressions": 1500}
    ],
    "GetKeyword": {"Query": "mock keyword", "BroadImpressions": 2000, "Impressions": 1300},
    "GetRelatedKeywords": [
        {"Query": "related mock keyword", "BroadImpressions": 500, "Impressions": 300}
    ],
    # Phase 3 -- URL normalization (query parameters).
    "GetQueryParameters": [
        {"Date": _MOCK_DATE, "IsEnabled": True, "Parameter": "utm_campaign", "Source": 0}
    ],
    "AddQueryParameter": None,
    "RemoveQueryParameter": None,
    "EnableDisableQueryParameter": None,
    # Phase 3 -- geo-targeting (country/region settings). Type=2 exercises the
    # Domain enum value (0=Page, 1=Directory, 2=Domain, 3=Subdomain).
    "GetCountryRegionSettings": [
        {"Date": _MOCK_DATE, "TwoLetterIsoCountryCode": "us", "Type": 2, "Url": "https://example.test/"}
    ],
    "AddCountryRegionSettings": None,
    "RemoveCountryRegionSettings": None,
    # Phase 3 -- connected pages. Real wire shape has 17 fields (ConnectedSite);
    # only a useful subset is modeled/exposed, extra fields included here to
    # confirm both clients tolerate an unmapped-field superset without crashing.
    "GetConnectedPages": [
        {
            "Url": "https://connected.example.test/",
            "IsVerified": True,
            "RequestedMasterSite": "https://example.test/",
            "ActualMasterSite": "https://example.test/",
            "HttpStatusCode": 200,
            "Market": "en-US",
            "IsBlocked": False,
            "LastSuccessfullyVerified": _MOCK_DATE,
            "AppId": "mock-app-id",
            "AppName": "mock-app-name",
            "ConsecutiveFailedAttempts": 0,
            "FailureCode": 0,
            "FirstSuccessfullyVerified": _MOCK_DATE,
            "UpdateTime": _MOCK_DATE,
        }
    ],
    "AddConnectedPage": None,
    # Phase 3 -- page preview blocks. BlockReason=4 ("Other") is the one value
    # confirmed via a real recorded cassette; Action/RefreshReason/Reason/
    # SiteUrl/UserId are real wire fields intentionally not exposed by either
    # client (lower-confidence or redundant), included here to confirm both
    # tolerate the extra fields without crashing.
    "GetActivePagePreviewBlocks": [
        {
            "__type": "PagePreview:#Microsoft.Bing.Webmaster.Shared.DataContracts",
            "Action": 0,
            "BlockReason": 4,
            "Reason": "4",
            "RefreshReason": 0,
            "SiteUrl": "https://example.test/",
            "SubmitDate": _MOCK_DATE,
            "Url": "https://example.test/blocked-preview",
            "UserId": "mock-user-id",
        }
    ],
    "AddPagePreviewBlock": None,
    "RemovePagePreviewBlock": None,
}

# Bing method names that are fire-and-forget commands: real Bing returns
# "d":null unreliably even on success, so both clients treat HTTP 2xx as the
# only success signal for these. Recorded here so tests can assert this
# subset never leaks a meaningful "d" value a client might start depending on.
FIRE_AND_FORGET_METHODS = frozenset({
    "AddSite", "RemoveSite", "AddSiteRoles", "RemoveSiteRole", "SubmitFeed",
    "RemoveFeed", "SubmitUrl", "SubmitUrlBatch", "SubmitContent", "AddBlockedUrl",
    "RemoveBlockedUrl", "FetchUrl", "SubmitSiteMove",
    "AddQueryParameter", "RemoveQueryParameter", "EnableDisableQueryParameter",
    "AddCountryRegionSettings", "RemoveCountryRegionSettings", "AddConnectedPage",
    "AddPagePreviewBlock", "RemovePagePreviewBlock",
})


class MockBingServer:
    """A local HTTP server simulating the Bing Webmaster JSON API + IndexNow."""

    def __init__(self) -> None:
        self._lock = threading.Lock()
        self.requests: list[dict] = []
        self._httpd = ThreadingHTTPServer(("127.0.0.1", 0), _make_handler(self))
        self._thread = threading.Thread(target=self._httpd.serve_forever, daemon=True)

    @property
    def port(self) -> int:
        return self._httpd.server_address[1]

    @property
    def webmaster_base_url(self) -> str:
        return f"http://127.0.0.1:{self.port}/webmaster/api.svc/json"

    @property
    def indexnow_url(self) -> str:
        return f"http://127.0.0.1:{self.port}/indexnow"

    def start(self) -> None:
        self._thread.start()

    def stop(self) -> None:
        self._httpd.shutdown()
        self._httpd.server_close()

    def record(self, entry: dict) -> None:
        with self._lock:
            self.requests.append(entry)

    def last_request_for(self, method_name: str) -> dict | None:
        with self._lock:
            for entry in reversed(self.requests):
                if entry["method_name"] == method_name:
                    return entry
        return None

    def requests_for(self, method_name: str) -> list[dict]:
        with self._lock:
            return [entry for entry in self.requests if entry["method_name"] == method_name]

    def clear_requests(self) -> None:
        with self._lock:
            self.requests.clear()


def _make_handler(server: "MockBingServer"):
    class Handler(BaseHTTPRequestHandler):
        protocol_version = "HTTP/1.1"

        def log_message(self, format, *args):  # noqa: A002 -- stdlib override signature
            pass  # keep test output focused on assertions, not raw HTTP access logs

        def _dispatch(self) -> None:
            parsed = urlparse(self.path)
            query = {k: v[0] for k, v in parse_qs(parsed.query).items()}

            length = int(self.headers.get("Content-Length", 0) or 0)
            raw_body = self.rfile.read(length) if length else b""

            if parsed.path.rstrip("/").endswith("/indexnow"):
                self._handle_indexnow(query, raw_body)
                return

            method_name = parsed.path.rstrip("/").rsplit("/", 1)[-1]
            self._handle_webmaster(method_name, query, raw_body)

        def _handle_indexnow(self, query: dict, raw_body: bytes) -> None:
            server.record(
                {
                    "method_name": "IndexNow",
                    "http_method": self.command,
                    "path": self.path,
                    "query": query,
                    "body": _safe_json(raw_body),
                }
            )
            self.send_response(200)
            self.send_header("Content-Length", "0")
            self.end_headers()

        def _handle_webmaster(self, method_name: str, query: dict, raw_body: bytes) -> None:
            body = _safe_json(raw_body)
            server.record(
                {
                    "method_name": method_name,
                    "http_method": self.command,
                    "path": self.path,
                    "query": query,
                    "body": body,
                }
            )

            haystack = json.dumps(body) + json.dumps(query)
            if ERROR_TRIGGER in haystack:
                self._send_json(
                    400, {"ErrorCode": 1, "Message": "ERROR!!! Simulated failure for testing."}
                )
                return

            if method_name not in WEBMASTER_RESPONSES:
                self._send_json(404, {"Message": f"Method '{method_name}' not found."})
                return

            self._send_json(200, {"d": WEBMASTER_RESPONSES[method_name]})

        def _send_json(self, status: int, payload: object) -> None:
            body = json.dumps(payload).encode("utf-8")
            self.send_response(status)
            self.send_header("Content-Type", "application/json")
            self.send_header("Content-Length", str(len(body)))
            self.end_headers()
            self.wfile.write(body)

        def do_GET(self) -> None:
            self._dispatch()

        def do_POST(self) -> None:
            self._dispatch()

    return Handler


def _safe_json(raw: bytes) -> object:
    if not raw:
        return None
    try:
        return json.loads(raw.decode("utf-8"))
    except (json.JSONDecodeError, UnicodeDecodeError):
        return None
