package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/lautarojayat/writing-an-interpreter-in-go/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("hi %s! wellcome to the REPL!\n", user.Name)
	repl.Start(os.Stdin, os.Stdout)

}
