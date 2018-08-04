package main

import "github.com/jinzhu/copier"

type Nested struct {
	value *int
}

type Wrapper struct {
	nested *Nested
}

func main()  {
	var a = 1
	n := Nested{&a}
	w := Wrapper{&n}
	var newIns = Wrapper{}
	copier.Copy(&newIns, &w)
	a = 3
	println(&w, w.nested)
	println(&newIns, newIns.nested)
}
