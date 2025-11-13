package main

import (
	"fmt"

	"github.com/maniartech/uexl/parser"
)

func demonstratePanicFixes() {
	fmt.Println("=== Demonstrating Parser Error Handling Improvements ===")

	testCases := []string{
		"",                    // empty string
		"r\"unterminated",     // unterminated raw string
		"\"unterminated",      // unterminated regular string
		"1 + 2",               // valid expression
		"@invalid",            // invalid character
		"func(a, b) + [1, 2]", // valid complex expression
	}

	for _, input := range testCases {
		fmt.Printf("Testing input: %q\n", input)

		// Test old approach (for backward compatibility)
		fmt.Println("  Old NewParser approach:")
		parser1 := parser.NewParser(input)
		result1, err1 := parser1.Parse()
		if err1 != nil {
			fmt.Printf("    Error: %v\n", err1)
		} else {
			fmt.Printf("    Success: %T\n", result1)
		}

		// Test new approach (with validation)
		fmt.Println("  New NewParserWithValidation approach:")
		parser2, err2 := parser.NewParserWithValidation(input)
		if err2 != nil {
			fmt.Printf("    Validation Error: %v\n", err2)
		} else {
			result2, parseErr := parser2.Parse()
			if parseErr != nil {
				fmt.Printf("    Parse Error: %v\n", parseErr)
			} else {
				fmt.Printf("    Success: %T\n", result2)
			}
		}

		fmt.Println()
	}
}
