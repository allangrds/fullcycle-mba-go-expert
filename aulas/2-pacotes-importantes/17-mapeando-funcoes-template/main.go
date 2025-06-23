package main

import (
	"html/template" //html/template sabe que o pacote template é usado para gerar HTML, por isso ela vai ter uns blindagens
	"os"
	"strings"
)

type Curso struct {
	Nome         string
	CargaHoraria int
}

type Cursos []Curso

func ToUpper(s string) string {
	return strings.ToUpper(s)
}

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
	templates := []string{
		"header.html",
		"content.html",
		"footer.html",
	}

	templat := template.Must(
		// template.New("content.html").Funcs(template.FuncMap{"ToUpper": strings.ToUpper}).ParseFiles(templates...))
		template.New("content.html").Funcs(template.FuncMap{"ToUpper": ToUpper}).ParseFiles(templates...))
	err := templat.Execute(os.Stdout, cursos)
	if err != nil {
		panic(err)
	}
}
