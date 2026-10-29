//go:build !nogui

package gui

import "github.com/wailsapp/wails/v3/pkg/application"

// A macOS window hidden with orderOut keeps its WKWebView's page visible, so
// tell the page before hiding it and tell it again when it is shown.
const (
	pageHiddenJS = "window.magpieWindowVisibility?.(false)"
	pageShownJS  = "window.magpieWindowVisibility?.(true)"
)

func hideWebview(w *application.WebviewWindow) {
	if w == nil {
		return
	}
	w.ExecJS(pageHiddenJS)
	w.Hide()
}

func showWebview(w *application.WebviewWindow) *application.WebviewWindow {
	if w == nil {
		return nil
	}
	w.Show()
	w.ExecJS(pageShownJS)
	return w
}
