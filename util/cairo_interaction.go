package util

// #cgo pkg-config: json-glib-1.0
// #include <gdk/gdk.h>
import "C"

import (
	"errors"
	"unsafe"
	
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/cairo"
)


func Dump(obj glib.IObject) (string, error) {
	if obj == nil || obj.toGObject() == nil {
		return "", nil
	}
	data := C.json_gobject_to_data((*C.GObject)(unsafe.Pointer(obj.GObject)), nil)
	if data == nil {
		return "", errors.New("json_gobject_to_data returned nil pointer")
	}
	return C.GoString((*C.char)(data)), nil
}
