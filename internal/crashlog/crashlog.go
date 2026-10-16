package crashlog

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"time"

	buildinfo "TcNo-Acc-Switcher/build"
	"TcNo-Acc-Switcher/internal/fsutil"
	"TcNo-Acc-Switcher/internal/logredact"
	"TcNo-Acc-Switcher/internal/logsanitize"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const (
	crashDumpFile  = "CrashDump.json"
	toastEventName = "toast"
)

type CrashDump struct {
	Stack     string `json:"stack"`
	Error     string `json:"error"`
	Version   string `json:"version"`
	OS        string `json:"os"`
	OSInfo    string `json:"osInfo"`
	Timestamp string `json:"timestamp"`
	UUID      string `json:"uuid"`
	Log       string `json:"log"`
}

func exeDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exe), nil
}

var (
	crashDumpDirResolver = exeDir
	osExit               = os.Exit // swapped in tests to verify CaptureFatal behavior
)

func crashDumpPath() (string, error) {
	dir, err := crashDumpDirResolver()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, crashDumpFile), nil
}

func writeCrashDump(dump CrashDump) error {
	path, err := crashDumpPath()
	if err != nil {
		return err
	}
	dump.Stack = logredact.RedactText(dump.Stack)
	dump.Error = logredact.RedactText(dump.Error)
	dump.OSInfo = logredact.RedactText(dump.OSInfo)
	dump.Log = logredact.RedactText(dump.Log)
	payload, err := json.MarshalIndent(dump, "", "  ")
	if err != nil {
		return err
	}
	return fsutil.WriteFileAtomic(path, payload, 0o644)
}

// Capture recovers from a panic, logs it, writes CrashDump.json, and returns.
// Use this for background goroutines where the process should stay alive.
func Capture() {
	r := recover()
	if r == nil {
		return
	}

	panicText := captureAndWrite(r)

	if app := application.Get(); app != nil {
		_ = app.Event.Emit(toastEventName, map[string]any{
			"type":     "error",
			"title":    "Background task failed",
			"message":  fmt.Sprintf("A background task crashed (%s). Restart if the app behaves oddly.", panicText),
			"duration": 6000,
		})
	}
}

func CaptureFatal() {
	r := recover()
	if r == nil {
		return
	}

	captureAndWrite(r)
	osExit(1)
}

func captureAndWrite(r any) string {
	stack := string(debug.Stack())
	panicText := logredact.FormatValue(r)
	slog.Error("panic recovered",
		"error", panicText,
		"stack", stack,
	)

	dump := CrashDump{
		Stack:     stack,
		Error:     panicText,
		Version:   buildinfo.Version(),
		OS:        runtime.GOOS + "/" + runtime.GOARCH,
		OSInfo:    osDisplayString(),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Log:       logsanitize.ActionLogForUpload(),
	}

	if err := writeCrashDump(dump); err != nil {
		slog.Warn("writing crash dump", "err", err)
	}
	return panicText
}

// HasPending reports whether a crash dump from a previous run is waiting locally.
func HasPending() bool {
	path, err := crashDumpPath()
	if err != nil {
		return false
	}
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}

// DiscardPending removes a pending crash dump without submitting it.
func DiscardPending() error {
	path, err := crashDumpPath()
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if err != nil && os.IsNotExist(err) {
		return nil
	}
	return err
}

// SubmitPending is retained for native API compatibility. Crash dumps stay local.
func SubmitPending() bool {
	return false
}
