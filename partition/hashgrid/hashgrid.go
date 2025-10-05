package hashgrid

import (
	"math"

	"github.com/adm87/finch-core/geom"
	"github.com/adm87/finch-core/hashset"
)

// ======================================================
// Cell Key
// ======================================================

// CellKey represents the coordinates of a cell in the hash grid.
type CellKey struct {
	X, Y int
}

// ======================================================
// Hash Grid
// ======================================================

// HashGrid is a spatial partitioning structure that divides space into a uniform grid of cells.
// Each cell can contain multiple items of type T, which must implement the geom.Bounded interface.
type HashGrid[T geom.Bounded] struct {
	cellSize float64
	cellKeys map[T]hashset.Set[CellKey]
	cells    map[CellKey]hashset.Set[T]
}

// New creates a new HashGrid with the specified cell size.
func New[T geom.Bounded](cellSize float64) *HashGrid[T] {
	return &HashGrid[T]{
		cellSize: cellSize,
		cellKeys: make(map[T]hashset.Set[CellKey]),
		cells:    make(map[CellKey]hashset.Set[T]),
	}
}

func (hg *HashGrid[T]) Insert(item T) bool {
	if hg.cellKeys[item] != nil {
		hg.Update(item)
		return false
	}

	bounds := item.Bounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		panic("cannot insert item with non-positive width or height bounding box")
	}

	keys := hg.getCellKeys(bounds)
	hg.cellKeys[item] = hashset.From(keys...)

	for _, key := range keys {
		if hg.cells[key] == nil {
			hg.cells[key] = hashset.New[T]()
		}
		hg.cells[key].AddDistinct(item)
	}

	return true
}

func (hg *HashGrid[T]) Remove(item T) bool {
	if hg.cellKeys[item] == nil {
		return false
	}

	keys := hg.cellKeys[item]
	for key := range keys {
		if cell, ok := hg.cells[key]; ok {
			cell.Remove(item)
			if len(cell) == 0 {
				delete(hg.cells, key)
			}
		}
	}

	delete(hg.cellKeys, item)
	return true
}

func (hg *HashGrid[T]) Query(area geom.Rect64) hashset.Set[T] {
	result := hashset.New[T]()
	keys := hg.getCellKeys(area)

	for _, key := range keys {
		if cell, exists := hg.cells[key]; exists {
			for item := range cell {
				result.AddDistinct(item)
			}
		}
	}

	return result
}

func (hg *HashGrid[T]) Update(item T) {
	hg.Remove(item)
	hg.Insert(item)
}

func (hg *HashGrid[T]) Clear() {
	hg.cellKeys = make(map[T]hashset.Set[CellKey])
	hg.cells = make(map[CellKey]hashset.Set[T])
}

func (hg *HashGrid[T]) Count() int {
	return len(hg.cellKeys)
}

func (hg *HashGrid[T]) Partitions(area geom.Rect64) hashset.Set[geom.Rect64] {
	result := hashset.New[geom.Rect64]()
	keys := hg.getCellKeys(area) // Only get relevant cells

	for _, key := range keys {
		if _, exists := hg.cells[key]; exists {
			rect := geom.Rect64{
				X:      float64(key.X) * hg.cellSize,
				Y:      float64(key.Y) * hg.cellSize,
				Width:  hg.cellSize,
				Height: hg.cellSize,
			}
			result.AddDistinct(rect)
		}
	}

	return result
}

func (hg *HashGrid[T]) getCellKeys(area geom.Rect64) []CellKey {
	minX := int(math.Floor((area.X - hg.cellSize*0.5) / hg.cellSize))
	minY := int(math.Floor((area.Y - hg.cellSize*0.5) / hg.cellSize))
	maxX := int(math.Floor((area.X+area.Width-hg.cellSize*0.5)/hg.cellSize)) + 1
	maxY := int(math.Floor((area.Y+area.Height-hg.cellSize*0.5)/hg.cellSize)) + 1

	keys := make([]CellKey, 0, (maxX-minX+1)*(maxY-minY+1))
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			keys = append(keys, CellKey{X: x, Y: y})
		}
	}

	return keys
}
