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

const hudZ = 90000

type HUDRegion struct {
	*Bubble
	Items []Drawable
	V, R  bool
}

func (h *HUDRegion) AddToScene(s *Scene) {
	h.Bubble.ChildOf = ChildOf{h}
	h.Bubble.AddToScene(s)
	for _, o := range h.Items {
		s.AddObject(&struct {
			Drawable
			ChildOf
		}{o, ChildOf{h.Bubble}})
	}
}

func (h *HUDRegion) Fixed() bool        { return true }
func (h *HUDRegion) Parent() Semiobject { return nil }
func (h *HUDRegion) Retire() bool       { return h.R }
func (h *HUDRegion) Visible() bool      { return h.V }
func (h *HUDRegion) Z() int             { return hudZ }
