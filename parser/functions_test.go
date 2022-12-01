package parser

import (
	"log"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	// Test for the Modulus (//) Operator.
	got, err := ParseReader("", strings.NewReader("10 // 4"))
	expected := 2
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}

	// Test for the Bitwise And (&) Operator.
	got, err = ParseReader("", strings.NewReader("10 & 11"))
	expected = 10
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}

	// Test for the Bitwise Or (|) Operator.
	got, err = ParseReader("", strings.NewReader("25 | 10"))
	expected = 27
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}

	// Test for the Bitwise Xor (^) Operator.
	got, err = ParseReader("", strings.NewReader("45 ^ 35"))
	expected = 14
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}

	// Test for the Equality Operator (==) Operator.
	got, err = ParseReader("", strings.NewReader("10 == 10"))
	expected = 1
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}

	// Test for the Inequality Operator (!=) Operator.
	got, err = ParseReader("", strings.NewReader("10 != 10"))
	expected = 0
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}

	// Test for the "Smaller than" Operator (<) Operator.
	got, err = ParseReader("", strings.NewReader("7 < 10"))
	expected = 1
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}

	// Test for the "Smaller than or equal to" Operator (<=) Operator.
	got, err = ParseReader("", strings.NewReader("10 <= 10"))
	expected = 1
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}

	// Test for the "Greater than" Operator (>) Operator.
	got, err = ParseReader("", strings.NewReader("10 > 7"))
	expected = 1
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}

	// Test for the "Greater than or equal to" Operator (>=) Operator.
	got, err = ParseReader("", strings.NewReader("10 >= 10"))
	expected = 1
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}

	// Test for the Logical OR Operator (||) Operator.
	got, err = ParseReader("", strings.NewReader("0 || 1"))
	expected = 1
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}

	// Test for the Logical AND Operator (&&) Operator.
	got, err = ParseReader("", strings.NewReader("0 && 1"))
	expected = 0
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}
}
