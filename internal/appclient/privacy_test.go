package appclient

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type countingTransport struct{ calls int }

func (t *countingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.calls++
	return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Header: make(http.Header)}, nil
}

func TestRedirectToOriginalDeveloperIsBlocked(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://api.tcno.co/sw/stats/", http.StatusFound)
	}))
	defer server.Close()
	client := &http.Client{Transport: offlineRoundTripper{base: http.DefaultTransport}}
	_, err := client.Get(server.URL)
	if !errors.Is(err, ErrOriginalDeveloperService) {
		t.Fatalf("redirect was not blocked: %v", err)
	}
}

func TestOriginalDeveloperRequestsNeverReachNetwork(t *testing.T) {
	for _, host := range []string{"tcno.co", "api.tcno.co", "s.tcno.co", "API.TCNO.CO."} {
		t.Run(host, func(t *testing.T) {
			base := &countingTransport{}
			req, err := http.NewRequest(http.MethodPost, "https://"+host+"/sw/crash/", http.NoBody)
			if err != nil {
				t.Fatal(err)
			}
			_, err = (offlineRoundTripper{base: base}).RoundTrip(req)
			if !errors.Is(err, ErrOriginalDeveloperService) || base.calls != 0 {
				t.Fatalf("request reached transport: calls=%d, err=%v", base.calls, err)
			}
		})
	}
}

func TestPrivacyPolicyPreservesGitHubAndPlatformRequests(t *testing.T) {
	for _, host := range []string{"api.github.com", "raw.githubusercontent.com", "api.steampowered.com", "example.com"} {
		base := &countingTransport{}
		req, err := http.NewRequest(http.MethodGet, "https://"+host+"/", nil)
		if err != nil {
			t.Fatal(err)
		}
		_, err = (offlineRoundTripper{base: base}).RoundTrip(req)
		if err != nil || base.calls != 1 {
			t.Fatalf("%s blocked: calls=%d err=%v", host, base.calls, err)
		}
	}
}
