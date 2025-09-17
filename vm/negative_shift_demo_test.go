package vm_test

import (
	"testing"
)

// TestNegativeShiftPanicDemo demonstrates the critical issue with negative shift amounts
// This test is expected to PANIC with current implementation - it's here to demonstrate the issue
func TestNegativeShiftPanicDemo(t *testing.T) {
	t.Skip("DEMO TEST: This would panic with current implementation - run manually to see the issue")
	
	// WARNING: This will cause a runtime panic!
	// Expected error: "panic: runtime error: negative shift amount"
	tests := []vmTestCase{
		{"8 << -1", 0.0}, // This will PANIC, not return an error gracefully
	}
	
	// If we reach this line, the implementation has been fixed
	runVmTests(t, tests)
}