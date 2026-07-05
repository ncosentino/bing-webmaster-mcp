using System.Net;
using System.Text;
using System.Text.Json;
using BingWebmasterMcp.BingWebmaster;
using Xunit;

namespace BingWebmasterMcp.Tests;

public sealed class BingWebmasterClientTests
{
    private sealed class FakeMessageHandler(Func<HttpRequestMessage, HttpResponseMessage> handler) : HttpMessageHandler
    {
        protected override Task<HttpResponseMessage> SendAsync(
            HttpRequestMessage request,
            CancellationToken cancellationToken)
            => Task.FromResult(handler(request));
    }

    private static HttpResponseMessage JsonResponse(HttpStatusCode statusCode, string json)
        => new(statusCode)
        {
            Content = new StringContent(json, Encoding.UTF8, "application/json")
        };

    [Fact]
    public async Task GetCrawlIssues_DeserializesRecordedEmptyFixture()
    {
        var handler = new FakeMessageHandler(_ => JsonResponse(HttpStatusCode.OK, """{"d": []}"""));
        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");

        var result = await client.GetCrawlIssuesAsync("https://example.test");

        Assert.Equal("https://example.test", result.SiteUrl);
        Assert.Empty(result.Issues);
    }

    [Fact]
    public async Task GetCrawlIssues_DecodesIssueBitFlags()
    {
        const string responseJson = """{"d":[{"Url":"https://example.test/a","HttpCode":404,"Issues":20,"InLinks":3}]}""";
        var handler = new FakeMessageHandler(_ => JsonResponse(HttpStatusCode.OK, responseJson));
        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");

        var result = await client.GetCrawlIssuesAsync("https://example.test");

        Assert.Single(result.Issues);
        Assert.Equal(["Code4xx", "BlockedByRobotsTxt"], result.Issues[0].Issues);
    }

    [Fact]
    public async Task GetUrlLinks_UsesLinkParameterName_AndDeserializesRecordedFixture()
    {
        HttpRequestMessage? capturedRequest = null;
        const string responseJson = """{"d": {"__type": "LinkDetails:#Microsoft.Bing.Webmaster.Api", "Details": [], "TotalPages": 0}}""";

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            return JsonResponse(HttpStatusCode.OK, responseJson);
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://ssl.bing.com/webmaster/api.svc/json");
        var result = await client.GetUrlLinksAsync("https://example.test", "https://example.test", 0);

        Assert.NotNull(capturedRequest);
        Assert.Equal(HttpMethod.Get, capturedRequest!.Method);
        Assert.Equal(
            "https://ssl.bing.com/webmaster/api.svc/json/GetUrlLinks?apikey=test-key&siteUrl=https%3A%2F%2Fexample.test&link=https%3A%2F%2Fexample.test&page=0",
            capturedRequest.RequestUri!.AbsoluteUri);
        Assert.Empty(result.Details);
        Assert.Equal(0, result.TotalPages);
    }

    [Fact]
    public async Task SubmitUrlBatch_SerializesCamelCaseBody_AndIncludesApiKeyInQuery()
    {
        HttpRequestMessage? capturedRequest = null;
        string? body = null;

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            body = request.Content is null
                ? null
                : request.Content.ReadAsStringAsync().GetAwaiter().GetResult();
            return JsonResponse(HttpStatusCode.OK, """{"d": true}""");
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.SubmitUrlBatchAsync("https://example.test", ["https://example.test/a", "https://example.test/b"]);

        Assert.NotNull(capturedRequest);
        Assert.Equal(HttpMethod.Post, capturedRequest!.Method);
        Assert.Equal("https://example.test/api/SubmitUrlBatch?apikey=test-key", capturedRequest.RequestUri!.AbsoluteUri);
        Assert.Equal("""{"siteUrl":"https://example.test","urlList":["https://example.test/a","https://example.test/b"]}""", body);
        Assert.True(result.Success);
        Assert.Equal(2, result.SubmittedCount);
    }

    [Fact]
    public async Task SubmitUrlBatch_Throws_WhenMoreThan500UrlsProvided()
    {
        var client = new BingWebmasterClient(
            new HttpClient(new FakeMessageHandler(_ => JsonResponse(HttpStatusCode.OK, """{"d": true}"""))),
            "test-key",
            "https://example.test/api");

        var urls = Enumerable.Range(1, 501).Select(i => $"https://example.test/{i}").ToArray();

        var exception = await Assert.ThrowsAsync<ArgumentException>(() => client.SubmitUrlBatchAsync("https://example.test", urls));
        Assert.Contains("maximum of 500 URLs", exception.Message);
    }

    [Fact]
    public async Task GetKeywordStats_UsesQParameter_WithoutSiteUrl()
    {
        HttpRequestMessage? capturedRequest = null;
        const string responseJson = """{"d":[{"Query":"bing webmaster tools","Date":"/Date(1732612952000+0000)/","Impressions":42,"BroadImpressions":99}]}""";

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            return JsonResponse(HttpStatusCode.OK, responseJson);
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.GetKeywordStatsAsync("bing webmaster tools", "US", "en-US");

        Assert.NotNull(capturedRequest);
        Assert.Equal(
            "https://example.test/api/GetKeywordStats?apikey=test-key&q=bing%20webmaster%20tools&country=US&language=en-US",
            capturedRequest!.RequestUri!.AbsoluteUri);
        Assert.DoesNotContain("siteUrl", capturedRequest.RequestUri!.Query, StringComparison.Ordinal);
        Assert.Single(result.Rows);
        Assert.Equal("bing webmaster tools", result.Rows[0].Query);
    }

    [Fact]
    public async Task GetPageStats_MapsWireQueryFieldToPage()
    {
        HttpRequestMessage? capturedRequest = null;
        const string responseJson = """{"d":[{"Query":"https://example.test/page-a","Date":"/Date(1732612952000+0000)/","Clicks":7,"Impressions":11,"AvgClickPosition":2,"AvgImpressionPosition":4}]}""";

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            return JsonResponse(HttpStatusCode.OK, responseJson);
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.GetPageStatsAsync("https://example.test");

        Assert.NotNull(capturedRequest);
        Assert.Equal(
            "https://example.test/api/GetPageStats?apikey=test-key&siteUrl=https%3A%2F%2Fexample.test",
            capturedRequest!.RequestUri!.AbsoluteUri);
        Assert.Single(result.Rows);
        Assert.Equal("https://example.test/page-a", result.Rows[0].Page);
        Assert.Equal(7, result.Rows[0].Clicks);
    }

