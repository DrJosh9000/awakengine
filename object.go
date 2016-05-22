package awakengine

import "github.com/hajimehoshi/ebiten"

// Object is everything, everything is an object.
type Object interface {
	ImageKey() string
	Dst() (x0, y0, x1, y1 int)
	InWorld() bool
	Src() (x0, y0, x1, y1 int) // relative to the image referred to by ImageKey()
	Visible() bool
	Z() int
}

type drawList []Object

// Sorting by Z.
func (d drawList) Len() int           { return len(d) }
func (d drawList) Less(i, j int) bool { return d[i].Z() < d[j].Z() }
func (d drawList) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }

// cull removes invisible objects. Visibility is given
func (d drawList) cull() drawList {
	l := make(drawList, 0, len(d))
	for _, o := range d {
		if !o.Visible() {
			continue
		}
		if x0, y0, x1, y1 := o.Dst(); x1 <= camPos.X || y1 <= camPos.Y || x0 > camPos.X+camSize.X || y0 > camPos.Y+camSize.Y {
			continue
		}
		l = append(l, o)
	}
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
