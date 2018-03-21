package builder

import (
	. "fmt"
	"log"
	"github.com/gotk3/gotk3/gtk"
	"karbiv/magicam/app"
	"karbiv/magicam/accel"
	"karbiv/magicam/keys"
	"karbiv/magicam/widgets/drawpix"
	"karbiv/magicam/widgets/drawvec"
	"karbiv/magicam/widgets/overlay"
	"karbiv/magicam/widgets/win" 
)

var (
	signalFuncs = make(map[string]interface{})
)

func InitApp() {
	Println("InitApp()")
	initBuilder()
	connectSignalFuncs()
	initWindow()
	initOverlay()
	initDrawPix()
	initDrawVec()
	initAccelGroup()
}

func connectSignalFuncs() {
	signalFuncs["winDestroy"] = win.DestroyHandler
	signalFuncs["drawPixbuf"] = drawpix.DrawPixbufHandler
	signalFuncs["drawVector"] = drawvec.DrawVectorHandler
	signalFuncs["pointerMotion"] = overlay.PointerMotionHandler
	signalFuncs["buttonPress"] = overlay.ButtonPressHandler
	signalFuncs["buttonRelease"] = overlay.ButtonReleaseHandler
	signalFuncs["sizeAllocate"] = overlay.SizeAllocateHandler
	signalFuncs["wheel"] = drawpix.WheelHandler
	signalFuncs["keyPress"] = keys.KeyPress
	signalFuncs["keyRelease"] = keys.KeyRelease
	app.Builder.ConnectSignals(signalFuncs)
}

func initBuilder() {
	builder, _ := gtk.BuilderNew()
	err := builder.AddFromFile("interface.xml")
	if err != nil {
		log.Fatal("Builder error: ", err)
	}
	app.Builder = builder
}

func initWindow() {
	_window, _ := app.Builder.GetObject("win")
	if window, ok := _window.(*gtk.Window); ok {
		app.Window = window
	} else {
		log.Fatal("Error: Unable to create GtkWindow.")
	}
}

func initOverlay() {
	_overlay, _ := app.Builder.GetObject("overlay")
	if overlay, ok := _overlay.(*gtk.Overlay); ok {
		app.Overlay = overlay
	} else {
		log.Fatal("Unable to fetch GtkOverlay from interface.xml")
	}
}

func initDrawPix() {
	_pixbuf, _ := app.Builder.GetObject("pixbuf")
	if pixbuf, ok := _pixbuf.(*gtk.DrawingArea); ok {
		app.Pixels = pixbuf
	} else {
		log.Fatal("Unable to fetch GtkDrawingArea from interface.xml")
	}
}

func initDrawVec() {
	_vector, _ := app.Builder.GetObject("vector")
	if vector, ok := _vector.(*gtk.DrawingArea); ok {
		app.Vector = vector
	} else {
		log.Fatal("Unable to fetch GtkDrawingArea from interface.xml")
	}
}

func initAccelGroup() {
	_accels, _ := app.Builder.GetObject("accelGroup1")
	if accels, ok := _accels.(*gtk.AccelGroup); ok {
		app.Accels = accels
		accel.SetAccels()
	} else {
		log.Fatal("Unable to create AccelGroup.")
	}
}
