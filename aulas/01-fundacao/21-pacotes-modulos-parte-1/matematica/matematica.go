package matematica

//Se estiver em maiscula, o pacote é exportado
//Se estiver em minuscula, o pacote é privado
func Soma[T int | float64](a, b T) T {
  return a + b
}

var A int = 10
var b int = 10