    [Fact]
    public async Task ListSites_UsesGet_AndDeserializesEnvelope()
    {
        HttpRequestMessage? capturedRequest = null;
        const string responseJson = """{"d":[{"Url":"https://example.test","IsVerified":true,"DnsVerificationCode":"dns-code","AuthenticationCode":"meta-code"}]}""";

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            return JsonResponse(HttpStatusCode.OK, responseJson);
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.ListSitesAsync();

        Assert.NotNull(capturedRequest);
        Assert.Equal(HttpMethod.Get, capturedRequest!.Method);
        Assert.Equal("https://example.test/api/GetUserSites?apikey=test-key", capturedRequest.RequestUri!.AbsoluteUri);
        Assert.Single(result.Sites);
        Assert.Equal("https://example.test", result.Sites[0].SiteUrl);
        Assert.True(result.Sites[0].IsVerified);
        Assert.Equal("dns-code", result.Sites[0].DnsVerificationCode);
        Assert.Equal("meta-code", result.Sites[0].AuthenticationCode);
    }

    [Fact]
    public async Task AddSite_TreatsNullPayloadAsSuccess_MatchingRealBingBehavior()
    {
        // Regression test: real Bing "d":null payload (observed live for both a fresh add and a
        // no-op repeat of an already-added site) previously threw a JsonException because the
        // client tried to deserialize it as a bool. AddSite has no reliable boolean signal --
        // success means the HTTP call completed without error.
        HttpRequestMessage? capturedRequest = null;
        string? body = null;

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            body = request.Content is null
                ? null
                : request.Content.ReadAsStringAsync().GetAwaiter().GetResult();
            return JsonResponse(HttpStatusCode.OK, """{"d": null}""");
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.AddSiteAsync("https://example.test");

        Assert.NotNull(capturedRequest);
        Assert.Equal(HttpMethod.Post, capturedRequest!.Method);
        Assert.Equal("https://example.test/api/AddSite?apikey=test-key", capturedRequest.RequestUri!.AbsoluteUri);
        Assert.Equal("""{"siteUrl":"https://example.test"}""", body);
        Assert.Equal("https://example.test", result.SiteUrl);
        Assert.True(result.Success);
    }

    [Fact]
    public async Task AddSite_ThrowsBingWebmasterApiException_OnNonSuccessResponse()
    {
        var client = new BingWebmasterClient(
            new HttpClient(new FakeMessageHandler(_ => JsonResponse(HttpStatusCode.BadGateway, """{"error":"upstream failed"}"""))),
            "test-key",
            "https://example.test/api");

        var exception = await Assert.ThrowsAsync<BingWebmasterApiException>(() => client.AddSiteAsync("https://example.test"));

        Assert.Equal(502, exception.StatusCode);
        Assert.Contains("HTTP 502 BadGateway", exception.Message);
    }

    [Fact]
    public async Task VerifySite_UsesPostAndCamelCaseBody_AndDeserializesEnvelope()
    {
        HttpRequestMessage? capturedRequest = null;
        string? body = null;

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            body = request.Content is null
                ? null
                : request.Content.ReadAsStringAsync().GetAwaiter().GetResult();
            return JsonResponse(HttpStatusCode.OK, """{"d": true}""");
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.VerifySiteAsync("https://example.test");

        Assert.NotNull(capturedRequest);
        Assert.Equal(HttpMethod.Post, capturedRequest!.Method);
        Assert.Equal("https://example.test/api/VerifySite?apikey=test-key", capturedRequest.RequestUri!.AbsoluteUri);
        Assert.Equal("""{"siteUrl":"https://example.test"}""", body);
        Assert.Equal("https://example.test", result.SiteUrl);
        Assert.True(result.Verified);
    }

    [Fact]
    public async Task ListSitemaps_UsesSiteUrlQuery_AndDeserializesSitemapInfo()
    {
        HttpRequestMessage? capturedRequest = null;
        const string responseJson = """{"d":[{"Url":"https://example.test/sitemap.xml","Type":"Sitemap","Compressed":false,"FileSize":2048,"LastCrawled":"/Date(1732612952000+0000)/","Submitted":"/Date(1732526552000+0000)/","Status":"Pending","UrlCount":12}]}""";

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            return JsonResponse(HttpStatusCode.OK, responseJson);
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.ListSitemapsAsync("https://example.test");

        Assert.NotNull(capturedRequest);
        Assert.Equal(HttpMethod.Get, capturedRequest!.Method);
        Assert.Equal(
            "https://example.test/api/GetFeeds?apikey=test-key&siteUrl=https%3A%2F%2Fexample.test",
            capturedRequest.RequestUri!.AbsoluteUri);
        Assert.Single(result.Sitemaps);
        Assert.Equal("https://example.test/sitemap.xml", result.Sitemaps[0].Url);
        Assert.Equal(DateTimeOffset.FromUnixTimeMilliseconds(1732612952000), result.Sitemaps[0].LastCrawled);
        Assert.Equal(12, result.Sitemaps[0].UrlCount);
    }

    [Fact]
    public async Task GetSitemapDetails_UsesFeedUrlQuery_AndDeserializesObjectPayload()
    {
        HttpRequestMessage? capturedRequest = null;
        const string responseJson = """{"d":{"Url":"https://example.test/sitemap.xml","Type":"SitemapIndex","Compressed":true,"FileSize":4096,"LastCrawled":"/Date(1732612952000+0000)/","Submitted":"/Date(1732526552000+0000)/","Status":"Success","UrlCount":44}}""";

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            return JsonResponse(HttpStatusCode.OK, responseJson);
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.GetSitemapDetailsAsync("https://example.test", "https://example.test/sitemap.xml");

        Assert.NotNull(capturedRequest);
        Assert.Equal(HttpMethod.Get, capturedRequest!.Method);
        Assert.Equal(
            "https://example.test/api/GetFeedDetails?apikey=test-key&siteUrl=https%3A%2F%2Fexample.test&feedUrl=https%3A%2F%2Fexample.test%2Fsitemap.xml",
            capturedRequest.RequestUri!.AbsoluteUri);
        Assert.NotNull(result.Sitemap);
        Assert.True(result.Sitemap!.Compressed);
        Assert.Equal("SitemapIndex", result.Sitemap.Type);
        Assert.Equal(DateTimeOffset.FromUnixTimeMilliseconds(1732526552000), result.Sitemap.Submitted);
    }

