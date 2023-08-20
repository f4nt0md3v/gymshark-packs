package services

import (
	"sort"

	"github.com/f4nt0md3v/gymshark-packs/pkg/utilx"
)

type PackService struct{}

func NewPackService() *PackService {
	return &PackService{}
}

type Packer interface {
	CalculateNumberOfPacks(items int, packSizes []int) map[int]int
}

// CalculateNumberOfPacks calculates the number of packs needed to ship the given number of items.
// This function calculates the number of packs needed to ship the given number of items.
func (s *PackService) CalculateNumberOfPacks(items int, packSizes []int) map[int]int {
	var (
		// The minimum number of packs is 1.
		minPacks = 1
		packMap  = make(map[int]int)
	)

	// Sort pack sizes first
	sort.Ints(packSizes)

	// Loop through the pack sizes.
	for i, packSize := range packSizes {
		if items == 0 {
			break
		}

		// Calculate the number of packs that can be created from the current pack size
		// and the remainder of items that cannot fit with the current pack size.
		packsFromCurrentSize := items / packSize
		remainderWithCurrentSize := items % packSize
		if packsFromCurrentSize > 1 && remainderWithCurrentSize >= 0 && i != len(packSizes)-1 {
			continue
		}

		// If the pack with current size can fit the remaining items, then that is the minimum number of packs needed.
		if packsFromCurrentSize == 0 {
			packMap[packSize] = minPacks
			break
		}

		// We need to ship at least that many packs.
		items -= packsFromCurrentSize * packSize
		packMap[packSize] = packsFromCurrentSize
		packMap = utilx.MergeMaps(packMap, s.CalculateNumberOfPacks(items, packSizes))

		break
	}

	return packMap
}
