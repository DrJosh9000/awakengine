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
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"

	"github.com/DrJosh9000/vec"
)

// TileInfo describes the properties of a tile.
type TileInfo struct {
	Name     string
	Blocking bool // Player is unable to walk through?
}

// ImageAsMap returns the contents and size of a paletted PNG file.
func ImageAsMap(imgkey string) ([]uint8, vec.I2, error) {
	pngData, ok := allData[imgkey]
	if !ok {
		return nil, vec.I2{}, fmt.Errorf("source %q not a registered image", imgkey)
	}
	i, err := png.Decode(bytes.NewReader(pngData))
	if err != nil {
		return nil, vec.I2{}, fmt.Errorf("loading source png: %v", err)
	}
	p, ok := i.(*image.Paletted)
	if !ok {
		return nil, vec.I2{}, fmt.Errorf("source png is not paletted [%T != *image.Paletted]", i)
	}
	//log.Printf("%s: loaded map %v", imgkey, p.Pix)
	return p.Pix, vec.I2(p.Rect.Max), nil
}

type tilePart struct {
	*Terrain
	i   int // Keep an index in case the map updates dynamically!
	d   vec.I2
	vis bool
}

func (t *tilePart) ImageKey() string { return t.TilesetKey }

func (t *tilePart) Dst() (x0, y0, x1, y1 int) {
	x0, y0 = t.d.C()
	x1, y1 = x0+t.TileSize, y0+t.TileSize
	return
}

func (t *tilePart) Src() (x0, y0, x1, y1 int) {
	x0, y0 = vec.Div(int(t.TileMap[t.i]), t.tilesetSize.X).Mul(t.TileSize).C()
	x1, y1 = x0+t.TileSize, y0+t.TileSize
	return
}

func (t *tilePart) Visible() bool { return t.vis && t.Terrain.Visible() }
func (t *tilePart) Z() int        { return -100 } // hax

type blockPart struct {
	*Terrain
	d    vec.I2
	i, z int
	vis  bool
}

func (b *blockPart) ImageKey() string { return b.BlocksetKey }

func (b *blockPart) Dst() (x0, y0, x1, y1 int) {
	x0, y0 = b.d.C()
	x1, y1 = b.d.Add(b.blockSize).C()
	return
}

func (b *blockPart) Src() (x0, y0, x1, y1 int) {
	x0, y0 = vec.Div(int(b.BlockMap[b.i]), b.blocksetSize.X).EMul(b.blockSize).C()
	x1, y1 = x0+b.blockSize.X, y0+b.blockSize.Y
	return
}

func (b *blockPart) Visible() bool { return b.vis && b.Terrain.Visible() }
func (b *blockPart) Z() int        { return b.z }

// Terrain is the base layer of the game world.
type Terrain struct {
	*View
	*Level

	blockSize    vec.I2 // full size of each block (frame size for blockset)
	blocksetSize vec.I2 // size of the block map in blocks.
	tilesetSize  vec.I2 // size of the tile map in tiles.

	tileParts  map[int]*tilePart
	blockParts map[int]*blockPart
}

// loadTerrain loads from a paletted image file.
func loadTerrain(level *Level, parent *View) (*Terrain, error) {
	bs := vec.I2{level.TileSize, level.TileSize + level.BlockHeight}
	t := &Terrain{
		View:       &View{},
		Level:      level,
		blockSize:  bs,
		tileParts:  make(map[int]*tilePart),
		blockParts: make(map[int]*blockPart),
	}
	t.View.SetParent(parent)
	if level.BlocksetKey != "" {
		t.blocksetSize = sizes[level.BlocksetKey].EDiv(bs)
	}
	if level.TilesetKey != "" {
		t.tilesetSize = sizes[level.TilesetKey].Div(level.TileSize)
	}

	for i := range t.TileMap {
		if t.TileMap[i] == 0 {
			continue
		}
		t.tileParts[i] = &tilePart{
			Terrain: t,
			i:       i,
			d:       vec.Div(i, t.MapSize.X).Mul(t.TileSize),
		}
	}
	for i := range t.BlockMap {
		if t.BlockMap[i] == 0 {
			continue
		}
		d := vec.Div(i, t.MapSize.X).Mul(t.TileSize)
		t.blockParts[i] = &blockPart{
			Terrain: t,
			i:       i,
			d:       d.Sub(vec.I2{0, t.BlockHeight}),
			z:       d.Y,
		}
	}
	return t, nil
}

// AddToScene adds terrain objects to the scene.
func (t *Terrain) AddToScene(s *Scene) {
	for _, p := range t.tileParts {
		s.AddPart(p)
	}
	for _, p := range t.blockParts {
		s.AddPart(p)
	}
}

