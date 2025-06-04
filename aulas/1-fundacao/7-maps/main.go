package main

import "fmt"

func main() {
  salarios := map[string]int{
    "João": 1500,
  }
  // fmt.Println("Salário de João:", salarios["João"])

  delete(salarios, "João")
  salarios["Maria"] = 2000
  // fmt.Println("Salário de Maria:", salarios["Maria"])

  salarios2 := make(map[string]int)
  salarios2["Pedro"] = 2500
  // fmt.Println("Salário de Pedro:", salarios2["Pedro"])

  //mostrar chave e valor
  // for nome, salario := range salarios {
  //   fmt.Printf("Salário de %s: %d\n", nome, salario)
  // }

  //apenas mostrar valor
  for _, salario := range salarios {
    fmt.Printf("Salário de %d\n", salario)
  }
}
