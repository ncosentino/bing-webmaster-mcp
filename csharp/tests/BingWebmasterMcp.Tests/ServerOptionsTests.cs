using Xunit;

namespace BingWebmasterMcp.Tests;

/// <summary>Tests shared-service command-line and environment configuration.</summary>
public sealed class ServerOptionsTests
{
    [Fact]
    public void Parse_DefaultsToStdioAndLoopback()
    {
        var options = ServerOptions.Parse([], _ => null);

        Assert.Equal("stdio", options.Transport);
        Assert.Equal(ServerOptions.DefaultListenAddress, options.ListenAddress);
        Assert.Equal(ServerOptions.DefaultPort, options.Port);
        Assert.Null(options.ShutdownToken);
    }

    [Fact]
    public void Parse_ArgumentsOverrideEnvironment()
    {
        var environment = new Dictionary<string, string?>
        {
            ["MCP_LISTEN_ADDRESS"] = "192.0.2.1",
            ["PORT"] = "9000",
            ["MCP_SHUTDOWN_TOKEN"] = "test-token",
        };

        var options = ServerOptions.Parse(
            [
                "--transport=http",
                "--listen-address",
                "127.0.0.2",
                "--port",
                "8083",
            ],
            name => environment.GetValueOrDefault(name));

        Assert.Equal("http", options.Transport);
        Assert.Equal("127.0.0.2", options.ListenAddress);
        Assert.Equal(8083, options.Port);
        Assert.Equal("test-token", options.ShutdownToken);
    }

    [Theory]
    [InlineData("--transport", "sse")]
    [InlineData("--listen-address", "")]
    [InlineData("--port", "0")]
    [InlineData("--port", "65536")]
    [InlineData("--port", "invalid")]
    public void Parse_InvalidConfiguration_IsRejected(string option, string value)
    {
        var args = option == "--transport"
            ? new[] { option, value }
            : new[] { "--transport", "http", option, value };

        Assert.Throws<ArgumentException>(() =>
            ServerOptions.Parse(args, _ => null));
    }

    [Fact]
    public void Parse_UnknownOption_IsRejected()
    {
        var exception = Assert.Throws<ArgumentException>(() =>
            ServerOptions.Parse(["--tranport", "http"], _ => null));

        Assert.Contains("--tranport", exception.Message, StringComparison.Ordinal);
    }

    [Fact]
    public void Parse_MistypedOptionCasing_IsRejected()
    {
        var exception = Assert.Throws<ArgumentException>(() =>
            ServerOptions.Parse(["--Transport=http"], _ => null));

        Assert.Contains("--Transport", exception.Message, StringComparison.Ordinal);
    }

    [Fact]
    public void GetOption_SupportsEqualsSyntax()
    {
        Assert.Equal(
            "test-key",
            ServerOptions.GetOption(["--api-key=test-key"], "--api-key"));
    }
}
