package main

import (
	"fmt"
)

type tsv struct {
}

var TSV tsv

func (t tsv) Greet() {
	fmt.Println("Hello Universe")
}

func (t tsv) Execute() error {
	fmt.Printf("WHOA! TSV plugin here we go\n")

	return nil
}
