using System.Globalization;

namespace BingWebmasterMcp.BingWebmaster;

/// <summary>Parses Bing Webmaster's legacy .NET JSON date format.</summary>
internal static class BingDateParser
{
    private const string Prefix = "/Date(";
    private const string Suffix = ")/";

    internal static DateTimeOffset Parse(string value)
        => TryParse(value, out var parsed)
            ? parsed
            : throw new FormatException($"Invalid Bing .NET date value: '{value}'.");

    internal static bool TryParse(string? value, out DateTimeOffset parsed)
    {
        parsed = default;

        if (string.IsNullOrWhiteSpace(value) ||
            !value.StartsWith(Prefix, StringComparison.Ordinal) ||
            !value.EndsWith(Suffix, StringComparison.Ordinal))
        {
            return false;
        }

        var inner = value[Prefix.Length..^Suffix.Length];
        var numericLength = inner.Length;
        for (var i = 1; i < inner.Length; i++)
        {
            if (inner[i] is '+' or '-')
            {
                numericLength = i;
                break;
            }
        }

        var millisecondsPart = inner[..numericLength];
        if (!long.TryParse(millisecondsPart, NumberStyles.Integer, CultureInfo.InvariantCulture, out var milliseconds))
            return false;

        parsed = DateTimeOffset.FromUnixTimeMilliseconds(milliseconds);
        return true;
    }

    internal static string Format(DateTimeOffset value)
        => string.Create(
            CultureInfo.InvariantCulture,
            $"/Date({value.ToUnixTimeMilliseconds()}+0000)/");
}
