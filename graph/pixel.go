package graph

// #cgo pkg-config: gtk+-3.0
// #include <gtk/gtk.h>
import "C"
import (
	"github.com/gotk3/gotk3/gtk"
	"unsafe"
	. "fmt"
	//"reflect"
)

var (
	numOfPoints = 7
	scales []*gtk.Scale
)

func PixelGraphMenuItem(graphItem *gtk.MenuItem) {
	Println("PixelGraphMenuItem()")
	graphWin, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	graphWin.SetDefaultSize(500, 300)
	scales = make([]*gtk.Scale, numOfPoints)
	grid, _ := gtk.GridNew()
	
	graphWin.Add(grid)

	for i := 0; i < numOfPoints; i++ {
		scale, _ := gtk.ScaleNewWithRange(gtk.ORIENTATION_VERTICAL, 0, 2.0, 0.05)
		scale.Range.SetValue(1.0) // TODO check buffered/saved values
		scale.Range.SetIncrements(0.01, 0.05)
		scale.SetHExpand(true)
		scale.SetVExpand(true)
		scales[i] = scale
		p := unsafe.Pointer(scale.GObject)
		C.gtk_range_set_inverted((*C.GtkRange)(p), C.gboolean(1))
		
		grid.Add(scale)
	}

	graphWin.ShowAll()
}
