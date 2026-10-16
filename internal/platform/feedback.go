package platform

import (
	"strings"

	"TcNo-Acc-Switcher/internal/api"
	"TcNo-Acc-Switcher/internal/appclient"
	"TcNo-Acc-Switcher/internal/stability"
)

// SubmitStabilityRating refuses upstream ratings from legacy native clients.
func (p *PlatformService) SubmitStabilityRating(platform string, working bool) error {
	if appclient.IsOfflineMode() {
		return appclient.ErrOfflineMode
	}
	platform = strings.TrimSpace(platform)
	if platform == "" {
		return nil
	}
	return api.ErrReportingDisabled
}

// SubmitFeedback refuses upstream issue and suggestion submissions.
// kind must be "switch_issue" or "feature_suggestion".
func (p *PlatformService) SubmitFeedback(kind, platform, text string, attachLog bool) error {
	if appclient.IsOfflineMode() {
		return appclient.ErrOfflineMode
	}
	kind = strings.TrimSpace(kind)
	text = strings.TrimSpace(text)
	if kind == "" || text == "" {
		return nil
	}
	return stability.SubmitFeedback(kind, strings.TrimSpace(platform), text, attachLog)
}
