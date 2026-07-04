// Command bing-webmaster-mcp is an MCP server that exposes Bing Webmaster Tools
// and IndexNow as tools for AI assistants. It communicates via STDIO using the MCP protocol.
//
// Usage:
//
//	bing-webmaster-mcp [--api-key <key>] [--indexnow-key <key>]
//
// Credential resolution order: CLI flags > environment variables > .env file.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ncosentino/bing-webmaster-mcp/go/internal/bingwebmaster"
	"github.com/ncosentino/bing-webmaster-mcp/go/internal/config"
	"github.com/ncosentino/bing-webmaster-mcp/go/internal/indexnow"
)

var version = "dev"

// listSitesInputSchema is the explicit no-argument JSON Schema for the list_sites tool.
// Strict MCP clients (e.g. Copilot CLI) reject tools whose schema omits explicit
// properties/required/additionalProperties fields, breaking the entire MCP session.
var listSitesInputSchema = json.RawMessage(`{"type":"object","properties":{},"required":[],"additionalProperties":false}`)

func main() {
	apiKey := flag.String("api-key", "", "Bing Webmaster API key")
	indexNowKey := flag.String("indexnow-key", "", "Optional default IndexNow key")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg := config.Resolve(*apiKey, *indexNowKey)
	if cfg.APIKey == "" {
		slog.Error("no Bing Webmaster API key provided",
			"hint", "set --api-key flag, BING_WEBMASTER_API_KEY env var, or BING_WEBMASTER_API_KEY in .env")
		os.Exit(1)
	}

	bingClient := bingwebmaster.NewClient(cfg.APIKey)
	indexNowClient := indexnow.NewClient(cfg.IndexNowKey)

	srv := newServer(bingClient, indexNowClient)

	slog.Info("bing-webmaster-mcp starting", "version", version, "transport", "stdio")
	if err := srv.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		slog.Error("server stopped with error", "err", err)
		os.Exit(1)
	}
}

func newServer(bingClient *bingwebmaster.Client, indexNowClient *indexnow.Client) *mcp.Server {
	srv := mcp.NewServer(&mcp.Implementation{
		Name:    "bing-webmaster-mcp",
		Version: version,
	}, nil)

	registerTools(srv, bingClient, indexNowClient)
	return srv
}

