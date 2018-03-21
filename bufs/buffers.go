package bufs

import (
	. "fmt"
	"github.com/gotk3/gotk3/gdk"
	"log"
	. "math"
	//"time"
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
	maxR                                 = 0.0
	pixmap                       *gdk.Pixbuf
	tm, rm, bm, lm               []int // margins
)

func InitBuffers() {
	GetOrig()
	DoPixels()
	UpdateViewPixels()
	VectorCenterX, VectorCenterY = View.GetWidth()/2, View.GetHeight()/2
	Println("InitBuffers()")
}

func GetOrig() {
	var err error
	Orig, err = gdk.PixbufNewFromFile("./img.png")
	if err != nil {
		log.Fatal("Unable to create Pixbuf: ", err)
	}
}

func GetTransform(width int, height int) {
	Transform, _ = gdk.PixbufNew(gdk.COLORSPACE_RGB, true, 8, width, height)
}

func UpdateViewPixels() {
	View, _ = gdk.PixbufCopy(Transform)
}

func getNewRadius(x, y float64) float64 {
	r := Hypot(float64(x), float64(y))
	newR := maxR * Tan(r/maxR)
	return newR
}

func DoPixels() {
	orig_bytes := Orig.GetPixels()
	rowstride := Orig.GetRowstride()
	width, height := Orig.GetWidth(), Orig.GetHeight()
	halfW, halfH := width/2, height/2
	maxR = Hypot(float64(halfW), float64(halfH))
	maxNewR := getNewRadius(float64(halfW), float64(halfH))
	maxAngle := Atan2(float64(halfH), float64(halfW))
	maxW, maxH := int(maxNewR*Cos(maxAngle))*2, int(maxNewR*Sin(maxAngle))*2
	GetTransform(maxW, maxH)
	halfTW, halfTH := float64(maxW/2), float64(maxH/2)
	transform_bytes := Transform.GetPixels()
	tRowstride := Transform.GetRowstride()

	pixmap := make([]byte, maxW*maxH)
	tm, bm = make([]int, maxW), make([]int, maxW)
	rm, lm = make([]int, maxH), make([]int, maxH)
	for i, _ := range tm {
		tm[i], bm[i] = -1, -1
	}
	for i, _ := range lm {
		lm[i], rm[i] = -1, -1
	}
	heightm1 := height - 1
	widthm1 := width - 1
	for y := 0; y < height; y++ {
		row_start := y * rowstride
		ry := y - halfH
		for x := 0; x < width; x++ {
			rx := x - halfW
			newR := getNewRadius(float64(rx), float64(ry))
			angle := Atan2(float64(ry), float64(rx))
			newX := int(newR*Cos(angle) + halfTW)
			newY := int(newR*Sin(angle) + halfTH)

			boffset := row_start + x*4
			doffset := newY*tRowstride + newX*4
			map_offset := newY*maxW + newX

			if y == 0 {
				tm[newX] = newY
			} else if y == heightm1 {
				bm[newX] = newY
			}
			if x == 0 {
				lm[newY] = newX
			} else if x == widthm1 {
				rm[newY] = newX
			}

			pixmap[map_offset] = 1
			transform_bytes[doffset] = orig_bytes[boffset]
			transform_bytes[doffset+1] = orig_bytes[boffset+1]
			transform_bytes[doffset+2] = orig_bytes[boffset+2]
			transform_bytes[doffset+3] = orig_bytes[boffset+3]
		}
	}
	InterpolateMargins() // tm, rm, bm, lm
	Interpolate(pixmap, maxW, maxH)
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
			if pixmap[map_pos+x] != 1 && y > tm[x] && y < bm[x] && x > lm[y] && x < rm[y] {
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
		pixels1 [4]rgba
		pixels2 [4]rgba
		count1, count2  byte
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
	tm[0] = 0
	rm[0] = 0
	bm[0] = 0
	lm[0] = 0
	currT := tm[0]
	currR := rm[0]
	currB := bm[0]
	currL := lm[0]
	for i := 0; i < len(tm); i++ {
		if tm[i] < 0 {
			tm[i] = currT
		} else {
			currT = tm[i]
		}
		if bm[i] < 0 {
			bm[i] = currB
		} else {
			currB = bm[i]
		}
	}
	for i := 0; i < len(rm); i++ {
		if rm[i] < 0 {
			rm[i] = currR
		} else {
			currR = rm[i]
		}
		if lm[i] < 0 {
			lm[i] = currL
		} else {
			currL = lm[i]
		}
	}
}