    [Fact]
    public async Task GetSitemapDetails_AcceptsArrayPayload_AndUsesFirstItem()
    {
        const string responseJson = """{"d":[{"Url":"https://example.test/sitemap-a.xml","Type":"Sitemap","Compressed":false,"FileSize":100,"LastCrawled":"/Date(1732612952000+0000)/","Submitted":"/Date(1732526552000+0000)/","Status":"Success","UrlCount":10},{"Url":"https://example.test/sitemap-b.xml","Type":"Sitemap","Compressed":false,"FileSize":200,"LastCrawled":"/Date(1732612952000+0000)/","Submitted":"/Date(1732526552000+0000)/","Status":"Ignored","UrlCount":20}]}""";
        var client = new BingWebmasterClient(
            new HttpClient(new FakeMessageHandler(_ => JsonResponse(HttpStatusCode.OK, responseJson))),
            "test-key",
            "https://example.test/api");

        var result = await client.GetSitemapDetailsAsync("https://example.test", "https://example.test/sitemap.xml");

        Assert.NotNull(result.Sitemap);
        Assert.Equal("https://example.test/sitemap-a.xml", result.Sitemap!.Url);
        Assert.Equal(10, result.Sitemap.UrlCount);
    }

    [Fact]
    public async Task SubmitSitemap_UsesPostAndCamelCaseBody_AndDeserializesEnvelope()
    {
        HttpRequestMessage? capturedRequest = null;
        string? body = null;

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            body = request.Content is null
                ? null
                : request.Content.ReadAsStringAsync().GetAwaiter().GetResult();
            return JsonResponse(HttpStatusCode.OK, """{"d": true}""");
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.SubmitSitemapAsync("https://example.test", "https://example.test/sitemap.xml");

        Assert.NotNull(capturedRequest);
        Assert.Equal(HttpMethod.Post, capturedRequest!.Method);
        Assert.Equal("https://example.test/api/SubmitFeed?apikey=test-key", capturedRequest.RequestUri!.AbsoluteUri);
        Assert.Equal("""{"siteUrl":"https://example.test","feedUrl":"https://example.test/sitemap.xml"}""", body);
        Assert.Equal("https://example.test/sitemap.xml", result.FeedUrl);
        Assert.True(result.Success);
    }

    [Fact]
    public async Task SubmitUrl_UsesPostAndCamelCaseBody_AndDeserializesEnvelope()
    {
        HttpRequestMessage? capturedRequest = null;
        string? body = null;

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            body = request.Content is null
                ? null
                : request.Content.ReadAsStringAsync().GetAwaiter().GetResult();
            return JsonResponse(HttpStatusCode.OK, """{"d": true}""");
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.SubmitUrlAsync("https://example.test", "https://example.test/page-a");

        Assert.NotNull(capturedRequest);
        Assert.Equal(HttpMethod.Post, capturedRequest!.Method);
        Assert.Equal("https://example.test/api/SubmitUrl?apikey=test-key", capturedRequest.RequestUri!.AbsoluteUri);
        Assert.Equal("""{"siteUrl":"https://example.test","url":"https://example.test/page-a"}""", body);
        Assert.Equal("https://example.test/page-a", result.Url);
        Assert.True(result.Success);
    }

    [Fact]
    public async Task GetUrlSubmissionQuota_UsesSiteUrlQuery_AndDeserializesQuota()
    {
        HttpRequestMessage? capturedRequest = null;
        const string responseJson = """{"d":{"DailyQuota":10,"MonthlyQuota":300}}""";

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            return JsonResponse(HttpStatusCode.OK, responseJson);
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.GetUrlSubmissionQuotaAsync("https://example.test");

        Assert.NotNull(capturedRequest);
        Assert.Equal(HttpMethod.Get, capturedRequest!.Method);
        Assert.Equal(
            "https://example.test/api/GetUrlSubmissionQuota?apikey=test-key&siteUrl=https%3A%2F%2Fexample.test",
            capturedRequest.RequestUri!.AbsoluteUri);
        Assert.Equal(10, result.DailyQuota);
        Assert.Equal(300, result.MonthlyQuota);
    }

    [Fact]
    public async Task GetCrawlStats_UsesSiteUrlQuery_AndParsesDates()
    {
        HttpRequestMessage? capturedRequest = null;
        const string responseJson = """{"d":[{"Date":"/Date(1732612952000+0000)/","CrawledPages":9,"CrawlErrors":1,"InIndex":7,"InLinks":6,"Code2xx":5,"Code301":1,"Code302":0,"Code4xx":2,"Code5xx":3,"AllOtherCodes":4,"BlockedByRobotsTxt":5,"ContainsMalware":6}]}""";

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            return JsonResponse(HttpStatusCode.OK, responseJson);
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.GetCrawlStatsAsync("https://example.test");

        Assert.NotNull(capturedRequest);
        Assert.Equal(HttpMethod.Get, capturedRequest!.Method);
        Assert.Equal(
            "https://example.test/api/GetCrawlStats?apikey=test-key&siteUrl=https%3A%2F%2Fexample.test",
            capturedRequest.RequestUri!.AbsoluteUri);
        Assert.Single(result.Rows);
        Assert.Equal(DateTimeOffset.FromUnixTimeMilliseconds(1732612952000), result.Rows[0].Date);
        Assert.Equal(5, result.Rows[0].BlockedByRobotsTxt);
        Assert.Equal(6, result.Rows[0].ContainsMalware);
    }

