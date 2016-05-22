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
	"sort"

	"github.com/DrJosh9000/vec"
	//"github.com/hajimehoshi/ebiten"
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
	return p.Pix, vec.I2(p.Rect.Max), nil
}

// loadTerrain loads from a paletted image file.
func loadTerrain(level *Level) (*Terrain, error) {
	terrain := &Terrain{
		Level:       level,
		blockSize:   vec.I2{level.TileSize, level.TileSize + level.BlockHeight},
		tilesetSize: sizes[level.TilesetKey].Div(level.TileSize),
	}
	/*
		// Prerender terrain layers to a texture.
		f, err := ebiten.NewImage(level.MapSize.X*level.TileSize, level.MapSize.Y*level.TileSize, ebiten.FilterNearest)
		if err != nil {
			return nil, fmt.Errorf("creating terrain texture: %v", err)
		}
		if err := Draw(f, level.TilesetKey, (*prerenderBaseLayer)(terrain)); err != nil {
			return nil, fmt.Errorf("drawing all tiles: %v", err)
		}
		terrain.baseLayer = f
		terrain.baseDrawOpts = &ebiten.DrawImageOptions{
			ImageParts: (*baseLayer)(terrain),
		}
	*/
	// Predraw all doodads, then do limited Z checking & redraw at draw time.
	sort.Sort(DoodadsByYPos(level.Doodads))
	/*
		for _, t := range level.Doodads {
			if err := (SpriteParts{t, false}.Draw(f)); err != nil {
				return nil, fmt.Errorf("drawing doodad: %v", err)
			}
		}*/

	return terrain, nil
}

func (t *Terrain) drawList() drawList {
	l := make(drawList, len(t.TileMap))
	for i, c := range t.TileMap {
		l[i] = &tileObject{
			Terrain: t,
			s:       vec.Div(int(c), t.tilesetSize.X).Mul(t.TileSize),
			d:       vec.Div(i, t.MapSize.X).Mul(t.TileSize),
		}
	}
	// TODO: add blocks here
	return l
}

type tileObject struct {
	*Terrain
	s, d vec.I2
}

func (t *tileObject) ImageKey() string { return t.TilesetKey }

func (t *tileObject) Dst() (x0, y0, x1, y1 int) {
	x0, y0 = t.d.C()
	x1, y1 = x0+t.TileSize, y0+t.TileSize
	return
}

func (t *tileObject) InWorld() bool { return true }

func (t *tileObject) Src() (x0, y0, x1, y1 int) {
	x0, y0 = t.s.C()
	x1, y1 = x0+t.TileSize, y0+t.TileSize
	return
}

func (t *tileObject) Visible() bool { return true }
func (t *tileObject) Z() int        { return -100 } // hax

