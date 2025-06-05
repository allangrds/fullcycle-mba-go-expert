package main

func SomaInteiro(varMap map[string]int) int {
  var soma int = 0

  for _, value := range varMap {
    soma += value
  }

  return soma
}

func SomaFloat(varMap map[string]float64) float64 {
  var soma float64 = 0

  for _, value := range varMap {
    soma += value
  }

  return soma
}

func Soma[T int | float64] (varMap map[string]T) T {
  var soma T = 0

  for _, value := range varMap {
    soma += value
  }

  return soma
}

type MyNumber int

type Number interface {
  // ~int não apenas int de forma direta, mas int indireto, como o MyNumber
  ~int | ~float64
}

func Soma2[T Number] (varMap map[string]T) T {
  var soma T = 0

  for _, value := range varMap {
    soma += value
  }

  return soma
}

func Compara[T Number](a T, b T) bool {
  if a == b {
    return true
  }

  return false
}

//Comparable compara a igualdade, não pode usar a > b
func Compara2[T comparable](a T, b T) bool {
  if a == b {
    return true
  }

  return false
}

// func Compara3[T comparable](a T, b T) bool {
//   if a > b {
//     return true
//   }

//   return false
// }

func main() {
  varMap := map[string]int{
    "João": 1000,
    "Maria": 2000,
    "Pedro": 3000,
  }
  varMapFloat := map[string]float64{
    "João": 1000.50,
    "Maria": 2000.50,
    "Pedro": 3000.50,
  }
  varMapMyNumber := map[string]MyNumber{
    "João": 1000,
    "Maria": 2000,
    "Pedro": 3000,
  }

  println("A soma dos valores do mapa é:", SomaInteiro(varMap))
  println("A soma dos valores do mapa float é:", SomaFloat(varMapFloat))
  println("A soma dos valores do mapa é:", Soma(varMap))
  println("A soma dos valores do mapa float é:", Soma(varMapFloat))
  println("A soma dos valores do mapa é:", Soma2(varMap))
  println("A soma dos valores do mapa float é:", Soma2(varMapFloat))
  println("A soma dos valores do mapa é:", Soma2(varMapMyNumber))
  println("Comparando 10 e 10:", Compara(10, 10))
  println("Comparando 10 e 10.0:", Compara(10, 10.0))
  println("Comparando2 10 e 10.0:", Compara2(10, 10.0))
  // println("Comparando3 10 e 10.0:", Compara3(10, 10.0))
}