    [Fact]
    public async Task GetUrlInfo_UsesUrlQuery_AndParsesDates()
    {
        HttpRequestMessage? capturedRequest = null;
        const string responseJson = """{"d":{"Url":"https://example.test/page-a","IsPage":true,"HttpStatus":200,"DocumentSize":1234,"AnchorCount":5,"DiscoveryDate":"/Date(1732612952000+0000)/","LastCrawledDate":"/Date(1732699352000+0000)/","TotalChildUrlCount":8}}""";

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            return JsonResponse(HttpStatusCode.OK, responseJson);
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.GetUrlInfoAsync("https://example.test", "https://example.test/page-a");

        Assert.NotNull(capturedRequest);
        Assert.Equal(HttpMethod.Get, capturedRequest!.Method);
        Assert.Equal(
            "https://example.test/api/GetUrlInfo?apikey=test-key&siteUrl=https%3A%2F%2Fexample.test&url=https%3A%2F%2Fexample.test%2Fpage-a",
            capturedRequest.RequestUri!.AbsoluteUri);
        Assert.Equal("https://example.test/page-a", result.Url);
        Assert.Equal(DateTimeOffset.FromUnixTimeMilliseconds(1732612952000), result.DiscoveryDate);
        Assert.Equal(DateTimeOffset.FromUnixTimeMilliseconds(1732699352000), result.LastCrawledDate);
        Assert.Equal(8, result.TotalChildUrlCount);
    }

    [Fact]
    public async Task GetUrlTrafficInfo_UsesUrlQuery_AndDeserializesMetrics()
    {
        HttpRequestMessage? capturedRequest = null;
        const string responseJson = """{"d":{"Url":"https://example.test/page-a","IsPage":false,"Clicks":12,"Impressions":34}}""";

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            return JsonResponse(HttpStatusCode.OK, responseJson);
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.GetUrlTrafficInfoAsync("https://example.test", "https://example.test/page-a");

        Assert.NotNull(capturedRequest);
        Assert.Equal(HttpMethod.Get, capturedRequest!.Method);
        Assert.Equal(
            "https://example.test/api/GetUrlTrafficInfo?apikey=test-key&siteUrl=https%3A%2F%2Fexample.test&url=https%3A%2F%2Fexample.test%2Fpage-a",
            capturedRequest.RequestUri!.AbsoluteUri);
        Assert.False(result.IsPage);
        Assert.Equal(12, result.Clicks);
        Assert.Equal(34, result.Impressions);
    }

    [Fact]
    public async Task GetLinkCounts_UsesPageParameter_AndDeserializesLinks()
    {
        HttpRequestMessage? capturedRequest = null;
        const string responseJson = """{"d":{"Links":[{"Count":9,"Url":"https://ref.example.test/a"}],"TotalPages":4}}""";

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            return JsonResponse(HttpStatusCode.OK, responseJson);
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.GetLinkCountsAsync("https://example.test", 2);

        Assert.NotNull(capturedRequest);
        Assert.Equal(HttpMethod.Get, capturedRequest!.Method);
        Assert.Equal(
            "https://example.test/api/GetLinkCounts?apikey=test-key&siteUrl=https%3A%2F%2Fexample.test&page=2",
            capturedRequest.RequestUri!.AbsoluteUri);
        Assert.Equal(4, result.TotalPages);
        Assert.Single(result.Links);
        Assert.Equal(9, result.Links[0].Count);
        Assert.Equal("https://ref.example.test/a", result.Links[0].Url);
    }

    [Fact]
    public async Task GetRankAndTrafficStats_UsesSiteUrlQuery_AndParsesRows()
    {
        HttpRequestMessage? capturedRequest = null;
        const string responseJson = """{"d":[{"Date":"/Date(1732612952000+0000)/","Clicks":17,"Impressions":29}]}""";

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            return JsonResponse(HttpStatusCode.OK, responseJson);
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.GetRankAndTrafficStatsAsync("https://example.test");

        Assert.NotNull(capturedRequest);
        Assert.Equal(HttpMethod.Get, capturedRequest!.Method);
        Assert.Equal(
            "https://example.test/api/GetRankAndTrafficStats?apikey=test-key&siteUrl=https%3A%2F%2Fexample.test",
            capturedRequest.RequestUri!.AbsoluteUri);
        Assert.Single(result.Rows);
        Assert.Equal(DateTimeOffset.FromUnixTimeMilliseconds(1732612952000), result.Rows[0].Date);
        Assert.Equal(29, result.Rows[0].Impressions);
    }

    [Fact]
    public async Task GetQueryStats_UsesSiteUrlQuery_AndParsesRows()
    {
        HttpRequestMessage? capturedRequest = null;
        const string responseJson = """{"d":[{"Query":"bing webmaster api","Date":"/Date(1732612952000+0000)/","Clicks":4,"Impressions":15,"AvgClickPosition":3,"AvgImpressionPosition":6}]}""";

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            return JsonResponse(HttpStatusCode.OK, responseJson);
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.GetQueryStatsAsync("https://example.test");

        Assert.NotNull(capturedRequest);
        Assert.Equal(HttpMethod.Get, capturedRequest!.Method);
        Assert.Equal(
            "https://example.test/api/GetQueryStats?apikey=test-key&siteUrl=https%3A%2F%2Fexample.test",
            capturedRequest.RequestUri!.AbsoluteUri);
        Assert.Single(result.Rows);
        Assert.Equal("bing webmaster api", result.Rows[0].Query);
        Assert.Equal(DateTimeOffset.FromUnixTimeMilliseconds(1732612952000), result.Rows[0].Date);
        Assert.Equal(6, result.Rows[0].AvgImpressionPosition);
    }

    [Fact]
    public async Task GetQueryPageStats_UsesQueryParameter_AndMapsWireQueryFieldToPage()
    {
        HttpRequestMessage? capturedRequest = null;
        const string responseJson = """{"d":[{"Query":"https://example.test/page-b","Date":"/Date(1732612952000+0000)/","Clicks":8,"Impressions":13,"AvgClickPosition":2,"AvgImpressionPosition":5}]}""";

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            return JsonResponse(HttpStatusCode.OK, responseJson);
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.GetQueryPageStatsAsync("https://example.test", "bing webmaster api");

        Assert.NotNull(capturedRequest);
        Assert.Equal(HttpMethod.Get, capturedRequest!.Method);
        Assert.Equal(
            "https://example.test/api/GetQueryPageStats?apikey=test-key&siteUrl=https%3A%2F%2Fexample.test&query=bing%20webmaster%20api",
            capturedRequest.RequestUri!.AbsoluteUri);
        Assert.Single(result.Rows);
        Assert.Equal("https://example.test/page-b", result.Rows[0].Page);
        Assert.Equal(8, result.Rows[0].Clicks);
    }

