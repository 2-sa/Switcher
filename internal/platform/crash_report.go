package platform

import (
	"TcNo-Acc-Switcher/internal/api"
	"TcNo-Acc-Switcher/internal/crashlog"
)

// HasPendingCrashReport reports whether a local crash dump is waiting for user action.
func (*PlatformService) HasPendingCrashReport() (bool, error) {
	return crashlog.HasPending(), nil
}

// SubmitPendingCrashReport refuses uploads from legacy native clients.
func (*PlatformService) SubmitPendingCrashReport() (bool, error) {
	return false, api.ErrReportingDisabled
}

// DiscardPendingCrashReport deletes a pending crash dump without uploading it.
func (*PlatformService) DiscardPendingCrashReport() error {
	return crashlog.DiscardPending()
}

func (p *PlatformService) GetCrashReportAutoSubmit() (bool, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	exeDir, err := ResolveExeDir()
	if err != nil {
		return false, err
	}
	s, err := loadSettings(exeDir)
	if err != nil {
		return false, err
	}
	return s.CrashReportAutoSubmit, nil
}

func (p *PlatformService) SetCrashReportAutoSubmit(enabled bool) error {
	if enabled {
		return api.ErrReportingDisabled
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	exeDir, err := ResolveExeDir()
	if err != nil {
		return err
	}
	s, err := loadSettings(exeDir)
	if err != nil {
		return err
	}
	s.CrashReportAutoSubmit = enabled
	return saveSettingsAtomic(exeDir, s)
}
