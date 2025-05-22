package main

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/agnivade/levenshtein"
)

func splitSortConcat(a string) string {
	aa := strings.Split(a, " ")
	slices.Sort(aa)
	return strings.Join(aa, " ")
}

func main() {
	log.Printf("HEHEHEHE\n")

	a := "Road Nam Sing"
	b := "Old Airport Road Nam Sing"

	distance := levenshtein.ComputeDistance(splitSortConcat(a), splitSortConcat(b))
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}
	similarity := 1 - float64(distance)/float64(maxLen)

	fmt.Printf("Similarity: %.2f\n", similarity)

}
