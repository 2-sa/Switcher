package app

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"sync/atomic"
	"testing"
	"time"
)

func TestUpdateRestartDeduplicatesButtonClicksAndReadyReplays(t *testing.T) {
	calls := 0
	host := &updateHost{restart: func() error {
		calls++
		return nil
	}}
	host.requestRestart()
	host.requestRestart()
	if calls != 1 {
		t.Fatalf("started %d update helpers, want one", calls)
	}
}

func TestUpdateQuitTerminatesBlockedProcess(t *testing.T) {
	const sentinel = "TCNO_TEST_UPDATE_SHUTDOWN"
	if os.Getenv(sentinel) == "1" {
		host := &updateHost{
			shutdownTimeout: 20 * time.Millisecond,
			quit:            func() { select {} },
			exit:            func() { os.Exit(0) },
		}
		host.Quit()
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestUpdateQuitTerminatesBlockedProcess$")
	cmd.Env = append(os.Environ(), sentinel+"=1")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("blocked update parent did not exit successfully: %v; %s", err, output)
	}
}

func TestUpdateRestartSpawnFailureAllowsRetry(t *testing.T) {
	calls, failures := 0, 0
	spawnError := errors.New("helper spawn denied")
	host := &updateHost{
		restart: func() error {
			calls++
			if calls == 1 {
				return spawnError
			}
			return nil
		},
		restartError: func(err error) {
			if !errors.Is(err, spawnError) {
				t.Fatalf("reported %v, want helper spawn error", err)
			}
			failures++
		},
	}
	host.requestRestart()
	host.requestRestart()
	host.requestRestart()
	if calls != 2 || failures != 1 {
		t.Fatalf("restart calls = %d, errors = %d; want 2, 1", calls, failures)
	}
}

func TestUpdateQuitExitsWhenNormalShutdownHangs(t *testing.T) {
	blocked := make(chan struct{})
	exited := make(chan struct{}, 1)
	returned := make(chan struct{})
	host := &updateHost{
		shutdownTimeout: 20 * time.Millisecond,
		quit:            func() { <-blocked },
		exit:            func() { exited <- struct{}{} },
	}
	go func() {
		host.Quit()
		close(returned)
	}()
	t.Cleanup(func() { close(blocked); <-returned })
	select {
	case <-exited:
	case <-time.After(2 * time.Second):
		t.Fatal("update helper would time out waiting for the blocked parent")
	}
}

func TestUpdateQuitKeepsDeadlineWhenQuitReturnsBeforeProcessExit(t *testing.T) {
	exited := make(chan struct{}, 1)
	var calls atomic.Int32
	host := &updateHost{
		shutdownTimeout: 20 * time.Millisecond,
		quit:            func() { calls.Add(1) },
		exit:            func() { exited <- struct{}{} },
	}
	host.Quit()
	host.Quit()
	if calls.Load() != 1 {
		t.Fatalf("normal shutdown requested %d times, want once", calls.Load())
	}
	select {
	case <-exited:
	case <-time.After(2 * time.Second):
		t.Fatal("returning from Quit must not cancel the process-exit deadline")
	}
}
