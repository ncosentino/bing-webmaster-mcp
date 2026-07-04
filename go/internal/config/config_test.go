package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveFlagTakesPriorityPerKey(t *testing.T) {
	t.Setenv(apiKeyEnvVar, "env-api")
	t.Setenv(indexNowEnvVar, "env-index")

	cfg := Resolve("flag-api", "flag-index")
	if cfg.APIKey != "flag-api" {
		t.Fatalf("APIKey = %q, want %q", cfg.APIKey, "flag-api")
	}
	if cfg.IndexNowKey != "flag-index" {
		t.Fatalf("IndexNowKey = %q, want %q", cfg.IndexNowKey, "flag-index")
	}
}

func TestResolveEnvVarFallback(t *testing.T) {
	t.Setenv(apiKeyEnvVar, "env-api")
	t.Setenv(indexNowEnvVar, "env-index")

	cfg := Resolve("", "")
	if cfg.APIKey != "env-api" {
		t.Fatalf("APIKey = %q, want %q", cfg.APIKey, "env-api")
	}
	if cfg.IndexNowKey != "env-index" {
		t.Fatalf("IndexNowKey = %q, want %q", cfg.IndexNowKey, "env-index")
	}
}

func TestResolveDotEnvFallback(t *testing.T) {
	dir := t.TempDir()
	dotEnv := filepath.Join(dir, dotEnvFile)
	contents := "BING_WEBMASTER_API_KEY=dotenv-api\nBING_INDEXNOW_KEY=dotenv-index\n"
	if err := os.WriteFile(dotEnv, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}

	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	cfg := Resolve("", "")
	if cfg.APIKey != "dotenv-api" {
		t.Fatalf("APIKey = %q, want %q", cfg.APIKey, "dotenv-api")
	}
	if cfg.IndexNowKey != "dotenv-index" {
		t.Fatalf("IndexNowKey = %q, want %q", cfg.IndexNowKey, "dotenv-index")
	}
}

func TestResolveKeysAreIndependent(t *testing.T) {
	t.Setenv(apiKeyEnvVar, "env-api")

	cfg := Resolve("", "flag-index")
	if cfg.APIKey != "env-api" {
		t.Fatalf("APIKey = %q, want %q", cfg.APIKey, "env-api")
	}
	if cfg.IndexNowKey != "flag-index" {
		t.Fatalf("IndexNowKey = %q, want %q", cfg.IndexNowKey, "flag-index")
	}
}

func TestResolveEmpty(t *testing.T) {
	cfg := Resolve("", "")
	if cfg.APIKey != "" {
		t.Fatalf("APIKey = %q, want empty", cfg.APIKey)
	}
	if cfg.IndexNowKey != "" {
		t.Fatalf("IndexNowKey = %q, want empty", cfg.IndexNowKey)
	}
}
