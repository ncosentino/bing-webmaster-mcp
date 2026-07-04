using BingWebmasterMcp.Config;
using Xunit;

namespace BingWebmasterMcp.Tests;

public sealed class ApiKeyResolverTests
{
    [Fact]
    public void Resolve_FlagValue_ReturnsFlag()
    {
        var resolver = new ApiKeyResolver("BING_WEBMASTER_API_KEY");
        var result = resolver.Resolve("flag-key");
        Assert.Equal("flag-key", result);
    }

    [Fact]
    public void Resolve_EnvValue_ReturnsEnvVar_WhenFlagMissing()
    {
        const string envVarName = "BING_WEBMASTER_API_KEY_TEST";
        Environment.SetEnvironmentVariable(envVarName, "env-key");
        try
        {
            var resolver = new ApiKeyResolver(envVarName);
            var result = resolver.Resolve(null);
            Assert.Equal("env-key", result);
        }
        finally
        {
            Environment.SetEnvironmentVariable(envVarName, null);
        }
    }

    [Fact]
    public void Resolve_DotEnvValue_ReturnsFileValue_WhenHigherPrioritySourcesMissing()
    {
        const string envVarName = "BING_INDEXNOW_KEY_TEST";
        var folder = Path.Combine(AppContext.BaseDirectory, "resolver-fixtures");
        Directory.CreateDirectory(folder);
        var dotEnvPath = Path.Combine(folder, "bing-webmaster-mcp.env");
        File.WriteAllText(dotEnvPath, """
# comment
BING_INDEXNOW_KEY_TEST="dotenv-key"
""");

        try
        {
            Environment.SetEnvironmentVariable(envVarName, null);
            var resolver = new ApiKeyResolver(envVarName, dotEnvPath);
            var result = resolver.Resolve(null);
            Assert.Equal("dotenv-key", result);
        }
        finally
        {
            if (File.Exists(dotEnvPath))
                File.Delete(dotEnvPath);
            if (Directory.Exists(folder))
                Directory.Delete(folder);
        }
    }
}
