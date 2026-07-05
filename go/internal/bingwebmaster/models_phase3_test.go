package bingwebmaster

import "testing"

func TestDecodePhase3Enums(t *testing.T) {
	t.Parallel()

	if got := decodeCountryRegionSettingsType(3); got != "Subdomain" {
		t.Fatalf("decodeCountryRegionSettingsType(3) = %q, want %q", got, "Subdomain")
	}
	if got := decodePagePreviewBlockReason(4); got != "Other" {
		t.Fatalf("decodePagePreviewBlockReason(4) = %q, want %q", got, "Other")
	}
}

func TestEncodePhase3Enums(t *testing.T) {
	t.Parallel()

	if got, err := encodeCountryRegionSettingsType("Domain"); err != nil {
		t.Fatalf("encodeCountryRegionSettingsType error = %v", err)
	} else if got != 2 {
		t.Fatalf("encodeCountryRegionSettingsType = %d, want 2", got)
	}

	if got, err := encodePagePreviewBlockReason("IllegalContent"); err != nil {
		t.Fatalf("encodePagePreviewBlockReason error = %v", err)
	} else if got != 3 {
		t.Fatalf("encodePagePreviewBlockReason = %d, want 3", got)
	}
}