    [Fact]
    public async Task GetPageQueryStats_UsesPageParameter_AndParsesRows()
    {
        HttpRequestMessage? capturedRequest = null;
        const string responseJson = """{"d":[{"Query":"bing webmaster api","Date":"/Date(1732612952000+0000)/","Clicks":6,"Impressions":14,"AvgClickPosition":1,"AvgImpressionPosition":4}]}""";

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            return JsonResponse(HttpStatusCode.OK, responseJson);
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.GetPageQueryStatsAsync("https://example.test", "https://example.test/page-b");

        Assert.NotNull(capturedRequest);
        Assert.Equal(HttpMethod.Get, capturedRequest!.Method);
        Assert.Equal(
            "https://example.test/api/GetPageQueryStats?apikey=test-key&siteUrl=https%3A%2F%2Fexample.test&page=https%3A%2F%2Fexample.test%2Fpage-b",
            capturedRequest.RequestUri!.AbsoluteUri);
        Assert.Single(result.Rows);
        Assert.Equal("bing webmaster api", result.Rows[0].Query);
        Assert.Equal(4, result.Rows[0].AvgImpressionPosition);
    }

    [Fact]
    public async Task RemoveSite_UsesPostAndTreatsNullPayloadAsSuccess()
    {
        HttpRequestMessage? capturedRequest = null;
        string? body = null;

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            body = request.Content is null ? null : request.Content.ReadAsStringAsync().GetAwaiter().GetResult();
            return JsonResponse(HttpStatusCode.OK, """{"d":null}""");
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.RemoveSiteAsync("https://example.test");

        Assert.NotNull(capturedRequest);
        Assert.Equal(HttpMethod.Post, capturedRequest!.Method);
        Assert.Equal("https://example.test/api/RemoveSite?apikey=test-key", capturedRequest.RequestUri!.AbsoluteUri);
        Assert.Equal("""{"siteUrl":"https://example.test"}""", body);
        Assert.True(result.Success);
    }

    [Fact]
    public async Task GetSiteRoles_UsesIncludeAllSubdomainsQuery_AndDecodesRole()
    {
        HttpRequestMessage? capturedRequest = null;
        const string responseJson = """{"d":[{"Date":"/Date(1732612952000+0000)/","DelegatedCode":"abc123","DelegatorEmail":"owner@example.test","DelegatedCodeOwnerEmail":"verifier@example.test","Email":"reader@example.test","Expired":false,"Role":2,"Site":"https://example.test","VerificationSite":"https://example.test"}]}""";

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            return JsonResponse(HttpStatusCode.OK, responseJson);
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.GetSiteRolesAsync("https://example.test", true);

        Assert.NotNull(capturedRequest);
        Assert.Equal(
            "https://example.test/api/GetSiteRoles?apikey=test-key&siteUrl=https%3A%2F%2Fexample.test&includeAllSubdomains=true",
            capturedRequest!.RequestUri!.AbsoluteUri);
        Assert.True(result.IncludeAllSubdomains);
        Assert.Single(result.Rows);
        Assert.Equal("ReadWrite", result.Rows[0].Role);
        Assert.Equal("owner@example.test", result.Rows[0].DelegatorEmail);
        Assert.Equal(DateTimeOffset.FromUnixTimeMilliseconds(1732612952000), result.Rows[0].Date);
    }

    [Fact]
    public async Task AddSiteRole_UsesPostAndCamelCaseBody()
    {
        HttpRequestMessage? capturedRequest = null;
        string? body = null;

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            body = request.Content is null ? null : request.Content.ReadAsStringAsync().GetAwaiter().GetResult();
            return JsonResponse(HttpStatusCode.OK, """{"d":null}""");
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.AddSiteRoleAsync(
            "https://example.test",
            "https://blog.example.test",
            "reader@example.test",
            "auth-code",
            true,
            false);

        Assert.NotNull(capturedRequest);
        Assert.Equal(HttpMethod.Post, capturedRequest!.Method);
        Assert.Equal("https://example.test/api/AddSiteRoles?apikey=test-key", capturedRequest.RequestUri!.AbsoluteUri);
        Assert.Equal("""{"siteUrl":"https://example.test","delegatedUrl":"https://blog.example.test","userEmail":"reader@example.test","authenticationCode":"auth-code","isAdministrator":true,"isReadOnly":false}""", body);
        Assert.Equal("https://blog.example.test", result.DelegatedUrl);
        Assert.True(result.IsAdministrator);
    }

    [Fact]
    public async Task RemoveSiteRole_UsesNestedSiteRoleBody_AndEncodesFriendlyRole()
    {
        HttpRequestMessage? capturedRequest = null;
        string? body = null;

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            body = request.Content is null ? null : request.Content.ReadAsStringAsync().GetAwaiter().GetResult();
            return JsonResponse(HttpStatusCode.OK, """{"d":null}""");
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.RemoveSiteRoleAsync("https://example.test", "reader@example.test", "ReadOnly");

        Assert.NotNull(capturedRequest);
        Assert.Equal(HttpMethod.Post, capturedRequest!.Method);
        Assert.Equal("https://example.test/api/RemoveSiteRole?apikey=test-key", capturedRequest.RequestUri!.AbsoluteUri);

        using var document = JsonDocument.Parse(body!);
        var root = document.RootElement;
        Assert.Equal("https://example.test", root.GetProperty("siteUrl").GetString());
        var siteRole = root.GetProperty("siteRole");
        Assert.Equal("reader@example.test", siteRole.GetProperty("Email").GetString());
        Assert.Equal(1, siteRole.GetProperty("Role").GetInt32());
        Assert.Equal("https://example.test", siteRole.GetProperty("Site").GetString());
        Assert.Equal("https://example.test", siteRole.GetProperty("VerificationSite").GetString());
        Assert.StartsWith("/Date(", siteRole.GetProperty("Date").GetString(), StringComparison.Ordinal);
        Assert.Equal("ReadOnly", result.Role);
    }

    [Fact]
    public async Task GetBlockedUrls_UsesSiteUrlQuery_AndDecodesEnums()
    {
        HttpRequestMessage? capturedRequest = null;
        const string responseJson = """{"d":[{"Date":"/Date(1732612952000+0000)/","EntityType":1,"RequestType":0,"Url":"https://example.test/private/"}]}""";

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            return JsonResponse(HttpStatusCode.OK, responseJson);
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.GetBlockedUrlsAsync("https://example.test");

        Assert.NotNull(capturedRequest);
        Assert.Equal(
            "https://example.test/api/GetBlockedUrls?apikey=test-key&siteUrl=https%3A%2F%2Fexample.test",
            capturedRequest!.RequestUri!.AbsoluteUri);
        Assert.Single(result.Rows);
        Assert.Equal("Directory", result.Rows[0].EntityType);
        Assert.Equal("CacheOnly", result.Rows[0].RequestType);
    }

