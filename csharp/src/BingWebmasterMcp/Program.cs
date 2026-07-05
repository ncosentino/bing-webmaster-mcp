using BingWebmasterMcp.BingWebmaster;
using BingWebmasterMcp.Config;
using BingWebmasterMcp.IndexNow;
using BingWebmasterMcp.Tools;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using ModelContextProtocol.Server;

var apiKeyResolver = new ApiKeyResolver("BING_WEBMASTER_API_KEY");
var indexNowKeyResolver = new ApiKeyResolver("BING_INDEXNOW_KEY");

var apiKey = apiKeyResolver.Resolve(
    args.SkipWhile(a => a != "--api-key").Skip(1).FirstOrDefault());
var indexNowKey = indexNowKeyResolver.Resolve(
    args.SkipWhile(a => a != "--indexnow-key").Skip(1).FirstOrDefault());

if (string.IsNullOrWhiteSpace(apiKey))
{
    await Console.Error.WriteLineAsync(
        "ERROR: No API key provided. Use --api-key <key>, set BING_WEBMASTER_API_KEY env var, or add it to a .env file.")
        .ConfigureAwait(false);
    return 1;
}

// Undocumented test-only hooks: point the compiled binary at a local mock server for
// end-to-end testing. Left unset, both clients target the real Bing endpoints.
var webmasterBaseUrlOverride = NonEmpty(Environment.GetEnvironmentVariable("BING_WEBMASTER_API_BASE_URL"));
var indexNowBaseUrlOverride = NonEmpty(Environment.GetEnvironmentVariable("BING_INDEXNOW_API_BASE_URL"));

static string? NonEmpty(string? value) => string.IsNullOrWhiteSpace(value) ? null : value;

var builder = Host.CreateApplicationBuilder(args);

// All logs must go to stderr to avoid corrupting the MCP STDIO stream.
builder.Logging.AddConsole(o => o.LogToStandardErrorThreshold = LogLevel.Trace);
builder.Logging.SetMinimumLevel(LogLevel.Warning);

builder.Services
    .AddHttpClient(nameof(BingWebmasterClient), http =>
    {
        http.Timeout = TimeSpan.FromSeconds(30);
    });

builder.Services
    .AddHttpClient(nameof(IndexNowClient), http =>
    {
        http.Timeout = TimeSpan.FromSeconds(30);
    });

builder.Services.AddTransient<BingWebmasterClient>(sp =>
{
    var factory = sp.GetRequiredService<IHttpClientFactory>();
    return new BingWebmasterClient(factory.CreateClient(nameof(BingWebmasterClient)), apiKey!, webmasterBaseUrlOverride);
});

builder.Services.AddTransient<IndexNowClient>(sp =>
{
    var factory = sp.GetRequiredService<IHttpClientFactory>();
    return new IndexNowClient(factory.CreateClient(nameof(IndexNowClient)), indexNowKey, indexNowBaseUrlOverride);
});

builder.Services
    .AddMcpServer()
    .WithStdioServerTransport()
    .WithTools<BingWebmasterTool>();

var host = builder.Build();
await host.RunAsync().ConfigureAwait(false);
return 0;
