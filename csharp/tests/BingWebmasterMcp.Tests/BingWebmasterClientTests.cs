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
}
