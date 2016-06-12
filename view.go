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

// View represents a rectangular region in a view hierarchy. It caches its
// real position information because why not.
type View struct {
	bounds    vec.Rect // Relative to parent hierarchy
	z         int      // Relative to parent
	invisible bool     // So we can avoid explicitly setting to visible.
	retire    bool     // So we can avoid explicitly setting to not-retired.

	parent     *View
	children   []*View
	childIndex int // My index within my parent's children slice.

	valid           bool
	cachedBounds    vec.Rect
	cachedRetire    bool
	cachedInvisible bool
	cachedZ         int
}

// Dispose retires this view and all subviews, and disconnects
// everything.
func (v *View) Dispose() {
	for len(v.children) > 0 {
		v.children[0].Dispose()
	}
	v.retire = true
	v.cachedRetire = true
	v.SetParent(nil)
}

func (v *View) invalidate() {
	if !v.valid {
		// Children should be invalid. If a child view is valid,
		// they will have forced us to become valid in doing so.
		return
	}
	v.valid = false
	for _, c := range v.children {
		c.invalidate()
	}
}

func (v *View) compute() {
	if v.valid {
		return
	}
	v.valid = true
	if v.parent == nil {
		v.cachedBounds = v.bounds
		v.cachedRetire = v.retire
		v.cachedInvisible = v.invisible
		v.cachedZ = v.z
		return
	}
	v.cachedBounds = v.bounds.Translate(v.parent.Position())
	v.cachedRetire = v.retire || v.parent.Retire()
	v.cachedInvisible = v.invisible || !v.parent.Visible()
	v.cachedZ = v.z + v.parent.Z()
}

func (v *View) LogicalBounds() vec.Rect { return v.bounds }
func (v *View) LogicalZ() int           { return v.z }
func (v *View) Size() vec.I2            { return v.bounds.Size() }

// Convenience methods for embedding in Parts.
func (v *View) Container() *View { return v }
func (v *View) Fixed() bool      { return true }

func (v *View) Retire() bool {
	v.compute()
	return v.cachedRetire
}

func (v *View) Visible() bool {
	v.compute()
	return !v.cachedInvisible
}

func (v *View) Z() int {
	v.compute()
	return v.cachedZ
}

func (v *View) Bounds() vec.Rect {
	v.compute()
	return v.cachedBounds
}

// Position is always the view's top-left corner.
func (v *View) Position() vec.I2 {
	v.compute()
	return v.cachedBounds.UL
}

func (v *View) SetBounds(bounds vec.Rect) {
	v.bounds = bounds
	v.invalidate()
}

// SetPosition moves the view to a new position without altering the size.
// It invalidates.
func (v *View) SetPosition(ul vec.I2) {
	v.bounds = v.bounds.Reposition(ul)
	v.invalidate()
}

// SetSize alters the size of the view without altering the top-left corner of the bounds.
// It invalidates nothing.
func (v *View) SetSize(sz vec.I2) {
	v.bounds = v.bounds.Resize(sz)
}

// SetPosition moves the view to a new position and changes its size.
// It invalidates.
func (v *View) SetPositionAndSize(pos, size vec.I2) {
	v.bounds = vec.Rect{UL: pos, DR: pos.Add(size)}
	v.invalidate()
}

func (v *View) SetRetire(ret bool) {
	v.retire = ret
	v.invalidate()
}

func (v *View) SetVisible(vis bool) {
	v.invisible = !vis
	v.invalidate()
}

func (v *View) SetZ(z int) {
	v.z = z
	v.invalidate()
}

func (v *View) removeFromParent() {
	if v.parent == nil {
		return
	}
	c := v.parent.children
	switch len(c) {
	case 0:
		panic("parent should have children but doesn't")
	case 1:
		if v.childIndex != 0 || c[0] != v {
			panic("my idea of which child I am is wrong")
		}
		v.parent.children = nil
	default:
		if c[v.childIndex] != v {
			panic("my idea of which child I am is wrong")
		}
		last := len(c) - 1
		c[last].childIndex = v.childIndex
		c[v.childIndex] = c[last]
		v.parent.children = c[:last]
	}
}

func (v *View) SetParent(parent *View) {
	v.invalidate()
	v.removeFromParent()
	v.parent = parent
	if parent == nil {
		return
	}
	v.childIndex = len(v.parent.children)
	v.parent.children = append(v.parent.children, v)
}
