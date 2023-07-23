package main

import (
	"fmt"
	"unsafe"
)

type person struct {
	age    int
	weight int
	height int
}

func main() {



	pos := unsafe.Offsetof(person{}.height)
	fmt.Println(pos)
	p := &person{
		age:    1,
		weight: 12,
		height: 14,
	}
	ppos := unsafe.Offsetof(p.height)

	fmt.Println(ppos)
	fmt.Println("-------------------------------")
	a := "aa"
	dst := (*[10]byte)(unsafe.Pointer(&a))[:5]
	fmt.Println(dst)
	fmt.Println("we have to save our life")
}
