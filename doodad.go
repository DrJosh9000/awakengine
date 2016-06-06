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

import "github.com/DrJosh9000/vec"

// BaseDoodad models a static object rendered at the player layer, but computed as
// obstacles (like terrain).
type BaseDoodad struct {
	*SheetFrame
	Offset vec.I2
	UL, DR vec.I2 // base obstacle box (frame relative)
}

// Doodad is an instance of a BaseDoodad in a specific location.
type Doodad struct {
	P vec.I2
	ChildOf
	*BaseDoodad
}

func (d *Doodad) Dst() (x0, y0, x1, y1 int) {
	x0, y0 = d.P.Sub(d.Offset).C()
	x1, y1 = x0+d.FrameSize.X, d.FrameSize.Y
	return
}

func (d *Doodad) Update(int)    {}
func (d *Doodad) Fixed() bool   { return true }
func (d *Doodad) Retire() bool  { return false }
func (d *Doodad) Visible() bool { return true }
func (d *Doodad) Z() int        { return d.P.Y }
