package partition

import (
	"github.com/adm87/finch-core/geom"
	"github.com/adm87/finch-core/hashset"
)

// Partitioning defines the interface for spatial partitioning structures.
type Partitioning[T geom.Bounded] interface {
	Insert(item T) bool                                   // Add an item to the partitioning structure
	Remove(item T) bool                                   // Remove an item from the partitioning structure
	Query(area geom.Rect64) hashset.Set[T]                // Query items within a specified area
	Update(item T)                                        // Update an item's position in the partitioning structure
	Clear()                                               // Clear all items from the partitioning structure
	Count() int                                           // Get the number of items in the partitioning structure
	Partitions(area geom.Rect64) hashset.Set[geom.Rect64] // Get the partitions used in the structure
}
