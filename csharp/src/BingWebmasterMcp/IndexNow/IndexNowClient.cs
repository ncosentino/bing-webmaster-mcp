using System.Net;
using System.Text;
using System.Text.Json;
using BingWebmasterMcp.BingWebmaster;

namespace BingWebmasterMcp.IndexNow;

/// <summary>Thrown when the IndexNow API returns a non-success status code.</summary>
internal sealed class IndexNowApiException : Exception
{
    internal int StatusCode { get; }

    internal IndexNowApiException(int statusCode, string message) : base(message)
        => StatusCode = statusCode;
}

/// <summary>Client for the Bing IndexNow endpoint.</summary>
internal sealed class IndexNowClient(HttpClient httpClient, string? configuredKey, string? baseUrlOverride = null)
{
    private const string DefaultUrl = "https://www.bing.com/indexnow";

    private readonly string? _configuredKey = string.IsNullOrWhiteSpace(configuredKey) ? null : configuredKey.Trim();
    private readonly HttpClient _httpClient = httpClient;
    private readonly string _url = baseUrlOverride ?? DefaultUrl;

    internal async Task<SubmitUrlIndexNowResponse> SubmitUrlIndexNowAsync(
        string host,
        IReadOnlyList<string> urlList,
        string? keyOverride,
        string? keyLocation,
        CancellationToken cancellationToken = default)
    {
        var normalizedHost = RequireText(host, nameof(host));
        var normalizedUrls = NormalizeUrlList(urlList);
        var resolvedKey = string.IsNullOrWhiteSpace(keyOverride) ? _configuredKey : keyOverride.Trim();

        if (string.IsNullOrWhiteSpace(resolvedKey))
        {
            throw new InvalidOperationException(
                "No IndexNow key provided. Use the tool's key argument, set BING_INDEXNOW_KEY, or add it to a .env file.");
        }

        var normalizedKeyLocation = string.IsNullOrWhiteSpace(keyLocation) ? null : keyLocation.Trim();
        var requestBody = new IndexNowRequest
        {
            Host = normalizedHost,
            Key = resolvedKey,
            KeyLocation = normalizedKeyLocation,
            UrlList = normalizedUrls
        };

        var jsonBody = JsonSerializer.Serialize(requestBody, BingWebmasterJsonContext.Default.IndexNowRequest);
        using var request = new HttpRequestMessage(HttpMethod.Post, _url)
        {
            Content = new StringContent(jsonBody, Encoding.UTF8, "application/json")
        };
        request.Content.Headers.ContentType!.CharSet = "utf-8";

        using var response = await _httpClient.SendAsync(request, cancellationToken).ConfigureAwait(false);
        await EnsureSuccessAsync(response, cancellationToken).ConfigureAwait(false);

        return new SubmitUrlIndexNowResponse(
            normalizedHost,
            normalizedUrls,
            normalizedKeyLocation,
            true,
            string.IsNullOrWhiteSpace(keyOverride) ? "configured" : "override",
            DateTimeOffset.UtcNow);
    }

    private static async Task EnsureSuccessAsync(HttpResponseMessage response, CancellationToken cancellationToken)
    {
        if (response.IsSuccessStatusCode)
            return;

        var body = await response.Content.ReadAsStringAsync(cancellationToken).ConfigureAwait(false);
        var reason = response.StatusCode switch
        {
            HttpStatusCode.BadRequest => "bad request",
            HttpStatusCode.Forbidden => "key invalid",
            (HttpStatusCode)422 => "URLs do not belong to host or key does not match",
            (HttpStatusCode)429 => "too many requests",
            _ => string.IsNullOrWhiteSpace(response.ReasonPhrase) ? "request failed" : response.ReasonPhrase
        };
        var snippet = body.Length > 300 ? body[..300] + "..." : body;
        var suffix = string.IsNullOrWhiteSpace(snippet) ? string.Empty : $": {snippet}";
        throw new IndexNowApiException(
            (int)response.StatusCode,
            $"IndexNow returned HTTP {(int)response.StatusCode} {response.StatusCode} ({reason}){suffix}");
    }

    private static string RequireText(string value, string paramName)
        => string.IsNullOrWhiteSpace(value)
            ? throw new ArgumentException("Value is required.", paramName)
            : value.Trim();

    private static IReadOnlyList<string> NormalizeUrlList(IReadOnlyList<string> urlList)
    {
        ArgumentNullException.ThrowIfNull(urlList);

        var normalized = urlList
            .Where(url => !string.IsNullOrWhiteSpace(url))
            .Select(url => url.Trim())
            .ToList();

        if (normalized.Count == 0)
            throw new ArgumentException("At least one URL is required.", nameof(urlList));

        return normalized;
    }
}