    [Fact]
    public async Task AddBlockedUrl_UsesNestedBlockedUrlBody_AndDefaults()
    {
        HttpRequestMessage? capturedRequest = null;
        string? body = null;

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            body = request.Content is null ? null : request.Content.ReadAsStringAsync().GetAwaiter().GetResult();
            return JsonResponse(HttpStatusCode.OK, """{"d":null}""");
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.AddBlockedUrlAsync("https://example.test", "https://example.test/private/");

        Assert.NotNull(capturedRequest);
        Assert.Equal("https://example.test/api/AddBlockedUrl?apikey=test-key", capturedRequest!.RequestUri!.AbsoluteUri);

        using var document = JsonDocument.Parse(body!);
        var blockedUrl = document.RootElement.GetProperty("blockedUrl");
        Assert.Equal(0, blockedUrl.GetProperty("EntityType").GetInt32());
        Assert.Equal(0, blockedUrl.GetProperty("RequestType").GetInt32());
        Assert.Equal("https://example.test/private/", blockedUrl.GetProperty("Url").GetString());
        Assert.Equal("Page", result.EntityType);
        Assert.Equal("CacheOnly", result.RequestType);
    }

    [Fact]
    public async Task RemoveBlockedUrl_UsesNestedBlockedUrlBody_AndDefaults()
    {
        HttpRequestMessage? capturedRequest = null;
        string? body = null;

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            body = request.Content is null ? null : request.Content.ReadAsStringAsync().GetAwaiter().GetResult();
            return JsonResponse(HttpStatusCode.OK, """{"d":null}""");
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.RemoveBlockedUrlAsync("https://example.test", "https://example.test/private/");

        Assert.NotNull(capturedRequest);
        Assert.Equal("https://example.test/api/RemoveBlockedUrl?apikey=test-key", capturedRequest!.RequestUri!.AbsoluteUri);

        using var document = JsonDocument.Parse(body!);
        var blockedUrl = document.RootElement.GetProperty("blockedUrl");
        Assert.Equal(0, blockedUrl.GetProperty("EntityType").GetInt32());
        Assert.Equal(1, blockedUrl.GetProperty("RequestType").GetInt32());
        Assert.Equal("FullRemoval", result.RequestType);
    }

    [Fact]
    public async Task GetQueryPageDetailStats_UsesQueryAndPageParameters_AndParsesRows()
    {
        HttpRequestMessage? capturedRequest = null;
        const string responseJson = """{"d":[{"Date":"/Date(1732612952000+0000)/","Clicks":9,"Impressions":21,"Position":3}]}""";

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            return JsonResponse(HttpStatusCode.OK, responseJson);
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.GetQueryPageDetailStatsAsync("https://example.test", "bing api", "https://example.test/page-a");

        Assert.NotNull(capturedRequest);
        Assert.Equal(
            "https://example.test/api/GetQueryPageDetailStats?apikey=test-key&siteUrl=https%3A%2F%2Fexample.test&query=bing%20api&page=https%3A%2F%2Fexample.test%2Fpage-a",
            capturedRequest!.RequestUri!.AbsoluteUri);
        Assert.Single(result.Rows);
        Assert.Equal(3, result.Rows[0].Position);
        Assert.Equal(DateTimeOffset.FromUnixTimeMilliseconds(1732612952000), result.Rows[0].Date);
    }

    [Fact]
    public async Task GetQueryTrafficStats_UsesQueryParameter_AndParsesRows()
    {
        HttpRequestMessage? capturedRequest = null;
        const string responseJson = """{"d":[{"Date":"/Date(1732612952000+0000)/","Clicks":5,"Impressions":14}]}""";

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            return JsonResponse(HttpStatusCode.OK, responseJson);
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.GetQueryTrafficStatsAsync("https://example.test", "bing api");

        Assert.NotNull(capturedRequest);
        Assert.Equal(
            "https://example.test/api/GetQueryTrafficStats?apikey=test-key&siteUrl=https%3A%2F%2Fexample.test&query=bing%20api",
            capturedRequest!.RequestUri!.AbsoluteUri);
        Assert.Single(result.Rows);
        Assert.Equal(14, result.Rows[0].Impressions);
    }

    [Fact]
    public async Task GetKeyword_UsesDateRangeParameters_AndMapsSingleObjectPayload()
    {
        HttpRequestMessage? capturedRequest = null;
        const string responseJson = """{"d":{"Query":"bing webmaster tools","BroadImpressions":99,"Impressions":42}}""";

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            return JsonResponse(HttpStatusCode.OK, responseJson);
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.GetKeywordAsync("bing webmaster tools", "US", "en-US", "2026-06-01", "2026-06-30");

        Assert.NotNull(capturedRequest);
        Assert.Equal(
            "https://example.test/api/GetKeyword?apikey=test-key&q=bing%20webmaster%20tools&country=US&language=en-US&startDate=2026-06-01&endDate=2026-06-30",
            capturedRequest!.RequestUri!.AbsoluteUri);
        Assert.True(result.Found);
        Assert.Equal(42, result.Impressions);
        Assert.Equal(99, result.BroadImpressions);
    }

    [Fact]
    public async Task GetKeyword_ReturnsFoundFalse_WhenBingOmitsQuery()
    {
        var client = new BingWebmasterClient(
            new HttpClient(new FakeMessageHandler(_ => JsonResponse(HttpStatusCode.OK, """{"d":{"BroadImpressions":99,"Impressions":42}}"""))),
            "test-key",
            "https://example.test/api");

        var result = await client.GetKeywordAsync("bing webmaster tools", "US", "en-US", "2026-06-01", "2026-06-30");

        Assert.False(result.Found);
        Assert.Equal(0, result.Impressions);
        Assert.Equal(0, result.BroadImpressions);
    }

