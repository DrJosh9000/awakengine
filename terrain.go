package awakengine

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"sort"

	"github.com/DrJosh9000/vec"
	"github.com/hajimehoshi/ebiten"
)

// loadTerrain loads from a paletted image file.
func loadTerrain(l Level) error {
	pngData, ok := allData[l.Source()]
	if !ok {
		return fmt.Errorf("level source %q not a registered image", l.Source())
	}

	i, err := png.Decode(bytes.NewReader(pngData))
	if err != nil {
		return err
	}
	p, ok := i.(*image.Paletted)
	if !ok {
		return fmt.Errorf("source png is not a paletted png [%T != *image.Paletted]", i)
	}
	terrain.size = vec.I2{p.Rect.Max.X, p.Rect.Max.Y}
	terrain.tiles = p.Pix

	tilesKey, tileSize := l.Tiles()
	tilesImg, ok := allImages[tilesKey]
	if !ok {
		return fmt.Errorf("level tiles %q not a registered image", tilesKey)
	}

	// Prerender terrain to a single texture.
	f, err := ebiten.NewImage(p.Rect.Max.X*tileSize, p.Rect.Max.Y*tileSize, ebiten.FilterNearest)
	if err != nil {
		return err
	}
	if err := f.DrawImage(Image(tileKey), &ebiten.DrawImageOptions{ImageParts: (*AllTerrainParts)(terrain)}); err != nil {
		return err
	}

	// Predraw all doodads, then do limited Z checking & redraw at draw time.
	sort.Sort(DoodadsByYPos(terrain.doodads))
	for _, d := range terrain.doodads {
		if err := (SpriteParts{d, false}.Draw(f)); err != nil {
			return err
		}
	}
	terrain.flat = f
	return nil
}

// TileInfo describes the properties of a tile.
type TileInfo struct {
	name  string
	block bool // Player unable to walk through
}

// AllTerrainParts implements ebiten.ImageParts for drawing the entire terrain.
type AllTerrainParts Terrain

// Len implements ebiten.ImageParts.
func (a *AllTerrainParts) Len() int {
	return a.size.Area()
}

// Dst implements ebiten.ImageParts.
func (a *AllTerrainParts) Dst(i int) (x0, y0, x1, y1 int) {
	x0, y0 = vec.Div(i, a.size.X).Mul(tileSize).C()
	x1, y1 = x0+tileSize, y0+tileSize
	return
}

// Src implements ebiten.ImageParts.
func (a *AllTerrainParts) Src(i int) (x0, y0, x1, y1 int) {
	x0, y0 = vec.Div(int(a.tiles[i]), tileMapWidth).Mul(tileSize).C()
	x1, y1 = x0+tileSize, y0+tileSize
	return
}

// Terrain is the base layer of the game world.
type Terrain struct {
	tiles    []uint8       // tilemap index at terrain position (x, y) is tiles[x+y*size.X].
	tileInfo []TileInfo    // info for tile index
	size     vec.I2        // in tiles, not pixels.
	doodads  []*Doodad     // for the doodad throne
	flat     *ebiten.Image // prerendered
}

// Draw draws the terrain to the screen.
func (t *Terrain) Draw(screen *ebiten.Image) error {
	return screen.DrawImage(t.flat, &ebiten.DrawImageOptions{ImageParts: t})
}

// Len implements ebiten.ImageParts.
func (t *Terrain) Len() int { return 1 }

// Dst implements ebiten.ImageParts.
func (t *Terrain) Dst(int) (x0, y0, x1, y1 int) {
	x1, y1 = camSize.C()
	return
}

// Src implements ebiten.ImageParts.
func (t *Terrain) Src(int) (x0, y0, x1, y1 int) {
	x0, y0 = camPos.C()
	x1, y1 = camPos.Add(camSize).C()
	return
}

// Query returns information about the tile at a world coordinate.
func (t *Terrain) Query(wc vec.I2) TileInfo {
	tx, ty := wc.Div(tileSize).C()
	if tx >= 0 && tx < t.size.X && ty >= 0 && ty < t.size.Y {
		return t.tileInfo[t.tiles[tx+ty*t.size.X]]
	}
	return t.tileInfo[blackTile]
}

// Tile is the
func (t *Terrain) Tile(x, y int) TileInfo {
	if x < 0 || x >= t.size.X || y < 0 || y >= t.size.Y {
		return t.TileInfo{block: true}
	}
	return t.tileInfo[t.tiles[x+t.size.X*y]]
}

// ObstaclesAndPaths constructs two graphs, the first describing terrain
// obsctacles, the second describing a network of valid paths around
// the obstacles.
func (t *Terrain) ObstaclesAndPaths(fatUL, fatDR vec.I2) (obstacles, paths *vec.Graph) {
	o := vec.NewGraph()
	// Store a separate vertex set for path generation, because we only care
	// about convex corners.
	pVerts := make(vec.VertexSet)
	fatUR := vec.I2{fatDR.X, fatUL.Y}
	fatDL := vec.I2{fatUL.X, fatDR.Y}

	ul, ur, dl, dr := vec.I2{-1, -1}, vec.I2{1, -1}, vec.I2{-1, 1}, vec.I2{1, 1}

	// Generate edges along rows.
	for j := 0; j <= t.size.Y; j++ {
		up, down := true, true
		u := vec.I2{}
		for i := 0; i < t.size.X; i++ {
			ut := vec.I2{i, j}.Mul(tileSize)
			cup := t.Tile(i, j-1).block
			cdown := t.Tile(i, j).block
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
	for i := 0; i <= t.size.X; i++ {
		left, right := true, true
		u := vec.I2{}
		for j := 0; j < t.size.Y; j++ {
			ut := vec.I2{i, j}.Mul(tileSize)
			cleft := t.Tile(i-1, j).block
			cright := t.Tile(i, j).block
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
	for _, d := range t.doodads {
		u := d.pos.Sub(d.anim.offset)
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
