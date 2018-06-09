package main

import (
	"fmt"
)

type A struct {
	a []int
}

type B struct {

}

func (a *A) ingre() {
	a.a = []int{1,2}
}

func main() {
	a := A{}
	a.ingre()
	for k, v := range a.a {
		if v == 2{
			a.a = append(a.a[:k], a.a[k+1:]...)
		}
	}
	fmt.Println(a.a)
}


