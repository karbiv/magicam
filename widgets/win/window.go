package win

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/karbiv/magicam/app"
	//. "fmt"
)

var (
	isFullscreen                    bool = false
	allocatedWidth, allocatedHeight int
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

// "check-resize" event for the window
func CheckResize(container *gtk.Window) {
	allocatedWidth = container.GetAllocatedWidth()
	allocatedHeight = container.GetAllocatedHeight()
}

func DestroyHandler() {
	gtk.MainQuit()
}
