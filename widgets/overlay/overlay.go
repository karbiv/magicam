package overlay

import (
	//. "fmt"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/karbiv/magicam/app"
	"github.com/karbiv/magicam/bufs"
	"github.com/karbiv/magicam/widgets/drawpix"
)

var (
	overlayDragged         = false
	prevX, prevY   float64 = 0, 0
)

// "size-allocate"
func SizeAllocateHandler(widget *gtk.Overlay, ptr uintptr) {
	rect := gdk.WrapRectangle(ptr)
	rw, rh := rect.GetWidth(), rect.GetHeight()
	bufs.UpdateViewPixels()
	transformRatio := float64(bufs.Transform.GetWidth()) / float64(bufs.Transform.GetHeight())
	rectRatio := float64(rw) / float64(rh)
	if rectRatio > transformRatio {
		bufs.View, _ = bufs.View.ScaleSimple(
			int(float64(rh)*transformRatio), rh, gdk.INTERP_BILINEAR)
	} else {
		bufs.View, _ = bufs.View.ScaleSimple(
			rw, int(float64(rw)*transformRatio), gdk.INTERP_BILINEAR)
	}
	bufs.CurrViewX, bufs.CurrViewY = 0, 0
	app.Overlay.QueueDraw()
}

// "button-press-event"
func ButtonPressHandler(widget *gtk.Overlay, ev *gdk.Event) bool {
	overlayDragged = true
	evm := gdk.EventMotionNewFromEvent(ev)
	prevX, prevY = evm.MotionVal()
	return false
}

// "button-release-event"
func ButtonReleaseHandler(widget *gtk.Overlay, ev *gdk.Event) bool {
	overlayDragged = false
	return false
}

// "motion-notify-element"
func PointerMotionHandler(widget *gtk.Overlay, ev *gdk.Event) bool {
	if overlayDragged {
		evm := gdk.EventMotionNewFromEvent(ev)
		clickedX, clickedY := evm.MotionVal()
		drawpix.DeltaX = int(clickedX - prevX)
		drawpix.DeltaY = int(clickedY - prevY)
		prevX, prevY = clickedX, clickedY
		widget.QueueDraw()
	}
	return false
}
