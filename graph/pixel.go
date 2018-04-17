package graph

// #cgo pkg-config: gtk+-3.0
// #include <gtk/gtk.h>
/*
int mytest(int i) {
  return i+1;
}
*/
import "C"
import (
	. "fmt"
	"github.com/gotk3/gotk3/gtk"
	"github.com/karbiv/magicam/app"
	"github.com/karbiv/magicam/bufs"
	"strconv"
	"unsafe"
	//"reflect"
	//"time"
)

var (
	numOfPoints = 7
	scales      []*gtk.Scale
)

func PixelGraphMenuItem(graphItem *gtk.MenuItem) {
	Println("PixelGraphMenuItem()")
	graphWin, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	graphWin.SetDefaultSize(500, 300)
	scales = make([]*gtk.Scale, numOfPoints)
	grid, _ := gtk.GridNew()

	graphWin.Add(grid)

	for i := 0; i < numOfPoints; i++ {
		scale, _ := gtk.ScaleNewWithRange(gtk.ORIENTATION_VERTICAL, 0, 3.2, 0.05)
		scale.SetValue(1.0) // TODO check buffered/saved values
		scale.SetIncrements(0.01, 0.05)
		scale.Connect("value-changed", valueChanged)
		scale.SetHExpand(true)
		scale.SetVExpand(true)
		scale.SetName(strconv.Itoa(i))

		scales[i] = scale
		p := unsafe.Pointer(scale.GObject)
		C.gtk_range_set_inverted((*C.GtkRange)(p), C.gboolean(1))

		grid.Add(scale)
	}

	graphWin.ShowAll()
}

func valueChanged(s *gtk.Scale) {
	// compute radius multiplication coefficient graph
	// updates radius coefficients graph

	scaleName, _ := s.GetName() // is index as string type
	scaleIndex, _ := strconv.ParseInt(scaleName, 10, 8)
	// if it's the rightmost scale, then enlarge Transform buffer
	if (int(scaleIndex) + 1) == numOfPoints {
		bufs.SetTransformBufSize(bufs.MaxR * s.GetValue())
	} else { // next scale(to the right) can't be lower
		nextScale := scales[int(scaleIndex)+1]
		currScaleVal := s.GetValue()
		if currScaleVal > nextScale.GetValue() {
			s.SetValue(nextScale.GetValue())
			return
		}
	}

	radiusSectionSize := int(bufs.MaxR / float64(len(scales)))
	graphIndexes := make([]int, numOfPoints+1)
	currGraphIndex := int(bufs.MaxR)
	endIndex := len(graphIndexes)
	for i, _ := range graphIndexes {
		graphIndexes[endIndex-1-i] = currGraphIndex
		currGraphIndex = currGraphIndex - radiusSectionSize
	}
	graphIndexes[0] = 0
	graph := make([]float64, int(bufs.MaxR)+1)
	for i, _ := range graph {
		graph[i] = 1.0
	}
	var leftVal float64 = 1.0

	for i := 0; i < numOfPoints; i++ {
		rightVal := scales[i].GetValue()

		// TODO now implies that endVal is always higher
		if (rightVal - leftVal) != 0 {
			step := (rightVal - leftVal) / float64(graphIndexes[i+1]-graphIndexes[i])
			for j := graphIndexes[i]; j < graphIndexes[i+1]; j++ {
				graph[j] = leftVal
				leftVal = leftVal + step
			}
		} else {
			for j := graphIndexes[i]; j < graphIndexes[i+1]; j++ {
				graph[j] = rightVal
			}
		}

	}

	bufs.Graph = graph
	bufs.DoPixels()
	bufs.UpdateViewPixels()
	app.Pixels.QueueDraw()
}
