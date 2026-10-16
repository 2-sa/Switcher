package app

import (
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/updater"
)

// Wails' swap helper abandons the update if its parent is still alive after
// 30 seconds. Give normal shutdown 20 seconds, only after Restart has spawned
// that helper. Ordinary application exits never use this deadline.
const updateShutdownTimeout = 20 * time.Second

type updateHost struct {
	app             *application.App
	shutdownTimeout time.Duration
	quit            func()
	exit            func()
	restart         func() error
	restartError    func(error)
	quitting        sync.Once
	restarting      atomic.Bool
}

var _ updater.Host = (*updateHost)(nil)

func newUpdateHost(app *application.App) *updateHost {
	return &updateHost{
		app:             app,
		shutdownTimeout: updateShutdownTimeout,
		quit:            app.Quit,
		exit:            func() { os.Exit(0) },
		restart:         func() error { return app.Updater.Restart(app.Context()) },
		restartError: func(err error) {
			app.Logger.Error("update: restart failed; application remains open", "error", err)
		},
	}
}

func (h *updateHost) Emit(name string, data ...any) bool {
	return h.app.Event.Emit(name, data...)
}

func (h *updateHost) OnEvent(name string, callback func(any)) func() {
	return h.app.Event.On(name, func(event *application.CustomEvent) {
		if name == updater.EventUserRestart {
			h.requestRestart()
			return
		}
		callback(event.Data)
	})
}

func (h *updateHost) requestRestart() {
	if !h.restarting.CompareAndSwap(false, true) {
		return
	}
	if err := h.restart(); err != nil {
		h.restarting.Store(false)
		h.restartError(err)
	}
}

func (h *updateHost) OpenWindow(opts updater.WindowOptions) updater.WindowHandle {
	win := h.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:                opts.Title,
		Width:                opts.Width,
		Height:               opts.Height,
		Frameless:            opts.Frameless,
		AlwaysOnTop:          opts.AlwaysOnTop,
		DisableResize:        opts.DisableResize,
		HTML:                 opts.InitialHTML,
		AllowSimpleEventEmit: true,
	})
	win.Show()
	return win.AsUpdaterWindow()
}

func (h *updateHost) Quit() {
	h.quitting.Do(func() {
		time.AfterFunc(h.shutdownTimeout, h.exit)
		if h.app != nil {
			h.app.Logger.Info("update: helper started; requesting shutdown", "exitDeadline", h.shutdownTimeout)
		}
		// Do not stop the timer when Quit returns: Wails can return before
		// Run and its deferred service shutdowns finish. Process exit is the
		// acknowledgement the helper needs, not a return from App.Quit.
		h.quit()
	})
}
