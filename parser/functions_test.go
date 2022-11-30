package parser

import (
	"fmt"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	got, err := ParseReader("", strings.NewReader("10 @ 4"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(got)
}
