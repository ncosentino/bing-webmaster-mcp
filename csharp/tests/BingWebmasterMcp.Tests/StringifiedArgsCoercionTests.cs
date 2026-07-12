using System.Text.Json;
using ModelContextProtocol.Protocol;
using Xunit;

namespace BingWebmasterMcp.Tests;

/// <summary>Tests defensive coercion of stringified URL arrays.</summary>
public sealed class StringifiedArgsCoercionTests
{
    [Fact]
    public void CoerceStringifiedArrayArgs_StringifiedArray_IsReplaced()
    {
        var request = new CallToolRequestParams
        {
            Name = "submit_url_batch",
            Arguments = new Dictionary<string, JsonElement>
            {
                ["url_list"] = JsonOf("\"[\\\"https://example.test/a\\\"]\""),
            },
        };

        StringifiedArgsCoercion.CoerceStringifiedArrayArgs(
            request,
            StringifiedArgsCoercion.ToolArrayFields);

        Assert.Equal(JsonValueKind.Array, request.Arguments["url_list"].ValueKind);
    }

    [Fact]
    public void CoerceStringifiedArrayArgs_GenuineArray_IsUnchanged()
    {
        var request = new CallToolRequestParams
        {
            Name = "submit_url_indexnow",
            Arguments = new Dictionary<string, JsonElement>
            {
                ["url_list"] = JsonOf("[\"https://example.test/a\"]"),
            },
        };

        StringifiedArgsCoercion.CoerceStringifiedArrayArgs(
            request,
            StringifiedArgsCoercion.ToolArrayFields);

        Assert.Equal(JsonValueKind.Array, request.Arguments["url_list"].ValueKind);
    }

    private static JsonElement JsonOf(string json)
    {
        using var document = JsonDocument.Parse(json);
        return document.RootElement.Clone();
    }
}
