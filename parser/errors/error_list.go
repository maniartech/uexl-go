package errors

import (
	"fmt"
	"sort"
	"strings"
)

// ErrorList represents a list of parsing errors, similar to scanner.ErrorList
// This follows the Go standard library pattern for error accumulation
type ErrorList []ParserError

// Len returns the number of errors in the list
func (p ErrorList) Len() int { return len(p) }

// Swap swaps two errors in the list
func (p ErrorList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// Less compares two errors by position for sorting
func (p ErrorList) Less(i, j int) bool {
	e := &p[i]
	f := &p[j]

	// Compare by line first
	if e.Line < f.Line {
		return true
	}
	if e.Line > f.Line {
		return false
	}

	// If lines are equal, compare by column
	return e.Column < f.Column
}

// Sort sorts the error list by position
func (p ErrorList) Sort() {
	sort.Sort(p)
}

// RemoveMultiples removes duplicate errors on the same position
func (p *ErrorList) RemoveMultiples() {
	if len(*p) <= 1 {
		return
	}

	sort.Sort(*p)

	i := 0
	for j := 1; j < len(*p); j++ {
		if (*p)[i].Line != (*p)[j].Line || (*p)[i].Column != (*p)[j].Column {
			i++
			(*p)[i] = (*p)[j]
		}
	}
	*p = (*p)[:i+1]
}

// Error returns a string representation of the error list
// This makes ErrorList implement the error interface
func (p ErrorList) Error() string {
	switch len(p) {
	case 0:
		return "no errors"
	case 1:
		return p[0].Error()
	}

	var buf strings.Builder
	fmt.Fprintf(&buf, "%d errors:", len(p))
	for _, err := range p {
		fmt.Fprintf(&buf, "\n\t%s", err.Error())
	}
	return buf.String()
}

// Err returns an error equivalent to this error list.
// If the list is empty, Err returns nil.
func (p ErrorList) Err() error {
	if len(p) == 0 {
		return nil
	}
	return p
}

// Add adds a new error to the list at the specified position
func (p *ErrorList) Add(line, column int, msg string) {
	*p = append(*p, NewParserError(ErrGeneric, line, column, msg))
}

// AddError adds a ParserError to the list
func (p *ErrorList) AddError(err ParserError) {
	*p = append(*p, err)
}

// Reset clears the error list
func (p *ErrorList) Reset() {
	*p = (*p)[:0]
}
