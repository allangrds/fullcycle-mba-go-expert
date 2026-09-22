package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
  //Criacao de arquivo
  file, err := os.Create("meu_arquivo.txt")
  if err != nil {
    panic(err)
  }

  // Escrevendo uma string no arquivo
  // tamanho, err := file.WriteString("Olá, mundo!")
  // Escrevendo bytes diretamente no arquivo
  tamanho, err := file.Write([]byte("Olá, mundo!"))
  if err != nil {
    panic(err)
  }

  fmt.Printf("Escrevi %d bytes no arquivo.\n", tamanho)

  //leitura
  // arquivo, err := os.Open("meu_arquivo.txt")
  arquivo, err := os.ReadFile("meu_arquivo.txt")
  if err != nil {
    panic(err)
  }
  fmt.Printf("Conteúdo do arquivo: %s\n", string(arquivo))

  //leitura de pouco em pouco abrindo o arquivo
  arquivo2, err := os.Open("meu_arquivo.txt")
  if err != nil {
    panic(err)
  }
  reader := bufio.NewReader(arquivo2)
  buffer := make([]byte, 2) // Lê 2 bytes por vez

  for {
    n, err := reader.Read(buffer)
    if err != nil {
      break
    }

    fmt.Println(string(buffer[:n]))
  }

  err = os.Remove("meu_arquivo.txt")
  if err != nil {
    panic(err)
  }

  file.Close()
}