    [Fact]
    public async Task GetRelatedKeywords_UsesDateRangeParameters_AndParsesRows()
    {
        HttpRequestMessage? capturedRequest = null;
        const string responseJson = """{"d":[{"Query":"bing webmaster api","BroadImpressions":50,"Impressions":20},{"Query":"bing seo tools","BroadImpressions":40,"Impressions":10}]}""";

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            return JsonResponse(HttpStatusCode.OK, responseJson);
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.GetRelatedKeywordsAsync("bing webmaster tools", "US", "en-US", "2026-06-01", "2026-06-30");

        Assert.NotNull(capturedRequest);
        Assert.Equal(
            "https://example.test/api/GetRelatedKeywords?apikey=test-key&q=bing%20webmaster%20tools&country=US&language=en-US&startDate=2026-06-01&endDate=2026-06-30",
            capturedRequest!.RequestUri!.AbsoluteUri);
        Assert.Equal(2, result.RowCount);
        Assert.Equal("bing webmaster api", result.Rows[0].Query);
        Assert.Equal(40, result.Rows[1].BroadImpressions);
    }

    [Fact]
    public async Task GetChildrenUrlInfo_UsesPostWithNestedFilters_AndParsesRows()
    {
        HttpRequestMessage? capturedRequest = null;
        string? body = null;
        const string responseJson = """{"d":[{"Url":"https://example.test/child-a","IsPage":true,"HttpStatus":200,"DocumentSize":123,"AnchorCount":7,"DiscoveryDate":"/Date(1732612952000+0000)/","LastCrawledDate":"/Date(1732699352000+0000)/","TotalChildUrlCount":2}]}""";

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            body = request.Content is null ? null : request.Content.ReadAsStringAsync().GetAwaiter().GetResult();
            return JsonResponse(HttpStatusCode.OK, responseJson);
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.GetChildrenUrlInfoAsync(
            "https://example.test",
            "https://example.test/section/",
            2,
            "LastTwoWeeks",
            "LastMonth",
            "IsMalware",
            "Code4xx");

        Assert.NotNull(capturedRequest);
        Assert.Equal(HttpMethod.Post, capturedRequest!.Method);
        Assert.Equal("https://example.test/api/GetChildrenUrlInfo?apikey=test-key", capturedRequest.RequestUri!.AbsoluteUri);

        using var document = JsonDocument.Parse(body!);
        var root = document.RootElement;
        Assert.Equal("https://example.test", root.GetProperty("siteUrl").GetString());
        Assert.Equal("https://example.test/section/", root.GetProperty("url").GetString());
        Assert.Equal(2, root.GetProperty("page").GetInt32());
        var filters = root.GetProperty("filterProperties");
        Assert.Equal(2, filters.GetProperty("CrawlDateFilter").GetInt32());
        Assert.Equal(2, filters.GetProperty("DiscoveredDateFilter").GetInt32());
        Assert.Equal(2, filters.GetProperty("DocFlagsFilters").GetInt32());
        Assert.Equal(16, filters.GetProperty("HttpCodeFilters").GetInt32());
        Assert.Single(result.Rows);
        Assert.Equal("LastTwoWeeks", result.CrawlDateFilter);
        Assert.Equal("Code4xx", result.HttpCodeFilter);
    }

    [Fact]
    public async Task GetChildrenUrlTrafficInfo_UsesPageParameter_AndParsesRows()
    {
        HttpRequestMessage? capturedRequest = null;
        const string responseJson = """{"d":[{"Url":"https://example.test/child-a","IsPage":false,"Clicks":7,"Impressions":19}]}""";

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            return JsonResponse(HttpStatusCode.OK, responseJson);
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.GetChildrenUrlTrafficInfoAsync("https://example.test", "https://example.test/section/", 3);

        Assert.NotNull(capturedRequest);
        Assert.Equal(
            "https://example.test/api/GetChildrenUrlTrafficInfo?apikey=test-key&siteUrl=https%3A%2F%2Fexample.test&url=https%3A%2F%2Fexample.test%2Fsection%2F&page=3",
            capturedRequest!.RequestUri!.AbsoluteUri);
        Assert.Single(result.Rows);
        Assert.False(result.Rows[0].IsPage);
        Assert.Equal(19, result.Rows[0].Impressions);
    }

    [Fact]
    public async Task FetchUrl_UsesPostAndCamelCaseBody()
    {
        HttpRequestMessage? capturedRequest = null;
        string? body = null;

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            body = request.Content is null ? null : request.Content.ReadAsStringAsync().GetAwaiter().GetResult();
            return JsonResponse(HttpStatusCode.OK, """{"d":null}""");
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.FetchUrlAsync("https://example.test", "https://example.test/page-a");

        Assert.NotNull(capturedRequest);
        Assert.Equal("https://example.test/api/FetchUrl?apikey=test-key", capturedRequest!.RequestUri!.AbsoluteUri);
        Assert.Equal("""{"siteUrl":"https://example.test","url":"https://example.test/page-a"}""", body);
        Assert.True(result.Success);
    }

    [Fact]
    public async Task ListFetchedUrls_UsesSiteUrlQuery_AndParsesRows()
    {
        HttpRequestMessage? capturedRequest = null;
        const string responseJson = """{"d":[{"Date":"/Date(1732612952000+0000)/","Expired":false,"Fetched":true,"Url":"https://example.test/page-a"}]}""";

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            return JsonResponse(HttpStatusCode.OK, responseJson);
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.ListFetchedUrlsAsync("https://example.test");

        Assert.NotNull(capturedRequest);
        Assert.Equal(
            "https://example.test/api/GetFetchedUrls?apikey=test-key&siteUrl=https%3A%2F%2Fexample.test",
            capturedRequest!.RequestUri!.AbsoluteUri);
        Assert.Single(result.Rows);
        Assert.True(result.Rows[0].Fetched);
        Assert.Equal(DateTimeOffset.FromUnixTimeMilliseconds(1732612952000), result.Rows[0].Date);
    }

    [Fact]
    public async Task GetFetchedUrlDetails_UsesUrlQuery_AndDeserializesPayload()
    {
        HttpRequestMessage? capturedRequest = null;
        const string responseJson = """{"d":{"Date":"/Date(1732612952000+0000)/","Document":"PGh0bWw+PC9odG1sPg==","Headers":"SFRUUC8xLjEgMjAwIE9LDQpDb250ZW50LVR5cGU6IHRleHQvaHRtbA==","Status":"Success","Url":"https://example.test/page-a"}}""";

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            return JsonResponse(HttpStatusCode.OK, responseJson);
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.GetFetchedUrlDetailsAsync("https://example.test", "https://example.test/page-a");

        Assert.NotNull(capturedRequest);
        Assert.Equal(
            "https://example.test/api/GetFetchedUrlDetails?apikey=test-key&siteUrl=https%3A%2F%2Fexample.test&url=https%3A%2F%2Fexample.test%2Fpage-a",
            capturedRequest!.RequestUri!.AbsoluteUri);
        Assert.Equal("Success", result.Status);
        Assert.Equal("PGh0bWw+PC9odG1sPg==", result.Document);
    }

