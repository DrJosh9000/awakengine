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

type GridDelegate interface {
	Columns() int
	Item(i int, par *View)
	ItemSize() vec.I2
	ItemHandle(i int, e *Event) bool // item #i should handles event e
	NumItems() int
}

type Grid struct {
	*View
	GridDelegate

	items []*View
}

func (g *Grid) Handle(e *Event) bool {
	if !g.Bounds().Contains(e.ScreenPos) {
		return false
	}
	// Convert the coordinate back into an item index.
	t := e.ScreenPos.Sub(g.Position()).EDiv(g.GridDelegate.ItemSize())
	i := t.X + t.Y*g.GridDelegate.Columns()
	if i < 0 && i >= g.GridDelegate.NumItems() {
		return false
	}
	return g.GridDelegate.ItemHandle(i, e)
}

// Reload all the items.
func (g *Grid) Reload() {
	n := g.GridDelegate.NumItems()
	c := g.GridDelegate.Columns()
	sz := g.GridDelegate.ItemSize()
	gs := vec.I2{}
	for len(g.items) < n {
		g.items = append(g.items, &View{})
	}
	for i := 0; i < n; i++ {
		item := g.items[i]
		item.SetVisible(true)
		item.SetRetire(false)
		item.SetParent(g.View)
		p := vec.Div(i, c).EMul(sz)
		item.SetPositionAndSize(p, sz)
		item.SetZ(1)
		g.Item(i, item)

		// Ensure the grid itself is sized sufficiently. There's a mathsier way of
		// doing it but this is simple.
		if w := p.X + sz.X; w > gs.X {
			gs.X = w
		}
		if h := p.Y + sz.Y; h > gs.Y {
			gs.Y = h
		}
	}
	g.View.SetSize(gs)
	for i := n; i < len(g.items); i++ {
		g.items[i].SetRetire(true)
	}
}
