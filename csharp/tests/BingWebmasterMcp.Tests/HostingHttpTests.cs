using System.Net;
using System.Text;
using System.Text.Json;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Http;
using ModelContextProtocol.Client;
using Xunit;

namespace BingWebmasterMcp.Tests;

/// <summary>Tests the real Streamable HTTP host and MCP request pipeline.</summary>
public sealed class HostingHttpTests
{
    private sealed class FakeHandler(string responseBody) : HttpMessageHandler
    {
        protected override Task<HttpResponseMessage> SendAsync(
            HttpRequestMessage request,
            CancellationToken cancellationToken) =>
            Task.FromResult(new HttpResponseMessage(HttpStatusCode.OK)
            {
                Content = new StringContent(responseBody, Encoding.UTF8, "application/json"),
            });
    }

    [Fact]
    public async Task BuildHttpHost_ServesAllTools()
    {
        await using var app = Hosting.BuildHttpHost([], "test-key", "indexnow-key", port: 0);
        await app.StartAsync();
        try
        {
            await using var client = await ConnectAsync(app);
            var tools = await client.ListToolsAsync();

            Assert.Equal(55, tools.Count);
            Assert.Contains(tools, tool => tool.Name == "list_sites");
            Assert.Contains(tools, tool => tool.Name == "submit_url_indexnow");
            Assert.Contains(tools, tool => tool.Name == "submit_content");
        }
        finally
        {
            await app.StopAsync();
        }
    }

    [Fact]
    public async Task BuildHttpHost_CallListSites_ReturnsSuccess()
    {
        var handler = new FakeHandler(
            """{"d":[{"Url":"https://example.test","IsVerified":true,"DnsVerificationCode":"dns","AuthenticationCode":"meta"}]}""");
        await using var app = Hosting.BuildHttpHost(
            [],
            "test-key",
            "indexnow-key",
            port: 0,
            httpMessageHandler: handler,
            webmasterBaseUrlOverride: "https://example.test/api");
        await app.StartAsync();
        try
        {
            await using var client = await ConnectAsync(app);
            var result = await client.CallToolAsync(
                "list_sites",
                new Dictionary<string, object?>());

            Assert.NotEqual(true, result.IsError);
        }
        finally
        {
            await app.StopAsync();
        }
    }

    [Fact]
    public async Task BuildHttpHost_ServesHealth()
    {
        await using var app = Hosting.BuildHttpHost([], "test-key", "indexnow-key", port: 0);
        await app.StartAsync();
        try
        {
            using var httpClient = new HttpClient();
            using var response = await httpClient.GetAsync(GetEndpoint(app, Hosting.HealthPath));
            response.EnsureSuccessStatusCode();

            using var document = JsonDocument.Parse(await response.Content.ReadAsStringAsync());
            Assert.Equal("ok", document.RootElement.GetProperty("status").GetString());
            Assert.Equal(
                "bing-webmaster-mcp",
                document.RootElement.GetProperty("service").GetString());
            Assert.Equal("http", document.RootElement.GetProperty("transport").GetString());
        }
        finally
        {
            await app.StopAsync();
        }
    }

    [Fact]
    public async Task BuildHttpHost_DefaultsToLoopbackBinding()
    {
        await using var app = Hosting.BuildHttpHost([], "test-key", "indexnow-key", port: 0);
        await app.StartAsync();
        try
        {
            var bound = new Uri(app.Urls.First());
            Assert.Equal(ServerOptions.DefaultListenAddress, bound.Host);
            Assert.Equal(Hosting.DefaultAllowedHosts, app.Configuration["AllowedHosts"]);
        }
        finally
        {
            await app.StopAsync();
        }
    }

    [Fact]
    public async Task BuildHttpHost_RejectsCrossSiteFetchMetadataWithoutOrigin()
    {
        await using var app = Hosting.BuildHttpHost([], "test-key", "indexnow-key", port: 0);
        await app.StartAsync();
        try
        {
            using var httpClient = new HttpClient();
            using var request = new HttpRequestMessage(
                HttpMethod.Post,
                GetEndpoint(app, Hosting.McpPath));
            request.Headers.Add("Sec-Fetch-Site", "cross-site");
            request.Content = new StringContent("{}", Encoding.UTF8, "application/json");

            using var response = await httpClient.SendAsync(request);

            Assert.Equal(HttpStatusCode.Forbidden, response.StatusCode);
        }
        finally
        {
            await app.StopAsync();
        }
    }

