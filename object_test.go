package awakengine

import "testing"

type rect struct{ x0, y0, x1, y1 int }

func (r rect) c() (x0, y0, x1, y1 int) { return r.x0, r.y0, r.x1, r.y1 }

type fakeObject struct {
	inWorld, retire, visible bool
	z                        int
	imageKey                 string
	dst, src                 rect
}

func (f *fakeObject) InWorld() bool             { return f.inWorld }
func (f *fakeObject) Retire() bool              { return f.retire }
func (f *fakeObject) Visible() bool             { return f.visible }
func (f *fakeObject) Z() int                    { return f.z }
func (f *fakeObject) ImageKey() string          { return f.imageKey }
func (f *fakeObject) Dst() (x0, y0, x1, y1 int) { return f.dst.c() }
func (f *fakeObject) Src() (x0, y0, x1, y1 int) { return f.src.c() }

func TestGC(t *testing.T) {
	objs := make(drawList, 8)
	for i := 0; i < 8; i++ {
		objs[i] = &fakeObject{}
	}
	for i := 0; i < (1 << 8); i++ {
		want := 0
		for j := uint32(0); j < 8; j++ {
			z := (i&(1<<j) != 0)
			objs[j].(*fakeObject).retire = z
			if !z {
				want++
			}
		}
		got := objs.gc()
		if len(got) != want {
			t.Errorf("Test %x: got length %d want %d", i, len(got), want)
		}
		for _, o := range got {
			if o.(*fakeObject).retire {
				t.Errorf("Test %x: object retired not removed", i)
			}
		}
	}
}