    [Fact]
    public async Task RemoveSitemap_UsesPostAndCamelCaseBody()
    {
        HttpRequestMessage? capturedRequest = null;
        string? body = null;

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            body = request.Content is null ? null : request.Content.ReadAsStringAsync().GetAwaiter().GetResult();
            return JsonResponse(HttpStatusCode.OK, """{"d":null}""");
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.RemoveSitemapAsync("https://example.test", "https://example.test/sitemap.xml");

        Assert.NotNull(capturedRequest);
        Assert.Equal("https://example.test/api/RemoveFeed?apikey=test-key", capturedRequest!.RequestUri!.AbsoluteUri);
        Assert.Equal("""{"siteUrl":"https://example.test","feedUrl":"https://example.test/sitemap.xml"}""", body);
        Assert.True(result.Success);
    }

    [Fact]
    public async Task GetSiteMoves_UsesSiteUrlQuery_AndDecodesEnums()
    {
        HttpRequestMessage? capturedRequest = null;
        const string responseJson = """{"d":[{"Date":"/Date(1732612952000+0000)/","MoveScope":2,"MoveType":1,"SourceUrl":"https://example.test/old/","TargetUrl":"https://example.test/new/"}]}""";

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            return JsonResponse(HttpStatusCode.OK, responseJson);
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.GetSiteMovesAsync("https://example.test");

        Assert.NotNull(capturedRequest);
        Assert.Equal(
            "https://example.test/api/GetSiteMoves?apikey=test-key&siteUrl=https%3A%2F%2Fexample.test",
            capturedRequest!.RequestUri!.AbsoluteUri);
        Assert.Single(result.Rows);
        Assert.Equal("Directory", result.Rows[0].MoveScope);
        Assert.Equal("Global", result.Rows[0].MoveType);
    }

    [Fact]
    public async Task SubmitSiteMove_UsesNestedSettingsBody_AndEncodesEnums()
    {
        HttpRequestMessage? capturedRequest = null;
        string? body = null;

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            body = request.Content is null ? null : request.Content.ReadAsStringAsync().GetAwaiter().GetResult();
            return JsonResponse(HttpStatusCode.OK, """{"d":null}""");
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.SubmitSiteMoveAsync(
            "https://example.test",
            "https://example.test/old/",
            "https://example.test/new/",
            "Global",
            "Directory");

        Assert.NotNull(capturedRequest);
        Assert.Equal("https://example.test/api/SubmitSiteMove?apikey=test-key", capturedRequest!.RequestUri!.AbsoluteUri);

        using var document = JsonDocument.Parse(body!);
        var settings = document.RootElement.GetProperty("settings");
        Assert.Equal(1, settings.GetProperty("MoveType").GetInt32());
        Assert.Equal(2, settings.GetProperty("MoveScope").GetInt32());
        Assert.Equal("https://example.test/old/", settings.GetProperty("SourceUrl").GetString());
        Assert.Equal("https://example.test/new/", settings.GetProperty("TargetUrl").GetString());
        Assert.Equal("Global", result.MoveType);
        Assert.Equal("Directory", result.MoveScope);
    }

    [Fact]
    public async Task SubmitContent_UsesCamelCaseBody_AndEncodesDynamicServing()
    {
        HttpRequestMessage? capturedRequest = null;
        string? body = null;

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            body = request.Content is null ? null : request.Content.ReadAsStringAsync().GetAwaiter().GetResult();
            return JsonResponse(HttpStatusCode.OK, """{"d":null}""");
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.SubmitContentAsync(
            "https://example.test",
            "https://example.test/page-a",
            "SFRUUC8xLjEgMjAwIE9LDQoNCjxodG1sPg==",
            "eyJAZ3JhcGgiOltdfQ==",
            "Tablet");

        Assert.NotNull(capturedRequest);
        Assert.Equal("https://example.test/api/SubmitContent?apikey=test-key", capturedRequest!.RequestUri!.AbsoluteUri);
        Assert.Equal("""{"siteUrl":"https://example.test","url":"https://example.test/page-a","httpMessage":"SFRUUC8xLjEgMjAwIE9LDQoNCjxodG1sPg==","structuredData":"eyJAZ3JhcGgiOltdfQ==","dynamicServing":4}""", body);
        Assert.Equal("Tablet", result.DynamicServing);
    }

    [Fact]
    public async Task SubmitContent_ThrowsForInvalidDynamicServing()
    {
        var client = new BingWebmasterClient(
            new HttpClient(new FakeMessageHandler(_ => JsonResponse(HttpStatusCode.OK, """{"d":null}"""))),
            "test-key",
            "https://example.test/api");

        var exception = await Assert.ThrowsAsync<ArgumentException>(() => client.SubmitContentAsync(
            "https://example.test",
            "https://example.test/page-a",
            "aGVsbG8=",
            "e30=",
            "SmartFridge"));

        Assert.Contains("Expected one of", exception.Message);
    }

    [Fact]
    public async Task GetContentSubmissionQuota_UsesSiteUrlQuery_AndDeserializesQuota()
    {
        HttpRequestMessage? capturedRequest = null;
        const string responseJson = """{"d":{"DailyQuota":12,"MonthlyQuota":345}}""";

        var handler = new FakeMessageHandler(request =>
        {
            capturedRequest = request;
            return JsonResponse(HttpStatusCode.OK, responseJson);
        });

        var client = new BingWebmasterClient(new HttpClient(handler), "test-key", "https://example.test/api");
        var result = await client.GetContentSubmissionQuotaAsync("https://example.test");

        Assert.NotNull(capturedRequest);
        Assert.Equal(
            "https://example.test/api/GetContentSubmissionQuota?apikey=test-key&siteUrl=https%3A%2F%2Fexample.test",
            capturedRequest!.RequestUri!.AbsoluteUri);
        Assert.Equal(12, result.DailyQuota);
        Assert.Equal(345, result.MonthlyQuota);
    }
}
