package drawpix

import (
	//"math"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/karbiv/magicam/app"
	"github.com/karbiv/magicam/bufs"
	//"unsafe"
	//. "fmt"
)

type ScaleMode uint

const (
	No ScaleMode = iota
	Up
	Down
)

var (
	DeltaX, DeltaY        int       = 0, 0
	defaultScalePercent   int       = 5
	minVisiblePercent               = 24
	KeyHorizontalMoveStep           = 32
	KeyVerticalMoveStep             = 16
	ToScale               ScaleMode = No
)

// "draw"
func DrawPixbufHandler(widget *gtk.DrawingArea, cairoCtx *cairo.Context) bool {
	IfMustScale()
	Move(widget, cairoCtx)
	cairoCtx.Paint()
	return false
}

func Move(widget *gtk.DrawingArea, cairoCtx *cairo.Context) {
	ncurrX, ncurrY := bufs.CurrViewX+DeltaX, bufs.CurrViewY+DeltaY
	// restrict dragging
	W, H := widget.GetAllocatedWidth(), widget.GetAllocatedHeight()
	reqW, reqH := int(W/100*minVisiblePercent), int(H/100*minVisiblePercent)
	leftX, rightX := reqW-bufs.View.GetWidth(), W-reqW
	topY, botY := reqH-bufs.View.GetHeight(), H-reqH
	if ncurrX > leftX && ncurrX < rightX && ncurrY > topY && ncurrY < botY {
		bufs.CurrViewX, bufs.CurrViewY = ncurrX, ncurrY
		gtk.GdkCairoSetSourcePixBuf(cairoCtx, bufs.View, float64(ncurrX), float64(ncurrY))
		bufs.VectorCenterX, bufs.VectorCenterY =
			ncurrX+(bufs.View.GetWidth()/2), ncurrY+(bufs.View.GetHeight()/2)
	} else {
		gtk.GdkCairoSetSourcePixBuf(
			cairoCtx, bufs.View, float64(bufs.CurrViewX), float64(bufs.CurrViewY))
	}

	DeltaX, DeltaY = 0, 0
}

func ScaleUp(widget *gtk.Window, ev *gdk.Event) bool {
	ToScale = Up
	app.Pixels.QueueDraw()
	return false
}

func ScaleDown(widget *gtk.Window, ev *gdk.Event) bool {
	ToScale = Down
	app.Pixels.QueueDraw()
	return false
}

func IfMustScale() {
	if ToScale == No {
		return
	}
	w, h := bufs.View.GetWidth(), bufs.View.GetHeight()
	trw := bufs.Transform.GetWidth()
	switch ToScale {
	case Up:
		nw, nh := int(w+w/100*defaultScalePercent), int(h+h/100*defaultScalePercent)
		if float64(nw) < float64(trw)*bufs.MaxViewScale {
			DeltaX, DeltaY = int(-(nw-w)/2), int(-(nh-h)/2)
			bufs.View, _ = bufs.Transform.ScaleSimple(nw, nh, gdk.INTERP_BILINEAR)
			//bufs.View, _ = bufs.Transform.ScaleSimple(nw, nh, gdk.INTERP_NEAREST) // faster
			bufs.CurrentViewScale = float64(w)/float64(bufs.Transform.GetWidth())
		}
	case Down:
		nw, nh := int(w-w/100*defaultScalePercent), int(h-h/100*defaultScalePercent)
		if float64(nw) > float64(trw)/bufs.MinViewScale {
			DeltaX, DeltaY = int((w-nw)/2), int((h-nh)/2)
			bufs.View, _ = bufs.Transform.ScaleSimple(nw, nh, gdk.INTERP_BILINEAR)
			//bufs.View, _ = bufs.Transform.ScaleSimple(nw, nh, gdk.INTERP_NEAREST) // faster
			bufs.CurrentViewScale = float64(w)/float64(bufs.Transform.GetWidth())
		}
	}
	ToScale = No
}

// "scroll-event"
func WheelHandler(widget *gtk.Window, ev *gdk.Event) bool {
	scroll := gdk.EventScrollNewFromEvent(ev)
	direct := scroll.Direction()
	if direct == 1 {
		ScaleDown(widget, ev)
	} else if direct == 0 {
		ScaleUp(widget, ev)
	}
	return false
}
