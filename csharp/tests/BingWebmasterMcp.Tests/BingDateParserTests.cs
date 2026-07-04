using BingWebmasterMcp.BingWebmaster;
using Xunit;

namespace BingWebmasterMcp.Tests;

public sealed class BingDateParserTests
{
    [Fact]
    public void Parse_WithOffsetSuffix_ReturnsUtcInstant()
    {
        var result = BingDateParser.Parse("/Date(1732612952000+0000)/");
        Assert.Equal(DateTimeOffset.FromUnixTimeMilliseconds(1732612952000), result);
    }

    [Fact]
    public void TryParse_InvalidValue_ReturnsFalse()
    {
        var parsed = BingDateParser.TryParse("2025-01-01", out var result);
        Assert.False(parsed);
        Assert.Equal(default, result);
    }
}
