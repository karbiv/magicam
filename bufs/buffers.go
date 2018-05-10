package bufs

// #cgo pkg-config: gtk+-3.0
// #include <gtk/gtk.h>
import "C"
import (
	"github.com/gotk3/gotk3/gdk"
	"log"
	. "math"
	//"time"
	//"os"
	"unsafe"
	//. "fmt"
)

var (
	Orig      *gdk.Pixbuf
	Transform *gdk.Pixbuf
	View      *gdk.Pixbuf
)

var (
	VectorCenterX, VectorCenterY         = 0, 0
	CurrViewX, CurrViewY         int     = 0, 0
	MaxViewScale                 float64 = 1.5
	MinViewScale                 float64 = 5
	coeff                                = 1.03
	MaxR                                 = 0.0
	MaxNewR                              = 0.0
	MaxRadiusAngle                       = 0.0
	pixmap                       *gdk.Pixbuf
	topMargin, rightMargin       []int // margins
	bottomMargin, leftMargin     []int // margins

	origBytes     []byte
	origNumChans  int
	rowstride     int
	width, height int
	halfW, halfH  int
	TW, TH        int
	Graph         []float64

	// enough length to contain all radius coeffs
	RadiusGraph          [2000]float64
	CurrentViewScale     float64 = 1.0 // scale coeff
)

func InitBuffers() {
	GetOrig()
	initFrameData()
	DoPixels()
	UpdateViewPixels()
	VectorCenterX, VectorCenterY = View.GetWidth()/2, View.GetHeight()/2
}

func GetOrig() {
	var err error
	Orig, err = gdk.PixbufNewFromFile("./img.png")
	if err != nil {
		log.Fatal("Unable to create Pixbuf: ", err)
	}
	origNumChans = Orig.GetNChannels()
}

func UpdateViewPixels() {
	// Restore View buffer scale after Transform buffer has been recreated
	w, h := Transform.GetWidth(), Transform.GetHeight()
	nw, nh := int(float64(w)*CurrentViewScale), int(float64(h)*CurrentViewScale)
	View, _ = Transform.ScaleSimple(nw, nh, gdk.INTERP_BILINEAR)
	//bufs.View, _ = bufs.Transform.ScaleSimple(nw, nh, gdk.INTERP_NEAREST) // faster
	
	//View, _ = gdk.PixbufCopy(Transform)
}

func getNewRadius(x, y float64) float64 {
	r := Hypot(float64(x), float64(y))
	//newR := MaxR * Tan(r/MaxR) // basic demo function
	newR := r * Graph[int(r)]

	return newR
}

func initFrameData() {
	origBytes = Orig.GetPixels()
	rowstride = Orig.GetRowstride()
	width, height = Orig.GetWidth(), Orig.GetHeight()
	halfW, halfH = width/2, height/2
	MaxR = Hypot(float64(halfW), float64(halfH))
	Graph = make([]float64, int(MaxR)+1)
	for i, _ := range Graph {
		Graph[i] = 1.0
	}
	MaxNewR = getNewRadius(float64(halfW), float64(halfH))
	SetTransformBufSize(MaxNewR)
}

func SetTransformBufSize(maxNewR float64) {
	MaxNewR = maxNewR
	MaxRadiusAngle = Atan2(float64(halfH), float64(halfW))
	TW, TH = int(maxNewR*Cos(MaxRadiusAngle))*2, int(maxNewR*Sin(MaxRadiusAngle))*2
	Transform, _ = gdk.PixbufNew(gdk.COLORSPACE_RGB, true, 8, TW, TH)
}

func DoPixels() {
	halfTW, halfTH := float64(TW/2), float64(TH/2)
	transformBytes := Transform.GetPixels()
	tRowstride := Transform.GetRowstride()

	// Clear Transform PixBuf
	p := unsafe.Pointer(Transform.Native())
	C.gdk_pixbuf_fill((*C.GdkPixbuf)(p), C.guint32(0))

	pixmap := make([]byte, TW*TH)
	topMargin, bottomMargin = make([]int, TW), make([]int, TW)
	rightMargin, leftMargin = make([]int, TH), make([]int, TH)
	for i, _ := range topMargin {
		topMargin[i], bottomMargin[i] = -1, -1
	}
	for i, _ := range leftMargin {
		leftMargin[i], rightMargin[i] = -1, -1
	}

	heightm1 := height - 1
	widthm1 := width - 1

	// loops over all pixels of original frame
	// set adjusted pixels in the Transform Pixbuf
	for y := 0; y < height; y++ {
		row_start := y * rowstride
		ry := y - halfH
		for x := 0; x < width; x++ {
			rx := x - halfW
			newR := getNewRadius(float64(rx), float64(ry))
			angle := Atan2(float64(ry), float64(rx))
			newX := int(newR*Cos(angle) + halfTW)
			newY := int(newR*Sin(angle) + halfTH)

			mapOffset := newY*TW + newX

			// margins pixels as boundaries, in 4 arrays
			if y == 0 {
				topMargin[newX] = newY
			} else if y == heightm1 {
				bottomMargin[newX] = newY
			}
			if x == 0 {
				leftMargin[newY] = newX
			} else if x == widthm1 {
				rightMargin[newY] = newX
			}

			pixmap[mapOffset] = 1

			offset := row_start + x*origNumChans
			tOffset := newY*tRowstride + newX*4

			// pixel channels in Orig pixbuf
			for channel := 0; channel < origNumChans; channel++ {
				transformBytes[tOffset+channel] = origBytes[offset+channel]
			}
			// Alpha channel, 8 bit sample
			transformBytes[tOffset+3] = 255
		}
	}
	InterpolateMargins() // topMargin, rightMargin, bottomMargin, leftMargin
	Interpolate(pixmap, TW, TH)
}

