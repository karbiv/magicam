package keys

import (
	//. "fmt"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/karbiv/magicam/app"
	"github.com/karbiv/magicam/widgets/drawpix"
	"github.com/karbiv/magicam/widgets/drawvec"
	"github.com/karbiv/magicam/widgets/win"
)

var (
	// key names are in gdkkeysyms.h
	Equal       = gdk.KeyvalFromName("equal")
	Minus       = gdk.KeyvalFromName("minus")
	Up          = gdk.KeyvalFromName("Up")
	Down        = gdk.KeyvalFromName("Down")
	Left        = gdk.KeyvalFromName("Left")
	Right       = gdk.KeyvalFromName("Right")
	KP_Add      = gdk.KeyvalFromName("KP_Add")
	KP_Subtract = gdk.KeyvalFromName("KP_Subtract")
	F11         = gdk.KeyvalFromName("F11")
	F1          = gdk.KeyvalFromName("F1")
)

func KeyPress(widget *gtk.Window, ev *gdk.Event) bool {
	ek := gdk.EventKeyNewFromEvent(ev)
	//Printf("%x\n", ek.KeyVal())
	switch ek.KeyVal() {
	case Equal, KP_Add:
		drawpix.ScaleUp(widget, ev)
	case Minus, KP_Subtract:
		drawpix.ScaleDown(widget, ev)
	case Up:
		drawpix.DeltaY = -drawpix.KeyVerticalMoveStep
		app.Pixels.QueueDraw()
	case Down:
		drawpix.DeltaY = drawpix.KeyVerticalMoveStep
		app.Pixels.QueueDraw()
	case Left:
		drawpix.DeltaX = -drawpix.KeyHorizontalMoveStep
		app.Pixels.QueueDraw()
	case Right:
		drawpix.DeltaX = drawpix.KeyHorizontalMoveStep
		app.Pixels.QueueDraw()
	case F11:
		win.ToggleFullscreen()
	case F1:
		drawvec.ToggleLegend()
	}
	return false
}

func KeyRelease(widget *gtk.Window, ev *gdk.Event) bool {
	return false
}
