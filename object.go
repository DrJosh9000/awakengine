package awakengine

import "github.com/hajimehoshi/ebiten"

// Semiobject is some import stuff for logical grouping.
type Semiobject interface {
	InWorld() bool // true if the object exists in world-coordinates, false if screen coordinates
	Retire() bool  // true if the object will never draw again and can be removed from the draw list
	Visible() bool
	Z() int
}

// Object is everything, everything is an object.
type Object interface {
	Semiobject
	ImageKey() string
	Dst() (x0, y0, x1, y1 int)
	Src() (x0, y0, x1, y1 int) // relative to the image referred to by ImageKey()
}

type drawList []Object

// Sorting by Z.
func (d drawList) Len() int           { return len(d) }
func (d drawList) Less(i, j int) bool { return d[i].Z() < d[j].Z() }
func (d drawList) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }

// cull removes invisible objects. Visibility is given
func (d drawList) cull() drawList {
	l := make(drawList, 0, len(d))
	for i, o := range d {
		if !o.Visible() {
			//log.Printf("culling %#v because visible=false", o)
			continue
		}
		if x0, y0, x1, y1 := d.Dst(i); x1 <= 0 || y1 <= 0 || x0 > camSize.X || y0 > camSize.Y {
			//log.Printf("culling %#v because off-screen", o)
			continue
		}
		l = append(l, o)
	}
	return l
}

// gc removes retired objects.
func (d drawList) gc() drawList {
	l := make(drawList, 0, len(d))
	for _, o := range d {
		if o.Retire() {
			continue
		}
		l = append(l, o)
	}
	//log.Printf("drawList gc'ed %d objects", len(d)-len(l))
	return l
}

func (d drawList) Dst(i int) (x0, y0, x1, y1 int) {
	x0, y0, x1, y1 = d[i].Dst()
	if !d[i].InWorld() {
		return
	}
	x0 -= camPos.X
	y0 -= camPos.Y
	x1 -= camPos.X
	y1 -= camPos.Y
	return
}

func (d drawList) Src(i int) (x0, y0, x1, y1 int) {
	x0, y0, x1, y1 = d[i].Src()
	o := compositeOffset[d[i].ImageKey()]
	x0 += o.X
	y0 += o.Y
	x1 += o.X
	y1 += o.Y
	return
}

func (d drawList) draw(screen *ebiten.Image) error {
	return screen.DrawImage(composite, &ebiten.DrawImageOptions{
		ImageParts: d,
	})
}
