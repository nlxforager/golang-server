package concurrency

import (
	"log"
	"testing"
)

func BenchmarkFire(b *testing.B) {
	log.Printf("b.N=%d\n", b.N)
	for i := 0; i < b.N; i++ {
		i := i
		log.Printf("hello world %d\n", i)
		//b.Run("__"+strconv.Itoa(i), func(b *testing.B) {

		Fire(i + 1)
		//})
	}
}
