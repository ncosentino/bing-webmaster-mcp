using System.Runtime.CompilerServices;

namespace BingWebmasterMcp.Config;

/// <summary>Resolves an API key from multiple sources.</summary>
/// <remarks>Priority: CLI argument &gt; environment variable &gt; .env file.</remarks>
internal sealed class ApiKeyResolver(string envVarName, string dotEnvFile = ".env")
{
    /// <summary>Returns the API key from the highest-priority available source.</summary>
    internal string? Resolve(string? flagValue)
    {
        if (!string.IsNullOrWhiteSpace(flagValue))
            return flagValue;

        var envValue = Environment.GetEnvironmentVariable(envVarName);
        if (!string.IsNullOrWhiteSpace(envValue))
            return envValue;

        return ReadFromDotEnv();
    }

    [MethodImpl(MethodImplOptions.NoInlining)]
    private string? ReadFromDotEnv()
    {
        if (!File.Exists(dotEnvFile))
            return null;

        foreach (var line in File.ReadLines(dotEnvFile))
        {
            var trimmed = line.Trim();
            if (trimmed.StartsWith('#') || trimmed.Length == 0)
                continue;

            var prefix = envVarName + "=";
            if (trimmed.StartsWith(prefix, StringComparison.Ordinal))
            {
                var value = trimmed[prefix.Length..].Trim('"', '\'');
                return string.IsNullOrWhiteSpace(value) ? null : value;
            }
        }

        return null;
    }
}
