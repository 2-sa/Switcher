package stability

import "TcNo-Acc-Switcher/internal/api"

// SubmitRating is retained for native API compatibility; no ratings leave this fork.
func SubmitRating(platform string, working bool) {}

// SubmitFeedback refuses legacy native submissions instead of sending text or logs.
func SubmitFeedback(kind, platform, text string, attachLog bool) error {
	return api.ErrReportingDisabled
}
