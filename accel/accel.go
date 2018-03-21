package accel

// info about keynames is in magicam/gdkkeysyms.h

import (
	. "fmt"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/gdk"
	"karbiv/magicam/app"
)

func SetAccels() {
	app.Accels.Connect(gdk.KeyvalFromName("f"), 0, 0, fitWindowAccel)
	app.Accels.Connect(gdk.KeyvalFromName("Escape"), 0, 0, quit)
	// app.Accels.Connect(gdk.KeyvalFromName("equal"), 0, 0, drawpix.ScaleUp)
	// app.Accels.Connect(gdk.KeyvalFromName("minus"), 0, 0, drawpix.ScaleDown)
}

func fitWindowAccel() bool {
	Println("fitWindowAccel")
	return true
}

func quit() {
	gtk.MainQuit()
}