func registerTools(srv *mcp.Server, bingClient *bingwebmaster.Client, indexNowClient *indexnow.Client) {
	mcp.AddTool(srv,
		&mcp.Tool{Name: "list_sites", Description: "List all Bing Webmaster Tools sites accessible to the configured API key.", InputSchema: listSitesInputSchema},
		toolHandler("listing sites", func(ctx context.Context, _ listSitesInput) (any, error) {
			return bingClient.ListSites(ctx)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "add_site", Description: "Add a site to Bing Webmaster Tools."},
		toolHandler("adding site", func(ctx context.Context, input addSiteInput) (any, error) {
			return bingClient.AddSite(ctx, input.SiteURL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "verify_site", Description: "Verify a site in Bing Webmaster Tools."},
		toolHandler("verifying site", func(ctx context.Context, input verifySiteInput) (any, error) {
			return bingClient.VerifySite(ctx, input.SiteURL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "list_sitemaps", Description: "List submitted sitemaps for a Bing Webmaster Tools site."},
		toolHandler("listing sitemaps", func(ctx context.Context, input listSitemapsInput) (any, error) {
			return bingClient.ListSitemaps(ctx, input.SiteURL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_sitemap_details", Description: "Get sitemap details for a specific submitted sitemap."},
		toolHandler("getting sitemap details", func(ctx context.Context, input getSitemapDetailsInput) (any, error) {
			return bingClient.GetSitemapDetails(ctx, input.SiteURL, input.FeedURL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "submit_sitemap", Description: "Submit a sitemap to Bing Webmaster Tools."},
		toolHandler("submitting sitemap", func(ctx context.Context, input submitSitemapInput) (any, error) {
			return bingClient.SubmitSitemap(ctx, input.SiteURL, input.FeedURL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "submit_url", Description: "Submit a single URL to Bing Webmaster Tools."},
		toolHandler("submitting URL", func(ctx context.Context, input submitURLInput) (any, error) {
			return bingClient.SubmitURL(ctx, input.SiteURL, input.URL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "submit_url_batch", Description: "Submit up to 500 URLs to Bing Webmaster Tools in one request."},
		toolHandler("submitting URL batch", func(ctx context.Context, input submitURLBatchInput) (any, error) {
			return bingClient.SubmitURLBatch(ctx, input.SiteURL, input.URLList)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "submit_url_indexnow", Description: "Submit one or more URLs to Bing using the IndexNow protocol."},
		toolHandler("submitting URLs via IndexNow", func(ctx context.Context, input submitURLIndexNowInput) (any, error) {
			return indexNowClient.SubmitURLs(ctx, input.Host, input.URLList, input.Key, input.KeyLocation)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_url_submission_quota", Description: "Get Bing Webmaster Tools URL submission quotas for a site."},
		toolHandler("getting URL submission quota", func(ctx context.Context, input getURLSubmissionQuotaInput) (any, error) {
			return bingClient.GetURLSubmissionQuota(ctx, input.SiteURL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_crawl_issues", Description: "Get crawl issues for a site."},
		toolHandler("getting crawl issues", func(ctx context.Context, input getCrawlIssuesInput) (any, error) {
			return bingClient.GetCrawlIssues(ctx, input.SiteURL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_crawl_stats", Description: "Get crawl statistics for a site."},
		toolHandler("getting crawl stats", func(ctx context.Context, input getCrawlStatsInput) (any, error) {
			return bingClient.GetCrawlStats(ctx, input.SiteURL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_url_info", Description: "Get Bing Webmaster Tools index information for a URL."},
		toolHandler("getting URL info", func(ctx context.Context, input getURLInfoInput) (any, error) {
			return bingClient.GetURLInfo(ctx, input.SiteURL, input.URL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_url_traffic_info", Description: "Get Bing Webmaster Tools traffic information for a URL."},
		toolHandler("getting URL traffic info", func(ctx context.Context, input getURLTrafficInfoInput) (any, error) {
			return bingClient.GetURLTrafficInfo(ctx, input.SiteURL, input.URL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_url_links", Description: "Get inbound links for a URL."},
		toolHandler("getting URL links", func(ctx context.Context, input getURLLinksInput) (any, error) {
			return bingClient.GetURLLinks(ctx, input.SiteURL, input.Link, input.Page)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_link_counts", Description: "Get pages with inbound link counts for a site."},
		toolHandler("getting link counts", func(ctx context.Context, input getLinkCountsInput) (any, error) {
			return bingClient.GetLinkCounts(ctx, input.SiteURL, input.Page)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_rank_and_traffic_stats", Description: "Get clicks and impressions over time for a site."},
		toolHandler("getting rank and traffic stats", func(ctx context.Context, input getRankAndTrafficStatsInput) (any, error) {
			return bingClient.GetRankAndTrafficStats(ctx, input.SiteURL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_query_stats", Description: "Get top search queries for a site."},
		toolHandler("getting query stats", func(ctx context.Context, input getQueryStatsInput) (any, error) {
			return bingClient.GetQueryStats(ctx, input.SiteURL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_page_stats", Description: "Get top pages for a site."},
		toolHandler("getting page stats", func(ctx context.Context, input getPageStatsInput) (any, error) {
			return bingClient.GetPageStats(ctx, input.SiteURL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_page_query_stats", Description: "Get top queries for a specific page."},
		toolHandler("getting page query stats", func(ctx context.Context, input getPageQueryStatsInput) (any, error) {
			return bingClient.GetPageQueryStats(ctx, input.SiteURL, input.Page)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_query_page_stats", Description: "Get top pages for a specific query."},
		toolHandler("getting query page stats", func(ctx context.Context, input getQueryPageStatsInput) (any, error) {
			return bingClient.GetQueryPageStats(ctx, input.SiteURL, input.Query)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_keyword_stats", Description: "Get market-wide Bing keyword statistics."},
		toolHandler("getting keyword stats", func(ctx context.Context, input getKeywordStatsInput) (any, error) {
			return bingClient.GetKeywordStats(ctx, input.Query, input.Country, input.Language)
		}),
	)
}

type listSitesInput struct{}

type addSiteInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL to add, for example 'https://example.com'."`
}

type verifySiteInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL to verify, for example 'https://example.com'."`
}

type listSitemapsInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL whose sitemaps should be listed."`
}

type getSitemapDetailsInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL that owns the sitemap."`
	FeedURL string `json:"feed_url" jsonschema:"The sitemap feed URL to inspect."`
}

type submitSitemapInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL that owns the sitemap."`
	FeedURL string `json:"feed_url" jsonschema:"The sitemap feed URL to submit."`
}

type submitURLInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL that owns the submitted URL."`
	URL     string `json:"url" jsonschema:"The URL to submit."`
}

type submitURLBatchInput struct {
	SiteURL string   `json:"site_url" jsonschema:"The site URL that owns the submitted URLs."`
	URLList []string `json:"url_list" jsonschema:"Up to 500 URLs to submit."`
}

type submitURLIndexNowInput struct {
	Host        string   `json:"host" jsonschema:"The host name that owns every URL in url_list, for example 'example.com'."`
	URLList     []string `json:"url_list" jsonschema:"One or more URLs to submit through IndexNow."`
	Key         string   `json:"key,omitempty" jsonschema:"Optional IndexNow key override. If omitted, the configured default key is used."`
	KeyLocation string   `json:"key_location,omitempty" jsonschema:"Optional fully qualified URL of the hosted IndexNow key file."`
}

type getURLSubmissionQuotaInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL whose Bing URL submission quota should be returned."`
}

type getCrawlIssuesInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL whose crawl issues should be returned."`
}

type getCrawlStatsInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL whose crawl statistics should be returned."`
}

type getURLInfoInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL that owns the requested URL."`
	URL     string `json:"url" jsonschema:"The URL whose Bing index information should be returned."`
}

type getURLTrafficInfoInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL that owns the requested URL."`
	URL     string `json:"url" jsonschema:"The URL whose Bing traffic information should be returned."`
}

type getURLLinksInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL that owns the requested link target."`
	Link    string `json:"link" jsonschema:"The URL whose inbound links should be returned."`
	Page    int    `json:"page,omitempty" jsonschema:"Zero-based page number. Defaults to 0 when omitted."`
}

type getLinkCountsInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL whose link counts should be returned."`
	Page    int    `json:"page,omitempty" jsonschema:"Zero-based page number. Defaults to 0 when omitted."`
}

type getRankAndTrafficStatsInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL whose clicks and impressions over time should be returned."`
}

type getQueryStatsInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL whose top queries should be returned."`
}

type getPageStatsInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL whose top pages should be returned."`
}

type getPageQueryStatsInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL whose page query stats should be returned."`
	Page    string `json:"page" jsonschema:"The page URL whose queries should be returned."`
}

type getQueryPageStatsInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL whose query page stats should be returned."`
	Query   string `json:"query" jsonschema:"The search query whose top pages should be returned."`
}

type getKeywordStatsInput struct {
	Query    string `json:"query" jsonschema:"The market-wide search query to analyze."`
	Country  string `json:"country" jsonschema:"The Bing country market code, for example 'US'."`
	Language string `json:"language" jsonschema:"The Bing language code, for example 'en'."`
}

func toolHandler[T any](operation string, fn func(context.Context, T) (any, error)) func(context.Context, *mcp.CallToolRequest, T) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input T) (result *mcp.CallToolResult, meta any, err error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				result = jsonResult(map[string]string{"error": fmt.Sprintf("%s: %v", operation, recovered)})
				meta = nil
				err = nil
			}
		}()

		output, callErr := fn(ctx, input)
		if callErr != nil {
			return jsonResult(map[string]string{"error": fmt.Sprintf("%s: %v", operation, callErr)}), nil, nil
		}

		return jsonResult(output), nil, nil
	}
}

func jsonResult(value any) *mcp.CallToolResult {
	payload, err := json.Marshal(value)
	if err != nil {
		fallback, _ := json.Marshal(map[string]string{"error": fmt.Sprintf("marshalling result: %v", err)})
		payload = fallback
	}

	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(payload)}}}
}
