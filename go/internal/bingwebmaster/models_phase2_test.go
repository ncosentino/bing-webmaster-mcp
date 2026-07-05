package bingwebmaster

import (
	"reflect"
	"testing"
	"time"
)

func TestFormatDotNetDate(t *testing.T) {
	t.Parallel()

	got := formatDotNetDate(time.UnixMilli(1732612952000).UTC())
	if got != "/Date(1732612952000+0000)/" {
		t.Fatalf("got %q, want %q", got, "/Date(1732612952000+0000)/")
	}
}

func TestDecodePhase2Enums(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		got  string
		want string
	}{
		{name: "site role", got: decodeSiteRole(1), want: "ReadOnly"},
		{name: "entity type", got: decodeBlockedURLEntityType(1), want: "Directory"},
		{name: "request type", got: decodeBlockedURLRequestType(1), want: "FullRemoval"},
		{name: "move scope", got: decodeMoveScope(2), want: "Directory"},
		{name: "move type", got: decodeMoveType(1), want: "Global"},
		{name: "unknown", got: decodeSiteRole(99), want: "99"},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if testCase.got != testCase.want {
				t.Fatalf("got %q, want %q", testCase.got, testCase.want)
			}
		})
	}
}

func TestEncodePhase2Enums(t *testing.T) {
	t.Parallel()

	got := []int{}
	var err error

	if value, encodeErr := encodeSiteRole("Administrator"); encodeErr != nil {
		err = encodeErr
	} else {
		got = append(got, value)
	}
	if value, encodeErr := encodeBlockedURLEntityType("Page"); encodeErr != nil {
		err = encodeErr
	} else {
		got = append(got, value)
	}
	if value, encodeErr := encodeBlockedURLRequestType("CacheOnly"); encodeErr != nil {
		err = encodeErr
	} else {
		got = append(got, value)
	}
	if value, encodeErr := encodeCrawlDateFilter("LastThreeWeeks"); encodeErr != nil {
		err = encodeErr
	} else {
		got = append(got, value)
	}
	if value, encodeErr := encodeDiscoveredDateFilter("LastMonth"); encodeErr != nil {
		err = encodeErr
	} else {
		got = append(got, value)
	}
	if value, encodeErr := encodeDocFlagsFilter("IsMalware"); encodeErr != nil {
		err = encodeErr
	} else {
		got = append(got, value)
	}
	if value, encodeErr := encodeHTTPCodeFilter("Code5xx"); encodeErr != nil {
		err = encodeErr
	} else {
		got = append(got, value)
	}
	if value, encodeErr := encodeMoveScope("Host"); encodeErr != nil {
		err = encodeErr
	} else {
		got = append(got, value)
	}
	if value, encodeErr := encodeMoveType("Local"); encodeErr != nil {
		err = encodeErr
	} else {
		got = append(got, value)
	}
	if value, encodeErr := encodeDynamicServing("Tablet"); encodeErr != nil {
		err = encodeErr
	} else {
		got = append(got, value)
	}

	if err != nil {
		t.Fatalf("unexpected encode error = %v", err)
	}

	want := []int{0, 0, 0, 4, 2, 2, 32, 1, 0, 4}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestEncodePhase2Enums_Invalid(t *testing.T) {
	t.Parallel()

	if _, err := encodeDynamicServing("Spaceship"); err == nil {
		t.Fatal("expected error, got nil")
	}
}
