package math

var X string = "hello X"
var x string = "hello x"

type Math struct {
	A int //se colocar minúsculo, não será exportado
	B int
}

// método Add() vinculado a struct Math
// Letra maiúscula, no começø, indica que estou exportando a função ou struct
func (m Math) Add() int {
	return m.A + m.B
}

// deixando mathB não é exportada
type mathB struct {
	a int
	b int
}

func NewMathB(a, b int) mathB {
	return mathB{a: a, b: b}
}

func (m mathB) AddB() int {
	return m.a + m.b
}
