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

	// Undocumented test-only hooks: point the compiled binary at a local mock
	// server for end-to-end testing. Left unset, both clients target the real
	// Bing endpoints.
	bingwebmaster.SetBaseURL(os.Getenv("BING_WEBMASTER_API_BASE_URL"))
	indexnow.SetBaseURL(os.Getenv("BING_INDEXNOW_API_BASE_URL"))

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

	mcp.AddTool(srv,
		&mcp.Tool{Name: "remove_site", Description: "Remove a site from Bing Webmaster Tools."},
		toolHandler("removing site", func(ctx context.Context, input removeSiteInput) (any, error) {
			return bingClient.RemoveSite(ctx, input.SiteURL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_site_roles", Description: "List delegated site roles for a Bing Webmaster Tools site."},
		toolHandler("getting site roles", func(ctx context.Context, input getSiteRolesInput) (any, error) {
			return bingClient.GetSiteRoles(ctx, input.SiteURL, input.IncludeAllSubdomains)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "add_site_role", Description: "Delegate Bing Webmaster Tools access for a site to another user."},
		toolHandler("adding site role", func(ctx context.Context, input addSiteRoleInput) (any, error) {
			// IsReadOnly is a *bool (not bool) specifically so an omitted argument
			// can be told apart from an explicit false -- bool's zero value would
			// otherwise silently default to false instead of the documented true.
			isReadOnly := true
			if input.IsReadOnly != nil {
				isReadOnly = *input.IsReadOnly
			}
			return bingClient.AddSiteRole(ctx, input.SiteURL, input.DelegatedURL, input.UserEmail, input.AuthenticationCode, input.IsAdministrator, isReadOnly)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "remove_site_role", Description: "Remove delegated Bing Webmaster Tools access for a site user."},
		toolHandler("removing site role", func(ctx context.Context, input removeSiteRoleInput) (any, error) {
			return bingClient.RemoveSiteRole(ctx, input.SiteURL, input.Email, input.Role)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_blocked_urls", Description: "List blocked URL removal requests for a site."},
		toolHandler("getting blocked URLs", func(ctx context.Context, input getBlockedURLsInput) (any, error) {
			return bingClient.GetBlockedURLs(ctx, input.SiteURL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "add_blocked_url", Description: "Add a blocked URL removal request for a site."},
		toolHandler("adding blocked URL", func(ctx context.Context, input addBlockedURLInput) (any, error) {
			return bingClient.AddBlockedURL(ctx, input.SiteURL, input.URL, input.EntityType, input.RequestType)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "remove_blocked_url", Description: "Remove a blocked URL removal request for a site."},
		toolHandler("removing blocked URL", func(ctx context.Context, input removeBlockedURLInput) (any, error) {
			return bingClient.RemoveBlockedURL(ctx, input.SiteURL, input.URL, input.EntityType, input.RequestType)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_query_parameters", Description: "List query parameter normalization settings for a site."},
		toolHandler("getting query parameters", func(ctx context.Context, input getQueryParametersInput) (any, error) {
			return bingClient.GetQueryParameters(ctx, input.SiteURL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "add_query_parameter", Description: "Add a query parameter normalization setting for a site."},
		toolHandler("adding query parameter", func(ctx context.Context, input addQueryParameterInput) (any, error) {
			return bingClient.AddQueryParameter(ctx, input.SiteURL, input.QueryParameter)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "remove_query_parameter", Description: "Remove a query parameter normalization setting for a site."},
		toolHandler("removing query parameter", func(ctx context.Context, input removeQueryParameterInput) (any, error) {
			return bingClient.RemoveQueryParameter(ctx, input.SiteURL, input.QueryParameter)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "enable_disable_query_parameter", Description: "Enable or disable a query parameter normalization setting for a site."},
		toolHandler("changing query parameter state", func(ctx context.Context, input enableDisableQueryParameterInput) (any, error) {
			return bingClient.EnableDisableQueryParameter(ctx, input.SiteURL, input.QueryParameter, input.IsEnabled)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_country_region_settings", Description: "List Bing geo-targeting country and region settings for a site."},
		toolHandler("getting country region settings", func(ctx context.Context, input getCountryRegionSettingsInput) (any, error) {
			return bingClient.GetCountryRegionSettings(ctx, input.SiteURL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "add_country_region_settings", Description: "Add a Bing geo-targeting country or region setting for a site."},
		toolHandler("adding country region settings", func(ctx context.Context, input addCountryRegionSettingsInput) (any, error) {
			return bingClient.AddCountryRegionSettings(ctx, input.SiteURL, input.TwoLetterIsoCountryCode, input.SettingsType, input.URL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "remove_country_region_settings", Description: "Remove a Bing geo-targeting country or region setting for a site."},
		toolHandler("removing country region settings", func(ctx context.Context, input removeCountryRegionSettingsInput) (any, error) {
			return bingClient.RemoveCountryRegionSettings(ctx, input.SiteURL, input.TwoLetterIsoCountryCode, input.SettingsType, input.URL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_connected_pages", Description: "List connected pages for a Bing Webmaster Tools site."},
		toolHandler("getting connected pages", func(ctx context.Context, input getConnectedPagesInput) (any, error) {
			return bingClient.GetConnectedPages(ctx, input.SiteURL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "add_connected_page", Description: "Add a connected page for a Bing Webmaster Tools site."},
		toolHandler("adding connected page", func(ctx context.Context, input addConnectedPageInput) (any, error) {
			return bingClient.AddConnectedPage(ctx, input.SiteURL, input.MasterURL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_active_page_preview_blocks", Description: "List active page preview blocks for a site."},
		toolHandler("getting active page preview blocks", func(ctx context.Context, input getActivePagePreviewBlocksInput) (any, error) {
			return bingClient.GetActivePagePreviewBlocks(ctx, input.SiteURL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "add_page_preview_block", Description: "Add a page preview block for a site URL."},
		toolHandler("adding page preview block", func(ctx context.Context, input addPagePreviewBlockInput) (any, error) {
			return bingClient.AddPagePreviewBlock(ctx, input.SiteURL, input.URL, input.Reason)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "remove_page_preview_block", Description: "Remove a page preview block for a site URL."},
		toolHandler("removing page preview block", func(ctx context.Context, input removePagePreviewBlockInput) (any, error) {
			return bingClient.RemovePagePreviewBlock(ctx, input.SiteURL, input.URL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_query_page_detail_stats", Description: "Get daily stats for a specific query and page combination."},
		toolHandler("getting query page detail stats", func(ctx context.Context, input getQueryPageDetailStatsInput) (any, error) {
			return bingClient.GetQueryPageDetailStats(ctx, input.SiteURL, input.Query, input.Page)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_query_traffic_stats", Description: "Get daily clicks and impressions for a specific query."},
		toolHandler("getting query traffic stats", func(ctx context.Context, input getQueryTrafficStatsInput) (any, error) {
			return bingClient.GetQueryTrafficStats(ctx, input.SiteURL, input.Query)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_keyword", Description: "Get market-wide keyword impressions for a date range."},
		toolHandler("getting keyword", func(ctx context.Context, input getKeywordInput) (any, error) {
			return bingClient.GetKeyword(ctx, input.Query, input.Country, input.Language, input.StartDate, input.EndDate)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_related_keywords", Description: "Get related market-wide keywords and impression counts for a date range."},
		toolHandler("getting related keywords", func(ctx context.Context, input getRelatedKeywordsInput) (any, error) {
			return bingClient.GetRelatedKeywords(ctx, input.Query, input.Country, input.Language, input.StartDate, input.EndDate)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_children_url_info", Description: "Get child URL crawl information for a parent URL."},
		toolHandler("getting child URL info", func(ctx context.Context, input getChildrenURLInfoInput) (any, error) {
			return bingClient.GetChildrenURLInfo(ctx, input.SiteURL, input.URL, input.Page, input.CrawlDateFilter, input.DiscoveredDateFilter, input.DocFlagsFilter, input.HTTPCodeFilter)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_children_url_traffic_info", Description: "Get child URL traffic information for a parent URL."},
		toolHandler("getting child URL traffic info", func(ctx context.Context, input getChildrenURLTrafficInfoInput) (any, error) {
			return bingClient.GetChildrenURLTrafficInfo(ctx, input.SiteURL, input.URL, input.Page)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "fetch_url", Description: "Ask Bing Webmaster Tools to fetch a URL."},
		toolHandler("fetching URL", func(ctx context.Context, input fetchURLInput) (any, error) {
			return bingClient.FetchURL(ctx, input.SiteURL, input.URL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "list_fetched_urls", Description: "List URLs previously fetched through Bing Webmaster Tools."},
		toolHandler("listing fetched URLs", func(ctx context.Context, input listFetchedURLsInput) (any, error) {
			return bingClient.ListFetchedURLs(ctx, input.SiteURL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_fetched_url_details", Description: "Get Bing's stored fetch response for a URL."},
		toolHandler("getting fetched URL details", func(ctx context.Context, input getFetchedURLDetailsInput) (any, error) {
			return bingClient.GetFetchedURLDetails(ctx, input.SiteURL, input.URL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "remove_sitemap", Description: "Remove a sitemap from Bing Webmaster Tools."},
		toolHandler("removing sitemap", func(ctx context.Context, input removeSitemapInput) (any, error) {
			return bingClient.RemoveSitemap(ctx, input.SiteURL, input.FeedURL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_site_moves", Description: "List Bing Webmaster Tools site move settings."},
		toolHandler("getting site moves", func(ctx context.Context, input getSiteMovesInput) (any, error) {
			return bingClient.GetSiteMoves(ctx, input.SiteURL)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "submit_site_move", Description: "Submit a Bing Webmaster Tools site move request."},
		toolHandler("submitting site move", func(ctx context.Context, input submitSiteMoveInput) (any, error) {
			return bingClient.SubmitSiteMove(ctx, input.SiteURL, input.SourceURL, input.TargetURL, input.MoveType, input.MoveScope)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "submit_content", Description: "Submit cached content and structured data for a URL."},
		toolHandler("submitting content", func(ctx context.Context, input submitContentInput) (any, error) {
			return bingClient.SubmitContent(ctx, input.SiteURL, input.URL, input.HTTPMessage, input.StructuredData, input.DynamicServing)
		}),
	)

	mcp.AddTool(srv,
		&mcp.Tool{Name: "get_content_submission_quota", Description: "Get Bing content submission quotas for a site."},
		toolHandler("getting content submission quota", func(ctx context.Context, input getContentSubmissionQuotaInput) (any, error) {
			return bingClient.GetContentSubmissionQuota(ctx, input.SiteURL)
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

type removeSiteInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL to remove from Bing Webmaster Tools."`
}

type getSiteRolesInput struct {
	SiteURL              string `json:"site_url" jsonschema:"The site URL whose delegated roles should be returned."`
	IncludeAllSubdomains bool   `json:"include_all_subdomains,omitempty" jsonschema:"When true, include delegated roles across all subdomains. Defaults to false."`
}

type addSiteRoleInput struct {
	SiteURL            string `json:"site_url" jsonschema:"The site URL that owns the delegated role."`
	DelegatedURL       string `json:"delegated_url" jsonschema:"The delegated site or subdomain URL being granted."`
	UserEmail          string `json:"user_email" jsonschema:"The email address receiving access."`
	AuthenticationCode string `json:"authentication_code" jsonschema:"The Bing verification or authentication code required by the API."`
	IsAdministrator    bool   `json:"is_administrator,omitempty" jsonschema:"Grant administrator access when true."`
	IsReadOnly         *bool  `json:"is_read_only,omitempty" jsonschema:"Grant read-only access when true. Defaults to true when omitted."`
}

type removeSiteRoleInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL whose delegated role should be removed."`
	Email   string `json:"email" jsonschema:"The email address whose role should be removed."`
	Role    string `json:"role" jsonschema:"The role to remove. Allowed values: Administrator, ReadOnly, ReadWrite."`
}

type getBlockedURLsInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL whose blocked URL removals should be listed."`
}

type addBlockedURLInput struct {
	SiteURL     string `json:"site_url" jsonschema:"The site URL that owns the blocked URL request."`
	URL         string `json:"url" jsonschema:"The URL or directory to block."`
	EntityType  string `json:"entity_type,omitempty" jsonschema:"Entity type. Allowed values: Page or Directory. Defaults to Page."`
	RequestType string `json:"request_type,omitempty" jsonschema:"Removal request type. Allowed values: CacheOnly or FullRemoval. Defaults to CacheOnly."`
}

type removeBlockedURLInput struct {
	SiteURL     string `json:"site_url" jsonschema:"The site URL that owns the blocked URL request."`
	URL         string `json:"url" jsonschema:"The URL or directory whose block should be removed."`
	EntityType  string `json:"entity_type,omitempty" jsonschema:"Entity type. Allowed values: Page or Directory. Defaults to Page."`
	RequestType string `json:"request_type,omitempty" jsonschema:"Removal request type. Allowed values: CacheOnly or FullRemoval. Defaults to FullRemoval."`
}

type getQueryParametersInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL whose query parameter normalization settings should be listed."`
}

type addQueryParameterInput struct {
	SiteURL        string `json:"site_url" jsonschema:"The site URL that owns the query parameter setting."`
	QueryParameter string `json:"query_parameter" jsonschema:"The query parameter name to add, for example 'utm_campaign'."`
}

type removeQueryParameterInput struct {
	SiteURL        string `json:"site_url" jsonschema:"The site URL that owns the query parameter setting."`
	QueryParameter string `json:"query_parameter" jsonschema:"The query parameter name to remove."`
}

type enableDisableQueryParameterInput struct {
	SiteURL        string `json:"site_url" jsonschema:"The site URL that owns the query parameter setting."`
	QueryParameter string `json:"query_parameter" jsonschema:"The query parameter name whose enabled state should be changed."`
	IsEnabled      bool   `json:"is_enabled" jsonschema:"Set to true to enable the query parameter or false to disable it."`
}

type getCountryRegionSettingsInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL whose Bing geo-targeting settings should be listed."`
}

type addCountryRegionSettingsInput struct {
	SiteURL                 string `json:"site_url" jsonschema:"The site URL that owns the geo-targeting setting."`
	TwoLetterIsoCountryCode string `json:"two_letter_iso_country_code" jsonschema:"The two-letter ISO country code to target, for example 'us'."`
	SettingsType            string `json:"settings_type" jsonschema:"Settings scope type. Allowed values: Page, Directory, Domain, Subdomain."`
	URL                     string `json:"url" jsonschema:"The page, directory, domain, or subdomain URL the setting applies to."`
}

type removeCountryRegionSettingsInput struct {
	SiteURL                 string `json:"site_url" jsonschema:"The site URL that owns the geo-targeting setting."`
	TwoLetterIsoCountryCode string `json:"two_letter_iso_country_code" jsonschema:"The two-letter ISO country code to remove."`
	SettingsType            string `json:"settings_type" jsonschema:"Settings scope type. Allowed values: Page, Directory, Domain, Subdomain."`
	URL                     string `json:"url" jsonschema:"The page, directory, domain, or subdomain URL the setting applies to."`
}

type getConnectedPagesInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL whose connected pages should be listed."`
}

type addConnectedPageInput struct {
	SiteURL   string `json:"site_url" jsonschema:"The site URL that owns the connected page relationship."`
	MasterURL string `json:"master_url" jsonschema:"The connected page or master URL to add."`
}

type getActivePagePreviewBlocksInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL whose active page preview blocks should be listed."`
}

type addPagePreviewBlockInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL that owns the blocked page preview."`
	URL     string `json:"url" jsonschema:"The page URL whose preview should be blocked."`
	Reason  string `json:"reason" jsonschema:"Block reason. Allowed values: AdultContent, Copyright, IllegalContent, Other."`
}

type removePagePreviewBlockInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL that owns the blocked page preview."`
	URL     string `json:"url" jsonschema:"The page URL whose preview block should be removed."`
}

type getQueryPageDetailStatsInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL whose detailed query/page stats should be returned."`
	Query   string `json:"query" jsonschema:"The search query to inspect."`
	Page    string `json:"page" jsonschema:"The page URL to inspect."`
}

type getQueryTrafficStatsInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL whose query traffic stats should be returned."`
	Query   string `json:"query" jsonschema:"The search query to inspect."`
}

type getKeywordInput struct {
	Query     string `json:"query" jsonschema:"The market-wide search query to inspect."`
	Country   string `json:"country" jsonschema:"The Bing country market code, for example 'US'."`
	Language  string `json:"language" jsonschema:"The Bing language code, for example 'en'."`
	StartDate string `json:"start_date" jsonschema:"The inclusive start date in YYYY-MM-DD format."`
	EndDate   string `json:"end_date" jsonschema:"The inclusive end date in YYYY-MM-DD format."`
}

type getRelatedKeywordsInput struct {
	Query     string `json:"query" jsonschema:"The market-wide search query whose related keywords should be returned."`
	Country   string `json:"country" jsonschema:"The Bing country market code, for example 'US'."`
	Language  string `json:"language" jsonschema:"The Bing language code, for example 'en'."`
	StartDate string `json:"start_date" jsonschema:"The inclusive start date in YYYY-MM-DD format."`
	EndDate   string `json:"end_date" jsonschema:"The inclusive end date in YYYY-MM-DD format."`
}

type getChildrenURLInfoInput struct {
	SiteURL              string `json:"site_url" jsonschema:"The site URL that owns the parent URL."`
	URL                  string `json:"url" jsonschema:"The parent URL whose child crawl information should be returned."`
	Page                 int    `json:"page,omitempty" jsonschema:"Zero-based page number. Defaults to 0 when omitted."`
	CrawlDateFilter      string `json:"crawl_date_filter,omitempty" jsonschema:"Allowed values: Any, LastWeek, LastTwoWeeks, LastThreeWeeks. Defaults to Any."`
	DiscoveredDateFilter string `json:"discovered_date_filter,omitempty" jsonschema:"Allowed values: Any, LastWeek, LastMonth. Defaults to Any."`
	DocFlagsFilter       string `json:"doc_flags_filter,omitempty" jsonschema:"Allowed values: Any, IsBlockedByRobotsTxt, IsMalware. Defaults to Any."`
	HTTPCodeFilter       string `json:"http_code_filter,omitempty" jsonschema:"Allowed values: Any, Code2xx, Code3xx, Code301, Code302, Code4xx, Code5xx, AllOthers. Defaults to Any."`
}

type getChildrenURLTrafficInfoInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL that owns the parent URL."`
	URL     string `json:"url" jsonschema:"The parent URL whose child traffic information should be returned."`
	Page    int    `json:"page,omitempty" jsonschema:"Zero-based page number. Defaults to 0 when omitted."`
}

type fetchURLInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL that owns the URL to fetch."`
	URL     string `json:"url" jsonschema:"The URL Bing should fetch."`
}

type listFetchedURLsInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL whose fetched URLs should be listed."`
}

type getFetchedURLDetailsInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL that owns the fetched URL."`
	URL     string `json:"url" jsonschema:"The URL whose fetch details should be returned."`
}

type removeSitemapInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL that owns the sitemap."`
	FeedURL string `json:"feed_url" jsonschema:"The sitemap feed URL to remove."`
}

type getSiteMovesInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL whose move settings should be returned."`
}

type submitSiteMoveInput struct {
	SiteURL   string `json:"site_url" jsonschema:"The site URL that owns the move submission."`
	SourceURL string `json:"source_url" jsonschema:"The source URL or scope being moved."`
	TargetURL string `json:"target_url" jsonschema:"The target URL or scope for the move."`
	MoveType  string `json:"move_type,omitempty" jsonschema:"Allowed values: Local or Global. Defaults to Local."`
	MoveScope string `json:"move_scope,omitempty" jsonschema:"Allowed values: Domain, Host, Directory. Defaults to Domain."`
}

type submitContentInput struct {
	SiteURL        string `json:"site_url" jsonschema:"The site URL that owns the submitted content."`
	URL            string `json:"url" jsonschema:"The URL whose content is being submitted."`
	HTTPMessage    string `json:"http_message" jsonschema:"The base64-encoded HTTP message payload."`
	StructuredData string `json:"structured_data" jsonschema:"The base64-encoded structured data payload."`
	DynamicServing string `json:"dynamic_serving,omitempty" jsonschema:"Allowed values: None, PcLaptop, Mobile, Amp, Tablet, NonVisualBrowser. Defaults to None."`
}

type getContentSubmissionQuotaInput struct {
	SiteURL string `json:"site_url" jsonschema:"The site URL whose content submission quota should be returned."`
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
