// Copyright 2016 Josh Deprez
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package awakengine

import "testing"

type fakeObject struct {
	inWorld, retire, visible bool
	z                        int
	imageKey                 string
	dst, src                 Rect
}

func (f *fakeObject) InWorld() bool             { return f.inWorld }
func (f *fakeObject) Retire() bool              { return f.retire }
func (f *fakeObject) Visible() bool             { return f.visible }
func (f *fakeObject) Z() int                    { return f.z }
func (f *fakeObject) ImageKey() string          { return f.imageKey }
func (f *fakeObject) Dst() (x0, y0, x1, y1 int) { return f.dst.C() }
func (f *fakeObject) Src() (x0, y0, x1, y1 int) { return f.src.C() }

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
