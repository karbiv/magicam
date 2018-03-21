package drawvec

import (
	//. "fmt"
	"math"
	"github.com/gotk3/gotk3/cairo"
	//"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"karbiv/magicam/app"
	"karbiv/magicam/bufs"
)

var (
	Legend           bool    = true
	ViewScaled               = 0.0
	ResultFrameW     float64 = 1920
	ResultFrameH     float64 = 1080
	centerX, centerY float64 = 0, 0
)

func DrawVectorHandler(widget *gtk.DrawingArea, c *cairo.Context) bool {
	if Legend {
		centerX, centerY = float64(bufs.VectorCenterX), float64(bufs.VectorCenterY)
		DrawLegend(c)
	}
	return false
}

func DrawLegend(c *cairo.Context) {
	ViewScaled = float64(bufs.View.GetWidth()) / float64(bufs.Transform.GetWidth())
	FrameBorder(c)
	CenterDot(c)
}

func FrameBorder(c *cairo.Context) {
	w := ResultFrameW * ViewScaled
	h := ResultFrameH * ViewScaled
	c.SetLineWidth(1)
	c.SetSourceRGB(1.0, 0.0, 0.0)
	c.Rectangle(centerX-w/2-1, centerY-h/2-1, w+2, h+2)
	c.Stroke()
}

func CenterDot(c *cairo.Context) {
	c.SetSourceRGB(1.0, 0.0, 0.0)
	c.SetLineWidth(5)
	c.Arc(centerX, centerY, 1, 0, math.Pi*2)
	c.Stroke()
}

func ToggleLegend() {
	Legend = !Legend
	app.Overlay.QueueDraw()
}