    [Fact]
    public async Task BuildHttpHost_RejectsCrossSiteOriginWithTrailingSlash()
    {
        await using var app = Hosting.BuildHttpHost([], "test-key", "indexnow-key", port: 0);
        await app.StartAsync();
        try
        {
            using var httpClient = new HttpClient();
            using var request = new HttpRequestMessage(
                HttpMethod.Post,
                GetEndpoint(app, Hosting.McpPath + "/"));
            request.Headers.Add("Origin", "https://evil.example");
            request.Content = new StringContent("{}", Encoding.UTF8, "application/json");

            using var response = await httpClient.SendAsync(request);

            Assert.Equal(HttpStatusCode.Forbidden, response.StatusCode);
        }
        finally
        {
            await app.StopAsync();
        }
    }

    [Fact]
    public async Task BuildHttpHost_AllowsSafeCrossSiteMethod()
    {
        await using var app = Hosting.BuildHttpHost([], "test-key", "indexnow-key", port: 0);
        await app.StartAsync();
        try
        {
            using var httpClient = new HttpClient();
            using var request = new HttpRequestMessage(
                HttpMethod.Get,
                GetEndpoint(app, Hosting.McpPath));
            request.Headers.Add("Sec-Fetch-Site", "cross-site");
            request.Headers.Add("Origin", "https://evil.example");

            using var response = await httpClient.SendAsync(request);

            Assert.NotEqual(HttpStatusCode.Forbidden, response.StatusCode);
        }
        finally
        {
            await app.StopAsync();
        }
    }

    [Fact]
    public void IsCrossOriginRequestAllowed_AcceptsSameOrigin()
    {
        var context = new DefaultHttpContext();
        context.Request.Method = HttpMethods.Post;
        context.Request.Scheme = "http";
        context.Request.Host = new HostString("127.0.0.1", 8080);
        context.Request.Headers.Origin = "http://127.0.0.1:8080";

        Assert.True(Hosting.IsCrossOriginRequestAllowed(context.Request));
    }

    [Fact]
    public async Task BuildHttpHost_RejectsDisallowedHost()
    {
        await using var app = Hosting.BuildHttpHost([], "test-key", "indexnow-key", port: 0);
        await app.StartAsync();
        try
        {
            using var httpClient = new HttpClient();
            using var request = new HttpRequestMessage(
                HttpMethod.Get,
                GetEndpoint(app, Hosting.HealthPath));
            request.Headers.Host = "evil.example";

            using var response = await httpClient.SendAsync(request);

            Assert.Equal(HttpStatusCode.BadRequest, response.StatusCode);
        }
        finally
        {
            await app.StopAsync();
        }
    }

    [Fact]
    public async Task BuildHttpHost_ShutdownRequiresBearerToken()
    {
        await using var app = Hosting.BuildHttpHost(
            [],
            "test-key",
            "indexnow-key",
            port: 0,
            shutdownToken: "secret-token");
        await app.StartAsync();

        using var httpClient = new HttpClient();
        using var rejectedRequest = new HttpRequestMessage(
            HttpMethod.Post,
            GetEndpoint(app, Hosting.ShutdownPath));
        rejectedRequest.Headers.Authorization =
            new System.Net.Http.Headers.AuthenticationHeaderValue("Bearer", "wrong-token");
        using var rejectedResponse = await httpClient.SendAsync(rejectedRequest);
        Assert.Equal(HttpStatusCode.Unauthorized, rejectedResponse.StatusCode);

        var stopping = new TaskCompletionSource<bool>(
            TaskCreationOptions.RunContinuationsAsynchronously);
        app.Lifetime.ApplicationStopping.Register(() => stopping.SetResult(true));
        using var acceptedRequest = new HttpRequestMessage(
            HttpMethod.Post,
            GetEndpoint(app, Hosting.ShutdownPath));
        acceptedRequest.Headers.Authorization =
            new System.Net.Http.Headers.AuthenticationHeaderValue("Bearer", "secret-token");
        using var acceptedResponse = await httpClient.SendAsync(acceptedRequest);

        Assert.Equal(HttpStatusCode.Accepted, acceptedResponse.StatusCode);
        await stopping.Task.WaitAsync(TimeSpan.FromSeconds(5));
    }

    private static async Task<McpClient> ConnectAsync(WebApplication app) =>
        await McpClient.CreateAsync(new HttpClientTransport(
            new HttpClientTransportOptions
            {
                Endpoint = GetEndpoint(app, Hosting.McpPath),
            }));

    private static Uri GetEndpoint(WebApplication app, string path)
    {
        var bound = new Uri(app.Urls.First());
        return new UriBuilder(bound)
        {
            Host = ServerOptions.DefaultListenAddress,
            Path = path,
        }.Uri;
    }
}
