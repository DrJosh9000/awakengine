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

import (
	"image/color"

	"github.com/DrJosh9000/vec"
)

// GraphView visualises a graph.
type GraphView struct {
	edges        []vec.Edge
	edgeColour   color.Color
	normalColour color.Color
}

// Len implements ebiten.Lines.
func (v GraphView) Len() int {
	return len(v.edges) * 2 // edges and normals.
}

// Points implements ebiten.Lines.
func (v GraphView) Points(i int) (x0, y0, x1, y1 int) {
	l := len(v.edges)
	if i < l {
		e := v.edges[i]
		a, b := e.U.Sub(camPos), e.V.Sub(camPos)
		return a.X, a.Y, b.X, b.Y
	}
	e := v.edges[i-l]
	a := e.U.Div(2).Add(e.V.Div(2)).Sub(camPos)
	b := e.V.Sub(e.U).Normal().Div(4).Add(a)
	return a.X, a.Y, b.X, b.Y
}

// Color implements ebiten.Lines.
func (v GraphView) Color(i int) color.Color {
	if i < len(v.edges) {
		return v.edgeColour
	}
	return v.normalColour
}
