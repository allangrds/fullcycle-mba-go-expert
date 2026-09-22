package main

import (
	"os"
	"text/template"
)

type Curso struct {
	Nome         string
	CargaHoraria int
}

type Cursos []Curso

func main() {
	cursos := Cursos{
		{
			Nome:         "Go",
			CargaHoraria: 40,
		},
		{
			Nome:         "Python",
			CargaHoraria: 30,
		},
	}

	//  Must é um método do pacote template que cria um novo template e trata erros de forma mais simples
	// Se ocorrer um erro, o programa irá panicar, o que é útil para evitar
	// verificações de erro repetitivas e deixar o código mais limpo.
	// ParseFiles é um método que permite carregar templates de arquivos externos.
	// Isso é útil para separar a lógica do template do código Go, tornando o código mais
	// organizado e fácil de manter.
	templat := template.Must(template.New("template.html").ParseFiles("template.html"))

	// O código abaixo executa o template e escreve a saída no os.Stdout
	err := templat.Execute(os.Stdout, cursos)
	if err != nil {
		panic(err)
	}
}
