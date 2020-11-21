package diff

import (
	"fmt"
	"testing"
)

func TestIntSlice(t *testing.T) {
	src := []int64{
		1,
		2,
		3,
	}
	in := []int64{
		1,
		2,
		3,
	}
	add, del := IntSlice(src, in)
	fmt.Println("add:", add)
	fmt.Println("del:", del)
	in = append(in, 4, 5, 6)
	add, del = IntSlice(src, in)
	fmt.Println("add:", add)
	fmt.Println("del:", del)
	in = []int64{
		7,
		8,
		9,
	}
	add, del = IntSlice(src, in)
	fmt.Println("add:", add)
	fmt.Println("del:", del)
}