/*
// prerenderBaseLayer implements ebiten.ImageParts for drawing the entire base layer.
type prerenderBaseLayer Terrain

// Len implements ebiten.ImageParts.
func (a *prerenderBaseLayer) Len() int {
	return a.MapSize.Area()
}

// Dst implements ebiten.ImageParts.
func (a *prerenderBaseLayer) Dst(i int) (x0, y0, x1, y1 int) {
	x0, y0 = vec.Div(i, a.MapSize.X).Mul(a.TileSize).C()
	x1, y1 = x0+a.TileSize, y0+a.TileSize
	return
}

// Src implements ebiten.ImageParts.
func (a *prerenderBaseLayer) Src(i int) (x0, y0, x1, y1 int) {
	x0, y0 = vec.Div(int(a.TileMap[i]), a.tilesetSize.X).Mul(a.TileSize).C()
	x1, y1 = x0+a.TileSize, y0+a.TileSize
	return
}

// baseLayer draws the base terrain layer to the screen.
type baseLayer Terrain

// Len implements ebiten.ImageParts.
func (a *baseLayer) Len() int {
	return 1
}

// Dst implements ebiten.ImageParts.
func (a *baseLayer) Dst(i int) (x0, y0, x1, y1 int) {
	x1, y1 = camSize.C()
	return
}

// Src implements ebiten.ImageParts.
func (a *baseLayer) Src(i int) (x0, y0, x1, y1 int) {
	x0, y0 = camPos.C()
	x1, y1 = camPos.Add(camSize).C()
	return
}

// blocksToY renders the block layer up to a Y limit.
type blocksToY struct {
	*Terrain
	y int
}

// Len implements ebiten.ImageParts.
func (a *blocksToY) Len() int {
	return 0
}

// Dst implements ebiten.ImageParts.
func (a *blocksToY) Dst(i int) (x0, y0, x1, y1 int) {
	return
}

// Src implements ebiten.ImageParts.
func (a *blocksToY) Src(i int) (x0, y0, x1, y1 int) {
	return
}

// blocksFromY renders the block layer after a Y limit.
type blocksFromY struct {
	*Terrain
	y int
}

// Len implements ebiten.ImageParts.
func (a *blocksFromY) Len() int {
	return 0
}

// Dst implements ebiten.ImageParts.
func (a *blocksFromY) Dst(i int) (x0, y0, x1, y1 int) {
	return
}

// Src implements ebiten.ImageParts.
func (a *blocksFromY) Src(i int) (x0, y0, x1, y1 int) {
	return
}
*/

// Terrain is the base layer of the game world.
type Terrain struct {
	*Level

	//baseLayer    *ebiten.Image // prerendered
	//baseDrawOpts *ebiten.DrawImageOptions
	blockSize vec.I2 // full size of each block (frame size for blockset)
	//blockSrc    *ebiten.Image
	tilesetSize vec.I2 // size of the tile map in tiles.
	//tileSrc     *ebiten.Image // tileset image
}

/*
func (t *Terrain) DrawBase(screen *ebiten.Image) error {
	return screen.DrawImage(t.baseLayer, t.baseDrawOpts)
}

func (t *Terrain) DrawBlocksToY(screen *ebiten.Image, y int) error {
	return screen.DrawImage(t.blockSrc, &ebiten.DrawImageOptions{
		ImageParts: &blocksToY{t, y},
	})
}

func (t *Terrain) DrawBlocksFromY(screen *ebiten.Image, y int) error {
	return screen.DrawImage(t.blockSrc, &ebiten.DrawImageOptions{
		ImageParts: &blocksFromY{t, y},
	})
}
*/

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

// ObstaclesAndPaths constructs two graphs, the first describing terrain
// obsctacles, the second describing a network of valid paths around
// the obstacles. Obstacles will be fattened according to the footprint
// fatUL, fatDR, and paths will be based on vertices at convex points of
// the obstacle graph plus 1 pixel in both dimensions outwards from the
// convex vertex.
func (t *Terrain) ObstaclesAndPaths(fatUL, fatDR vec.I2) (obstacles, paths *vec.Graph) {
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
			cup := t.Tile(i, j-1).Blocking || t.Block(i, j-1).Blocking
			cdown := t.Tile(i, j).Blocking || t.Block(i, j).Blocking
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
			cleft := t.Tile(i-1, j).Blocking || t.Block(i-1, j).Blocking
			cright := t.Tile(i, j).Blocking || t.Block(i, j).Blocking
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
		u := d.Pos().Sub(d.Anim().Offset)
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

	if Debug {
		log.Printf("generated %d obstacle edges", o.NumEdges())
	}

	// Precompute paths.
	p := vec.NewGraph()
	for u := range pVerts {
		for v := range pVerts {
			// Cull edges that are too tall/wide for the viewport.
			if vec.Abs(u.X-v.X) > camSize.X {
				continue
			}
			if vec.Abs(u.Y-v.Y) > camSize.Y {
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
	if Debug {
		log.Printf("generated %d paths edges", p.NumEdges())
	}
	return o, p
}
