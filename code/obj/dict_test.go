package obj

import (
	"fmt"
	"testing"
	"time"
)

func TestDict(t *testing.T) {
	const size = 1000000

	begin := time.Now()
	dict := NewDict()
	for i := 0; i < size; i++ {
		dict.Set([]byte(fmt.Sprintf("%d", i)), i)
	}
	fmt.Printf("dict cost=%s\n", time.Since(begin).String())

	begin = time.Now()
	mp := make(map[string]int)
	for i := 0; i < size; i++ {
		mp[fmt.Sprintf("%d", i)] = i
	}
	fmt.Printf("map cost=%s\n", time.Since(begin).String())
}
