package math

type Math struct {
	A int
	B int
}

// func Add(a, b int) int { - mesma coisa que abaixo, mas sem a struct
func (m Math) Add() int {
	return m.A + m.B
}
