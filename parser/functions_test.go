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

	// Test for the Equality (==) Operator.
	got, err = ParseReader("", strings.NewReader("10 == 10"))
	expected = 1
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}

	// Test for the Inequality (!=) Operator.
	got, err = ParseReader("", strings.NewReader("10 != 10"))
	expected = 0
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}

	// Test for the "Smaller than" (<) Operator.
	got, err = ParseReader("", strings.NewReader("7 < 10"))
	expected = 1
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}

	// Test for the "Smaller than or equal to" (<=) Operator.
	got, err = ParseReader("", strings.NewReader("10 <= 10"))
	expected = 1
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}

	// Test for the "Greater than" (>) Operator.
	got, err = ParseReader("", strings.NewReader("10 > 7"))
	expected = 1
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}

	// Test for the "Greater than or equal to" (>=) Operator.
	got, err = ParseReader("", strings.NewReader("10 >= 10"))
	expected = 1
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}

	// Test for the Logical OR (||) Operator.
	got, err = ParseReader("", strings.NewReader("0 || 1"))
	expected = 1
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}

	// Test for the Logical AND (&&) Operator.
	got, err = ParseReader("", strings.NewReader("0 && 1"))
	expected = 0
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}

	// Test for the Left Shift (<<) Operator.
	got, err = ParseReader("", strings.NewReader("25 << 2"))
	expected = 100
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}

	// Test for the Left Shift (<<) Operator.
	got, err = ParseReader("", strings.NewReader("25 >> 2"))
	expected = 6
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}

	// Tests for complicated expressions.
	got, err = ParseReader("", strings.NewReader("25 >> 2 == 7 - 1"))
	expected = 1
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}

	got, err = ParseReader("", strings.NewReader("15 + 5 >= 10"))
	expected = 1
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}

	got, err = ParseReader("", strings.NewReader("14 * 10 + 5"))
	expected = 145
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}

	got, err = ParseReader("", strings.NewReader("2 * 10 // 4"))
	expected = 4
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}

	got, err = ParseReader("", strings.NewReader("45 | 25 == 61"))
	expected = 1
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}

	got, err = ParseReader("", strings.NewReader("45 | (25 == 61)"))
	expected = 45
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}

	got, err = ParseReader("", strings.NewReader("25 < 50 && 60 >= 30 || 10 != 10"))
	expected = 1
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		t.Errorf("Expected %v, got, %v\n", expected, got)
	}
}
