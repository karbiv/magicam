package main

import (
	//. "fmt"
	"log"
	"os"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"karbiv/magicam/app"
	"karbiv/magicam/builder"
	"karbiv/magicam/bufs"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	application, err := gtk.ApplicationNew("ak.magicam", glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Fatal("Application error: ", err)
	}

	application.Connect("activate", activateHandler)
	status := application.Run(os.Args)
	app.Application = application	
	gtk.Main()
	os.Exit(status)
}

func activateHandler(application *gtk.Application) {
	builder.InitApp()
	bufs.InitBuffers()
	app.Window.ShowAll()
	app.Window.Maximize()
}