// MakeAllVisible enables visiblity for all parts of the terrain.
func (t *Terrain) MakeAllVisible() {
	for _, p := range t.tileParts {
		p.vis = true
	}
	for _, p := range t.blockParts {
		p.vis = true
	}
}

func (t *Terrain) Fixed() bool  { return true }
func (t *Terrain) Retire() bool { return false }

type ray struct {
	*Terrain
	vis     bool
	n, dist int
}

func (r *ray) touch(p vec.I2) bool {
	//log.Printf("touching %v", p)
	if p.X < 0 || p.X >= r.MapSize.X || p.Y < 0 || p.Y >= r.MapSize.Y {
		return false
	}
	i := p.X + r.MapSize.X*p.Y
	if tp, ok := r.tileParts[i]; ok {
		tp.vis = tp.vis || r.vis
	}
	if bp, ok := r.blockParts[i]; ok {
		bp.vis = bp.vis || r.vis
	}
	if r.BlockInfos[r.BlockMap[i]].Blocking {
		r.vis = false
	}
	r.n++
	return r.n <= r.dist
}

// UpdatePartVisibility makes tile & block parts visible if they can be seen from origin.
func (t *Terrain) UpdatePartVisibility(origin vec.I2, dist int) {
	/*for _, p := range t.tileParts {
		p.vis = false
	}
	for _, p := range t.blockParts {
		p.vis = false
	}*/
	var r *ray
	originCell := t.TileCoord(origin)
	cellSize := vec.I2{t.TileSize, t.TileSize}
	offset := cellSize.Div(2)
	for x := originCell.X - dist; x <= originCell.X+dist; x++ {
		end := vec.I2{x, originCell.Y - dist}.Mul(t.TileSize).Add(vec.I2{offset.X, cellSize.Y - 1})
		r = &ray{t, true, 0, dist}
		vec.CellsTouchingSegment(cellSize, origin, end, r.touch)
		end = vec.I2{x, originCell.Y + dist}.Mul(t.TileSize).Add(vec.I2{offset.X, 0})
		r = &ray{t, true, 0, dist}
		vec.CellsTouchingSegment(cellSize, origin, end, r.touch)
	}
	for y := originCell.Y - dist; y <= originCell.Y+dist; y++ {
		end := vec.I2{originCell.X - dist, y}.Mul(t.TileSize).Add(vec.I2{cellSize.X - 1, offset.Y})
		r = &ray{t, true, 0, dist}
		vec.CellsTouchingSegment(cellSize, origin, end, r.touch)
		end = vec.I2{originCell.X + dist, y}.Mul(t.TileSize).Add(vec.I2{cellSize.X - 1, 0})
		r = &ray{t, true, 0, dist}
		vec.CellsTouchingSegment(cellSize, origin, end, r.touch)
	}
}

// TileCoord returns information about the tile at a world coordinate.
func (t *Terrain) TileCoord(wc vec.I2) vec.I2 { return wc.Div(t.TileSize) }

// Size returns the world size in pixels.
func (t *Terrain) Size() vec.I2 { return t.MapSize.Mul(t.TileSize) }

// Tile gets the information about the tile at a tile coordinate.
func (t *Terrain) Tile(x, y int) TileInfo {
	if x < 0 || x >= t.MapSize.X || y < 0 || y >= t.MapSize.Y {
		return TileInfo{Name: "out-of-bounds", Blocking: true}
	}
	i := x + t.MapSize.X*y
	n := t.TileMap[i]
	return t.TileInfos[n]
}

// Block gets the information about the block at a tile coordinate.
func (t *Terrain) Block(x, y int) TileInfo {
	if x < 0 || x >= t.MapSize.X || y < 0 || y >= t.MapSize.Y {
		return TileInfo{Name: "out-of-bounds", Blocking: true}
	}
	i := x + t.MapSize.X*y
	n := t.BlockMap[i]
	return t.BlockInfos[n]
}

func (t *Terrain) Blocking(i, j int) bool {
	if t.TileMap != nil && t.Tile(i, j).Blocking {
		return true
	}
	if t.BlockMap != nil && t.Block(i, j).Blocking {
		return true
	}
	return false
}

