package bingwebmaster

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestParseDotNetDate(t *testing.T) {
	t.Parallel()

	got, err := parseDotNetDate("/Date(1732612952000+0000)/")
	if err != nil {
		t.Fatalf("parseDotNetDate error = %v", err)
	}

	want := time.UnixMilli(1732612952000).UTC()
	if !got.Equal(want) {
		t.Fatalf("got %s, want %s", got.Format(time.RFC3339), want.Format(time.RFC3339))
	}
}

func TestParseDotNetDate_IgnoresOffsetSuffix(t *testing.T) {
	t.Parallel()

	got, err := parseDotNetDate("/Date(1732612952000-0700)/")
	if err != nil {
		t.Fatalf("parseDotNetDate error = %v", err)
	}

	want := time.UnixMilli(1732612952000).UTC()
	if !got.Equal(want) {
		t.Fatalf("got %s, want %s", got.Format(time.RFC3339), want.Format(time.RFC3339))
	}
}

func TestParseDotNetDate_Invalid(t *testing.T) {
	t.Parallel()

	if _, err := parseDotNetDate("2024-11-26T00:00:00Z"); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestWireTime_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	var got wireTime
	if err := json.Unmarshal([]byte(`"/Date(1732612952000+0000)/"`), &got); err != nil {
		t.Fatalf("Unmarshal error = %v", err)
	}
	want := time.UnixMilli(1732612952000).UTC()
	if !got.Equal(want) {
		t.Fatalf("got %s, want %s", got.Format(time.RFC3339), want.Format(time.RFC3339))
	}
}

func TestMapFeed_ConvertsDates(t *testing.T) {
	t.Parallel()

	raw := rawFeed{
		URL:         "https://example.test/sitemap.xml",
		Type:        "Sitemap",
		Compressed:  false,
		FileSize:    12,
		LastCrawled: wireTime{Time: time.UnixMilli(1732612952000).UTC()},
		Submitted:   wireTime{Time: time.UnixMilli(1732612953000).UTC()},
		Status:      "OK",
		URLCount:    7,
	}

	mapped := mapFeed(raw)
	if mapped.LastCrawled == nil || mapped.Submitted == nil {
		t.Fatal("expected converted time pointers")
	}
	if mapped.URL != raw.URL {
		t.Fatalf("URL = %q, want %q", mapped.URL, raw.URL)
	}
}

func TestDecodeCrawlIssueFlags(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		value int
		want  []string
	}{
		{
			name:  "multiple flags",
			value: 20,
			want:  []string{"Code4xx", "BlockedByRobotsTxt"},
		},
		{
			name:  "single flag",
			value: 128,
			want:  []string{"DnsErrors"},
		},
		{
			name:  "none",
			value: 0,
			want:  []string{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := decodeCrawlIssueFlags(testCase.value)
			if !reflect.DeepEqual(got, testCase.want) {
				t.Fatalf("decodeCrawlIssueFlags(%d) = %#v, want %#v", testCase.value, got, testCase.want)
			}
			if got == nil {
				t.Fatal("decodeCrawlIssueFlags returned nil; want empty slice for none")
			}
		})
	}
}
