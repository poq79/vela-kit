package runtime

import (
	"container/ring"
	"fmt"
	"testing"
)

func TestRing(t *testing.T) {
	r := ring.New(5)

	// Get the length of the ring
	n := r.Len()

	// Initialize the ring with some integer values
	for i := 0; i < n+3; i++ {
		r.Value = i
		r = r.Next()
	}

	r.Do(func(p any) {
		fmt.Println(p.(int))
	})

	fmt.Println("-----")

	for i := 0; i < n-2; i++ {
		r.Value = i
		r = r.Next()
	}

	// Iterate through the ring and print its contents
	r.Do(func(p any) {
		fmt.Println(p.(int))
	})
}
