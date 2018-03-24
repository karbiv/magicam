package win

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/karbiv/magicam/app"
)

var (
	isFullscreen bool = false
)

func ToggleFullscreen() {
	if isFullscreen {
		app.Window.Unfullscreen()
		isFullscreen = false
	} else {
		app.Window.Fullscreen()
		isFullscreen = true
	}
}

func DestroyHandler() {
	gtk.MainQuit()
}
