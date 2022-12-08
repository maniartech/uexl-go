package main

import (
	"log"
	"strings"

	"github.com/maniartech/uexl_go/parser"
)

// for testing purpose

func main() {

	_, err := parser.ParseReader("", strings.NewReader("10 + 4"))
	if err != nil {
		log.Fatal(err)
	}

	// if len(os.Args) != 2 {
	// 	log.Fatal("Usage: calculator 'EXPR'")
	// }
	// got, err := parser.ParseReader("", strings.NewReader(os.Args[1]))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("=", got)
}