// ObstaclesAndPaths constructs two graphs, the first describing terrain
// obsctacles, the second describing a network of valid paths around
// the obstacles. Obstacles will be fattened according to the footprint
// fatUL, fatDR, and paths will be based on vertices at convex points of
// the obstacle graph plus 1 pixel in both dimensions outwards from the
// convex vertex.
func (t *Terrain) ObstaclesAndPaths(fatUL, fatDR, limit vec.I2) (obstacles, paths *vec.Graph) {
	o := vec.NewGraph()
	// Store a separate vertex set for path generation, because we only care
	// about convex corners.
	pVerts := make(vec.VertexSet)
	fatUR, fatDL := vec.I2{fatDR.X, fatUL.Y}, vec.I2{fatUL.X, fatDR.Y}
	ul, ur, dl, dr := vec.I2{-1, -1}, vec.I2{1, -1}, vec.I2{-1, 1}, vec.I2{1, 1}

	// Generate edges along rows.
	for j := 0; j <= t.MapSize.Y; j++ {
		up, down := true, true
		u := vec.I2{}
		for i := 0; i < t.MapSize.X; i++ {
			ut := vec.I2{i, j}.Mul(t.TileSize)
			cup := t.Blocking(i, j-1)
			cdown := t.Blocking(i, j)
			if up != cup || down != cdown {
				if up && !down {
					if cdown {
						// concave
						v := ut.Add(fatDL)
						o.AddEdge(u, v)
					} else {
						// convex
						v := ut.Add(fatDR)
						o.AddEdge(u, v)
						pVerts[v.Add(dr)] = true
					}
				}
				if !up && down {
					if cup {
						// concave
						v := ut.Add(fatUL)
						o.AddEdge(v, u)
					} else {
						v := ut.Add(fatUR)
						o.AddEdge(v, u)
						pVerts[v.Add(ur)] = true
					}
				}
				if cup && !cdown {
					if down {
						// concave
						u = ut.Add(fatDR)
					} else {
						u = ut.Add(fatDL)
						pVerts[u.Add(dl)] = true
					}
				}
				if !cup && cdown {
					if up {
						// concave
						u = ut.Add(fatUR)
					} else {
						u = ut.Add(fatUL)
						pVerts[u.Add(ul)] = true
					}
				}
			}
			up, down = cup, cdown
		}
	}

	// Generate edges along columns.
	for i := 0; i <= t.MapSize.X; i++ {
		left, right := true, true
		u := vec.I2{}
		for j := 0; j < t.MapSize.Y; j++ {
			ut := vec.I2{i, j}.Mul(t.TileSize)
			cleft := t.Blocking(i-1, j)
			cright := t.Blocking(i, j)
			if left != cleft || right != cright {
				if left && !right {
					if cright {
						// concave
						v := ut.Add(fatUR)
						o.AddEdge(v, u)
					} else {
						v := ut.Add(fatDR)
						o.AddEdge(v, u)
						pVerts[v.Add(dr)] = true
					}
				}
				if !left && right {
					if cleft {
						// concave
						v := ut.Add(fatUL)
						o.AddEdge(u, v)
					} else {
						v := ut.Add(fatDL)
						o.AddEdge(u, v)
						pVerts[v.Add(dl)] = true
					}
				}
				if cleft && !cright {
					if right {
						// concave
						u = ut.Add(fatDR)
					} else {
						u = ut.Add(fatUR)
						pVerts[u.Add(ur)] = true
					}
				}
				if !cleft && cright {
					if left {
						// concave
						u = ut.Add(fatDL)
					} else {
						u = ut.Add(fatUL)
						pVerts[u.Add(ul)] = true
					}
				}
			}
			left, right = cleft, cright
		}
	}

	// Generate doodad edges
	for _, d := range t.Doodads {
		u := d.P.Sub(d.Offset)
		u, v := u.Add(d.UL).Add(fatUL), u.Add(d.DR).Add(fatDR)
		uv, vu := vec.I2{u.X, v.Y}, vec.I2{v.X, u.Y}
		o.AddEdge(u, uv)
		o.AddEdge(uv, v)
		o.AddEdge(v, vu)
		o.AddEdge(vu, u)
		pVerts[u.Add(ul)] = true
		pVerts[uv.Add(dl)] = true
		pVerts[v.Add(dr)] = true
		pVerts[vu.Add(ur)] = true
	}

	if config.Debug {
		log.Printf("generated %d vertices", len(pVerts))
		log.Printf("generated %d obstacle edges", o.NumEdges())
	}

	// Precompute paths.
	p := vec.NewGraph()
	for u := range pVerts {
		for v := range pVerts {
			// Cull edges that are too tall/wide for the viewport.
			if vec.Abs(u.X-v.X) > limit.X {
				continue
			}
			if vec.Abs(u.Y-v.Y) > limit.Y {
				continue
			}
			// Cull edges that intersect an obstacle. Do this for backfacing obstacle edges,
			// because u might be contained in another obstacle.
			if o.FullyBlocks(u, v) {
				continue
			}
			p.AddEdge(u, v)
		}
	}
	if config.Debug {
		log.Printf("generated %d paths edges", p.NumEdges())
	}
	return o, p
}
