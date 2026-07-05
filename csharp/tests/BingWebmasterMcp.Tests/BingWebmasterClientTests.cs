using System.Net;
using System.Text;
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
    public async Task AddSite_UsesPostAndCamelCaseBody_AndDeserializesEnvelope()
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
        Assert.True(result.Success);
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
}
