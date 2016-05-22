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

import "github.com/hajimehoshi/ebiten"

// PartWindow windows an ImageParts to the index range [Start,Start+N).
type PartWindow struct {
	ebiten.ImageParts
	Start, N int
}

// Len implements ebiten.ImageParts.
func (w *PartWindow) Len() int { return w.N }

// Dst implements ebiten.ImageParts.
func (w *PartWindow) Dst(i int) (x0, y0, x1, y1 int) {
	return w.ImageParts.Dst(i - w.Start)
}

// Src implements ebiten.ImageParts.
func (w *PartWindow) Src(i int) (x0, y0, x1, y1 int) {
	return w.ImageParts.Src(i - w.Start)
}
