package api

import (
	"errors"
	"strings"
)

var ErrReportingDisabled = errors.New("diagnostic and statistics uploads are disabled in this fork")

func UserAgent(version string) string {
	return "TcNo-Acc-Switcher/" + strings.TrimSpace(version)
}
