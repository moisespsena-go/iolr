package main

import (
	"fmt"

	"github.com/moisespsena-go/iolr"
)

var i int

func check(f func() (string, error)) {
	i++
	fmt.Println("--> EXAMPLE ", i, "<--")
	line, err := f()
	if err != nil {
		panic(err)
	}
	if iolr.IsEmptyInput(line) {
		fmt.Printf("Voce não digitou nada.\n\n")
	} else {
		fmt.Printf("Voce digitou: %q\n\n", line)
	}
}

func main() {
	l := iolr.STDMessageLR
	check(func() (string, error) {
		fmt.Println("Digite alguma coisa:")
		return iolr.StdinLR.ReadLineS()
	})
	check(func() (string, error) {
		return l.ReadS("Digite alguma coisa 2")
	})
	check(func() (string, error) {
		return l.ReadS("Digite alguma coisa OU apenas dê ENTER e veja o valor padrão", "Viva o Brasil!!")
	})
	msg := "Escolha uma Opção ou deixe vazio"
	check(func() (string, error) {
		return l.ReadFS(&iolr.FOptions{Message: msg, Options: []string{"a", "b", "c"}})
	})
	check(func() (string, error) {
		return l.ReadFS(&iolr.FOptions{Message: msg, Options: []string{"a", "b", "c"}, Default: "b"})
	})
	check(func() (string, error) {
		return l.ReadFS(&iolr.FOptions{Message: msg, Options: []string{"a", "b", "c"}}, "c")
	})
	check(func() (string, error) {
		return l.ReadFS(&iolr.FOptionsPairs{Message: msg, Options: iolr.MapToPairs(map[interface{}]string{"B": "Brazil", "E": "EUA"})})
	})
	check(func() (string, error) {
		return l.ReadFS(&iolr.FOptionsPairs{Message: msg, Options: iolr.MapToPairs(map[interface{}]string{"B": "Brazil", "E": "EUA"}), Default: "B"})
	})
	check(func() (string, error) {
		return l.ReadFS(&iolr.FOptionsPairs{Message: msg, Options: iolr.MapToPairs(map[interface{}]string{"B": "Brazil", "E": "EUA"}), Default: "B"}, "c")
	})

	msgObrigatorio := "Escolha uma Opção (OBS: NÃO pode ser vazio)"
	check(func() (string, error) {
		return l.RequireFS(&iolr.FOptions{Message: msgObrigatorio, Options: []string{"a", "b", "c"}})
	})
	check(func() (string, error) {
		return l.RequireFS(&iolr.FOptionsPairs{Message: msgObrigatorio, Options: iolr.MapToPairs(map[interface{}]string{"B": "Brazil", "E": "EUA"})})
	})
}
