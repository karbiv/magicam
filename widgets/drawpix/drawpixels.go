package drawpix

import (
	. "fmt"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/karbiv/magicam/app"
	"github.com/karbiv/magicam/bufs"
)

type ScaleMode uint

const (
	No ScaleMode = iota
	Up
	Down
)

var (
	DeltaX, DeltaY        int       = 0, 0
	scalePercent          int       = 5
	minVisiblePercent               = 24
	KeyHorizontalMoveStep           = 32
	KeyVerticalMoveStep             = 16
	ToScale               ScaleMode = No
)

// "draw"
func DrawPixbufHandler(widget *gtk.DrawingArea, cr *cairo.Context) bool {
	IfMustScale()
	ncurrX, ncurrY := bufs.CurrViewX+DeltaX, bufs.CurrViewY+DeltaY
	// restrict dragging
	W, H := widget.GetAllocatedWidth(), widget.GetAllocatedHeight()
	reqW, reqH := int(W/100*minVisiblePercent), int(H/100*minVisiblePercent)
	leftX, rightX := reqW-bufs.View.GetWidth(), W-reqW
	topY, botY := reqH-bufs.View.GetHeight(), H-reqH
	if ncurrX > leftX && ncurrX < rightX && ncurrY > topY && ncurrY < botY {
		bufs.CurrViewX, bufs.CurrViewY = ncurrX, ncurrY
		gtk.GdkCairoSetSourcePixBuf(cr, bufs.View, float64(ncurrX), float64(ncurrY))
		bufs.VectorCenterX, bufs.VectorCenterY =
			ncurrX+(bufs.View.GetWidth()/2), ncurrY+(bufs.View.GetHeight()/2)
	} else {
		gtk.GdkCairoSetSourcePixBuf(
			cr, bufs.View, float64(bufs.CurrViewX), float64(bufs.CurrViewY))
	}

	cr.Paint()
	DeltaX, DeltaY = 0, 0
	return false
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
		nw, nh := int(w+w/100*scalePercent), int(h+h/100*scalePercent)
		if float64(nw) < float64(trw)*bufs.MaxViewScale {
			DeltaX, DeltaY = int(-(nw-w)/2), int(-(nh-h)/2)
			bufs.View, _ = bufs.Transform.ScaleSimple(nw, nh, gdk.INTERP_BILINEAR)
			//bufs.View, _ = bufs.Transform.ScaleSimple(nw, nh, gdk.INTERP_NEAREST) // faster
		}
	case Down:
		nw, nh := int(w-w/100*scalePercent), int(h-h/100*scalePercent)
		if float64(nw) > float64(trw)/bufs.MinViewScale {
			DeltaX, DeltaY = int((w-nw)/2), int((h-nh)/2)
			bufs.View, _ = bufs.Transform.ScaleSimple(nw, nh, gdk.INTERP_BILINEAR)
			//bufs.View, _ = bufs.Transform.ScaleSimple(nw, nh, gdk.INTERP_NEAREST) // faster
		}
	}
	ToScale = No
}

// "scroll-event"
func WheelHandler(widget *gtk.DrawingArea, ev *gdk.Event) bool {
	widget.Native()
	scroll := gdk.EventScrollNewFromEvent(ev)
	direct := scroll.Direction()
	Println(direct)
	return false
}
