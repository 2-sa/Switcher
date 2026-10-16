package platform

import (
	"encoding/json"
	"errors"
	"testing"

	"TcNo-Acc-Switcher/internal/api"
	"TcNo-Acc-Switcher/internal/stability"
	"TcNo-Acc-Switcher/internal/stats"
)

func TestLegacyReportingCannotBeEnabled(t *testing.T) {
	dir := testExeDirWithPortable(t)
	s := defaultSettings()
	if s.StatsEnabled || s.StatsShare || s.CrashReportAutoSubmit {
		t.Fatal("fresh settings must not collect or upload usage data")
	}
	s.StatsEnabled = true // Explicit local-only opt-in remains supported.
	s.StatsShare = true
	s.CrashReportAutoSubmit = true
	if err := SaveAppSettings(dir, s); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadAppSettings(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !loaded.StatsEnabled || loaded.StatsShare || loaded.CrashReportAutoSubmit {
		t.Fatalf("local opt-in/reporting policy lost: %+v", loaded)
	}
	var old AppSettings
	if err := json.Unmarshal([]byte(`{"statsShare":true,"crashReportAutoSubmit":true}`), &old); err != nil {
		t.Fatal(err)
	}
	normalizeAppSettingsDefaults(&old, map[string]json.RawMessage{
		"statsShare": json.RawMessage("true"), "crashReportAutoSubmit": json.RawMessage("true"),
	})
	if old.StatsShare || old.CrashReportAutoSubmit || old.StatsEnabled {
		t.Fatal("old settings must not restore upstream reporting or missing collection opt-in")
	}
}

func TestNativeReportingAPIsCannotRestoreUploads(t *testing.T) {
	service := &PlatformService{}
	if err := service.SetStatsShare(true); !errors.Is(err, api.ErrReportingDisabled) {
		t.Fatalf("sharing could be enabled: %v", err)
	}
	if err := service.SetCrashReportAutoSubmit(true); !errors.Is(err, api.ErrReportingDisabled) {
		t.Fatalf("crash submission could be enabled: %v", err)
	}
	if submitted, err := service.SubmitPendingCrashReport(); submitted || !errors.Is(err, api.ErrReportingDisabled) {
		t.Fatalf("crash upload accepted: submitted=%v err=%v", submitted, err)
	}
	if err := stability.SubmitFeedback("switch_issue", "Steam", "private text", true); !errors.Is(err, api.ErrReportingDisabled) {
		t.Fatalf("legacy feedback upload accepted: %v", err)
	}
	if err := stats.TryUploadDaily(true, true, false); err != nil {
		t.Fatalf("disabled legacy upload performed work: %v", err)
	}
	on := true
	settings := defaultSettings()
	applySettingsBatchUpdate(&settings, SettingsBatchUpdate{StatsShare: &on, CrashReportAutoSubmit: &on})
	if settings.StatsShare || settings.CrashReportAutoSubmit {
		t.Fatal("batch update restored upstream reporting")
	}
}
