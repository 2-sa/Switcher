package crashlog

import (
	"bytes"
	"os"
	"testing"
	"time"
)

func TestDisabledSubmissionPreservesLocalCrashDump(t *testing.T) {
	dir := t.TempDir()
	previous := crashDumpDirResolver
	crashDumpDirResolver = func() (string, error) { return dir, nil }
	t.Cleanup(func() { crashDumpDirResolver = previous })
	if err := writeCrashDump(CrashDump{Error: "local diagnostic"}); err != nil {
		t.Fatal(err)
	}
	path, err := crashDumpPath()
	if err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if SubmitPending() {
		t.Fatal("crash report was submitted")
	}
	after, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatalf("local diagnostic was removed or changed: %v", err)
	}
}

func TestCapture_DoesNotExit(t *testing.T) {
	dir := t.TempDir()
	origResolver := crashDumpDirResolver
	crashDumpDirResolver = func() (string, error) { return dir, nil }
	t.Cleanup(func() { crashDumpDirResolver = origResolver })

	origExit := osExit
	exited := make(chan int, 1)
	osExit = func(code int) { exited <- code }
	t.Cleanup(func() { osExit = origExit })

	done := make(chan struct{})
	go func() {
		defer close(done)
		defer Capture()
		panic("non-fatal background panic")
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Capture() goroutine did not return")
	}

	select {
	case code := <-exited:
		t.Fatalf("Capture() called os.Exit(%d)", code)
	default:
	}

	if !HasPending() {
		t.Fatal("expected crash dump to be written")
	}
}

func TestCaptureFatal_Exits(t *testing.T) {
	dir := t.TempDir()
	origResolver := crashDumpDirResolver
	crashDumpDirResolver = func() (string, error) { return dir, nil }
	t.Cleanup(func() { crashDumpDirResolver = origResolver })

	origExit := osExit
	exited := make(chan int, 1)
	osExit = func(code int) { exited <- code }
	t.Cleanup(func() { osExit = origExit })

	done := make(chan struct{})
	go func() {
		defer close(done)
		defer CaptureFatal()
		panic("fatal main panic")
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("CaptureFatal() goroutine did not return")
	}

	select {
	case code := <-exited:
		if code != 1 {
			t.Fatalf("CaptureFatal() called os.Exit(%d), want 1", code)
		}
	default:
		t.Fatal("CaptureFatal() did not call os.Exit")
	}

	if !HasPending() {
		t.Fatal("expected crash dump to be written")
	}
}

func TestHasPendingAndDiscardPending(t *testing.T) {
	dir := t.TempDir()
	orig := crashDumpDirResolver
	crashDumpDirResolver = func() (string, error) { return dir, nil }
	t.Cleanup(func() { crashDumpDirResolver = orig })

	if HasPending() {
		t.Fatal("expected no pending crash dump")
	}

	dumpPath, err := crashDumpPath()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dumpPath, []byte(`{"error":"test"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if !HasPending() {
		t.Fatal("expected pending crash dump")
	}
	if err := DiscardPending(); err != nil {
		t.Fatal(err)
	}
	if HasPending() {
		t.Fatal("expected crash dump to be discarded")
	}
}
