using BingWebmasterMcp;
using BingWebmasterMcp.Config;
using BingWebmasterMcp.Tools;
using Microsoft.Extensions.Hosting;
using ModelContextProtocol.Server;

if (ServerOptions.IsHelpRequested(args))
{
    await Console.Out.WriteLineAsync(ServerOptions.Usage).ConfigureAwait(false);
    return 0;
}

ServerOptions options;
try
{
    options = ServerOptions.Parse(args);
}
catch (ArgumentException exception)
{
    await Console.Error.WriteLineAsync($"ERROR: {exception.Message}").ConfigureAwait(false);
    return 1;
}

var apiKeyResolver = new ApiKeyResolver("BING_WEBMASTER_API_KEY");
var indexNowKeyResolver = new ApiKeyResolver("BING_INDEXNOW_KEY");
var apiKey = apiKeyResolver.Resolve(
    ServerOptions.GetOption(args, "--api-key"));
var indexNowKey = indexNowKeyResolver.Resolve(
    ServerOptions.GetOption(args, "--indexnow-key"));

if (string.IsNullOrWhiteSpace(apiKey))
{
    await Console.Error.WriteLineAsync(
        "ERROR: No API key provided. Use --api-key <key>, set BING_WEBMASTER_API_KEY env var, or add it to a .env file.")
        .ConfigureAwait(false);
    return 1;
}

var webmasterBaseUrlOverride = NonEmpty(
    Environment.GetEnvironmentVariable("BING_WEBMASTER_API_BASE_URL"));
var indexNowBaseUrlOverride = NonEmpty(
    Environment.GetEnvironmentVariable("BING_INDEXNOW_API_BASE_URL"));

if (options.Transport == "http")
{
    var app = Hosting.BuildHttpHost(
        args,
        apiKey,
        indexNowKey,
        options.Port,
        webmasterBaseUrlOverride: webmasterBaseUrlOverride,
        indexNowBaseUrlOverride: indexNowBaseUrlOverride,
        listenAddress: options.ListenAddress,
        shutdownToken: options.ShutdownToken);
    await app.RunAsync().ConfigureAwait(false);
    return 0;
}

var builder = Host.CreateApplicationBuilder(args);
Hosting.ConfigureCommonServices(
    builder,
    apiKey,
    indexNowKey,
    webmasterBaseUrlOverride: webmasterBaseUrlOverride,
    indexNowBaseUrlOverride: indexNowBaseUrlOverride);
Hosting.ConfigureMcpServer(builder.Services)
    .WithStdioServerTransport()
    .WithTools<BingWebmasterTool>();

var host = builder.Build();
await host.RunAsync().ConfigureAwait(false);
return 0;

static string? NonEmpty(string? value) =>
    string.IsNullOrWhiteSpace(value) ? null : value;
