package main

import (
	"fmt"

	"github.com/moisespsena/go-ioutil"
)

var i int

func check(f func() (string, error)) {
	i++
	fmt.Println("--> EXAMPLE ", i, "<--")
	line, err := f()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Voce digitou: %q\n\n", line)
}

func main() {
	check(func() (string, error) {
		fmt.Println("Digite alguma coisa:")
		return ioutil.StdinLR.ReadLineS()
	})
	l := ioutil.STDStringMessageLR
	check(func() (string, error) {
		return l.Read("Digite alguma coisa 2")
	})
	check(func() (string, error) {
		return l.Read("Digite alguma coisa OU apenas dê ENTER e veja o valor padrão", "Viva o Brasil!!")
	})

	msg := "Escolha uma Opção ou deixe vazio"
	check(func() (string, error) {
		return l.ReadF(&ioutil.FOptions{Message: msg, Options: []string{"a", "b", "c"}})
	})
	check(func() (string, error) {
		return l.ReadF(&ioutil.FOptions{Message: msg, Options: []string{"a", "b", "c"}, Default: "b"})
	})
	check(func() (string, error) {
		return l.ReadF(&ioutil.FOptions{Message: msg, Options: []string{"a", "b", "c"}}, "c")
	})
	check(func() (string, error) {
		return l.ReadF(&ioutil.FOptionsMap{Message: msg, Options: map[string]string{"B": "Brazil", "E": "EUA"}})
	})
	check(func() (string, error) {
		return l.ReadF(&ioutil.FOptionsMap{Message: msg, Options: map[string]string{"B": "Brazil", "E": "EUA"}, Default: "B"})
	})
	check(func() (string, error) {
		return l.ReadF(&ioutil.FOptionsMap{Message: msg, Options: map[string]string{"B": "Brazil", "E": "EUA"}, Default: "B"}, "c")
	})

	msgObrigatorio := "Escolha uma Opção (OBS: NÃO pode ser vazio)"
	check(func() (string, error) {
		return l.RequireF(&ioutil.FOptions{Message: msgObrigatorio, Options: []string{"a", "b", "c"}})
	})
	check(func() (string, error) {
		return l.RequireF(&ioutil.FOptionsMap{Message: msgObrigatorio, Options: map[string]string{"B": "Brazil", "E": "EUA"}})
	})
}