type rgba struct {
	r, g, b, a byte
}

func Interpolate(pixmap []byte, w int, h int) {
	rowstride := Transform.GetRowstride()
	chans := Transform.GetPixels()
	for y := 0; y < h; y++ {
		map_pos := y * w
		row_pos := y * rowstride
		for x := 0; x < w; x++ {
			if pixmap[map_pos+x] != 1 &&
				y > topMargin[x] && y < bottomMargin[x] && x > leftMargin[y] && x < rightMargin[y] {
				dest := row_pos + x*4
				ip := getInterpolatedPixel(dest, chans, w, rowstride, pixmap, map_pos+x)
				chans[dest] = ip.r
				chans[dest+1] = ip.g
				chans[dest+2] = ip.b
				chans[dest+3] = ip.a
				pixmap[map_pos+x] = 1
			}
		}
	}
}

var (
	//   212
	//   1X1
	//   212
	ring1Wght = 0.60
	ring2Wght = 0.40
)

func getInterpolatedPixel(dest int, chans []byte, mapWidth int, rowstride int,
	pixmap []byte, mapIdx int) (ip rgba) {

	var (
		pixels1        [4]rgba
		pixels2        [4]rgba
		count1, count2 byte
	)

	// ring 1
	if pixmap[mapIdx-1] == 1 { // west
		pixels1[count1] = rgba{chans[dest-4], chans[dest-3], chans[dest-2], chans[dest-1]}
		count1++
	}
	if pixmap[mapIdx+1] == 1 { //east
		pixels1[count1] = rgba{chans[dest+4], chans[dest+5], chans[dest+6], chans[dest+7]}
		count1++
	}
	if pixmap[mapIdx-mapWidth] == 1 { //north
		idx := dest - rowstride
		pixels1[count1] = rgba{chans[idx], chans[idx+1], chans[idx+2], chans[idx+3]}
		count1++
	}
	if pixmap[mapIdx+mapWidth] == 1 { // south
		idx := dest + rowstride
		pixels1[count1] = rgba{chans[idx], chans[idx+1], chans[idx+2], chans[idx+3]}
		count1++
	}
	// ring 2
	if pixmap[mapIdx-mapWidth-1] == 1 { //north-west
		idx := dest - rowstride - 4
		pixels2[count2] = rgba{chans[idx], chans[idx+1], chans[idx+2], chans[idx+3]}
		count2++
	}
	if pixmap[mapIdx-mapWidth+1] == 1 { //north-east
		idx := dest - rowstride + 4
		pixels2[count2] = rgba{chans[idx], chans[idx+1], chans[idx+2], chans[idx+3]}
		count2++
	}
	if pixmap[mapIdx+mapWidth-1] == 1 { //south-west
		idx := dest + rowstride - 4
		pixels2[count2] = rgba{chans[idx], chans[idx+1], chans[idx+2], chans[idx+3]}
		count2++
	}
	if pixmap[mapIdx+mapWidth+1] == 1 { //south-east
		idx := dest + rowstride + 4
		pixels2[count2] = rgba{chans[idx], chans[idx+1], chans[idx+2], chans[idx+3]}
		count2++
	}

	var (
		r1, g1, b1, r2, g2, b2 float64
	)

	for i := byte(0); i < count1; i++ {
		r1 += float64(pixels1[i].r)
		g1 += float64(pixels1[i].g)
		b1 += float64(pixels1[i].b)
	}
	for i := byte(0); i < count2; i++ {
		r2 += float64(pixels2[i].r)
		g2 += float64(pixels2[i].g)
		b2 += float64(pixels2[i].b)
	}

	countF1, countF2 := float64(count1), float64(count2)

	r1 = (r1 / countF1) * ring1Wght
	g1 = (g1 / countF1) * ring1Wght
	b1 = (b1 / countF1) * ring1Wght
	r2 = (r2 / countF2) * ring2Wght
	g2 = (g2 / countF2) * ring2Wght
	b2 = (b2 / countF2) * ring2Wght

	// rounding
	ip.r = byte(r1 + r2 + 0.5)
	ip.g = byte(g1 + g2 + 0.5)
	ip.b = byte(b1 + b2 + 0.5)
	ip.a = 255

	return ip
}

func InterpolateMargins() {
	topMargin[0] = 0
	rightMargin[0] = 0
	bottomMargin[0] = 0
	leftMargin[0] = 0
	currT := topMargin[0]
	currR := rightMargin[0]
	currB := bottomMargin[0]
	currL := leftMargin[0]
	for i := 0; i < len(topMargin); i++ {
		if topMargin[i] < 0 {
			topMargin[i] = currT
		} else {
			currT = topMargin[i]
		}
		if bottomMargin[i] < 0 {
			bottomMargin[i] = currB
		} else {
			currB = bottomMargin[i]
		}
	}
	for i := 0; i < len(rightMargin); i++ {
		if rightMargin[i] < 0 {
			rightMargin[i] = currR
		} else {
			currR = rightMargin[i]
		}
		if leftMargin[i] < 0 {
			leftMargin[i] = currL
		} else {
			currL = leftMargin[i]
		}
	}
}
