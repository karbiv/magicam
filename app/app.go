package app

import (
	"github.com/gotk3/gotk3/gtk"
	//"github.com/gotk3/gotk3/gdk"
)

var (
	Application *gtk.Application
	Builder     *gtk.Builder
	Window      *gtk.Window
	Overlay     *gtk.Overlay
	Pixels      *gtk.DrawingArea
	Vector      *gtk.DrawingArea
	Accels      *gtk.AccelGroup
)
