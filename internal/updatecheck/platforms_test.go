package updatecheck

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestForkUpdateURLs(t *testing.T) {
	if got, want := latestReleaseAPIURL(), "https://api.github.com/repos/2-sa/Switcher/releases/latest"; got != want {
		t.Fatalf("latestReleaseAPIURL() = %q, want %q", got, want)
	}
	if got, want := LatestReleaseURL(), "https://github.com/2-sa/Switcher/releases/latest"; got != want {
		t.Fatalf("LatestReleaseURL() = %q, want %q", got, want)
	}
	if got, want := PlatformsJSONRawURL("Platforms.json"), "https://raw.githubusercontent.com/2-sa/Switcher/refs/heads/go/Platforms.json"; got != want {
		t.Fatalf("PlatformsJSONRawURL() = %q, want %q", got, want)
	}
}

func TestFetchLatestVersionReadsGitHubRelease(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/2-sa/Switcher/releases/latest" {
			t.Errorf("request path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"tag_name":"v2026.10.16","name":"October update","body":"Release notes"}`))
	}))
	defer server.Close()

	version, message, err := fetchLatestVersion(context.Background(), server.Client(), "2026.10.15", server.URL+"/repos/2-sa/Switcher/releases/latest")
	if err != nil {
		t.Fatal(err)
	}
	if version != "v2026.10.16" || message != "Release notes" {
		t.Fatalf("got version %q and message %q", version, message)
	}
}

func TestParsePlatformsJSONVersion(t *testing.T) {
	t.Parallel()
	v, err := ParsePlatformsJSONVersion([]byte(`{"Platforms":{},"Version":"4.0.1"}`))
	if err != nil || v != "4.0.1" {
		t.Fatalf("got %q err=%v", v, err)
	}
}

func TestIsVersionNewer(t *testing.T) {
	t.Parallel()
	if !IsVersionNewer("4.0.1", "4.0.0") {
		t.Fatal("expected 4.0.1 newer than 4.0.0")
	}
	if IsVersionNewer("4.0.0", "4.0.0") {
		t.Fatal("expected equal versions not newer")
	}
	if IsVersionNewer("4.0.0", "4.0.1") {
		t.Fatal("expected older remote rejected")
	}
	if !IsVersionNewer("4.0.0", "") {
		t.Fatal("expected empty current to accept remote")
	}
}
