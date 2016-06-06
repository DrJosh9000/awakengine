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

type Drawable interface {
	ImageKey() string
	Dst() (x0, y0, x1, y1 int) // relative to something else, or maybe nothing.
	Src() (x0, y0, x1, y1 int) // relative to the image referred to by ImageKey()
}

type Rect struct{ X0, Y0, X1, Y1 int }

func (r *Rect) C() (x0, y0, x1, y1 int) { return r.X0, r.Y0, r.X1, r.Y1 }

// StaticDrawable implements Drawable with struct fields.
type StaticDrawable struct {
	Key  string
	S, D Rect
}

func (s *StaticDrawable) ImageKey() string          { return s.Key }
func (s *StaticDrawable) Dst() (x0, y0, x1, y1 int) { return s.D.C() }
func (s *StaticDrawable) Src() (x0, y0, x1, y1 int) { return s.S.C() }

// Semiobject is some import stuff.
type Semiobject interface {
	Fixed() bool        // true if the object never moves - X, Y, or Z.
	Parent() Semiobject // for drawing and Z purposes.
	Retire() bool       // true if the object will never draw again.
	Visible() bool      // true if the object is visible if we are looking at it.
	Z() int
}

// StaticSemiobject implements Semiobjects with struct fields.
type StaticSemiobject struct {
	F, R, V bool
	P       Semiobject
	Zed     int
}

func (s *StaticSemiobject) Parent() Semiobject { return s.P }
func (s *StaticSemiobject) Fixed() bool        { return s.F }
func (s *StaticSemiobject) Retire() bool       { return s.R }
func (s *StaticSemiobject) Visible() bool      { return s.V }
func (s *StaticSemiobject) Z() int             { return s.Zed }

// Object is everything, everything is an object.
type Object interface {
	Drawable
	Semiobject
}

// ChildOf can be used to ensure an object is drawn relative to another.
type ChildOf struct{ Semiobject }

func (c ChildOf) Parent() Semiobject { return c.Semiobject }
func (c ChildOf) Z() int             { return c.Semiobject.Z() + 1 }
