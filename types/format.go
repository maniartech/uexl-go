package types

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// This file implements duration interchange for the datetime standard library: ISO 8601 duration
// format/parse (datetime-spec §10.3) — English/invariant, pure, and dependency-free. DateTime (NITES)
// formatting and pattern-directed parsing live in the datetime standard library (builtins/datetime),
// backed by the gotime reference implementation; they are not part of the dependency-free core.

// FormatISODuration renders a duration as an ISO 8601 duration using only the exact components
// (week/day/hour/minute/second), e.g. "PT1H30M", "P7D", "PT0.5S" (datetime-spec §10.3). Negative spans
// are prefixed with "-".
func FormatISODuration(ms int64) string {
	if ms == 0 {
		return "PT0S"
	}
	var b strings.Builder
	if ms < 0 {
		b.WriteByte('-')
		ms = -ms
	}
	b.WriteByte('P')
	days := ms / msPerDay
	ms %= msPerDay
	if days > 0 {
		fmt.Fprintf(&b, "%dD", days)
	}
	hours := ms / 3600000
	ms %= 3600000
	minutes := ms / 60000
	ms %= 60000
	seconds := ms / 1000
	millis := ms % 1000
	if hours > 0 || minutes > 0 || seconds > 0 || millis > 0 {
		b.WriteByte('T')
		if hours > 0 {
			fmt.Fprintf(&b, "%dH", hours)
		}
		if minutes > 0 {
			fmt.Fprintf(&b, "%dM", minutes)
		}
		if seconds > 0 || millis > 0 {
			if millis > 0 {
				frac := strings.TrimRight(fmt.Sprintf("%03d", millis), "0")
				fmt.Fprintf(&b, "%d.%sS", seconds, frac)
			} else {
				fmt.Fprintf(&b, "%dS", seconds)
			}
		}
	}
	return b.String()
}

// ParseISODuration parses an ISO 8601 duration restricted to the exact components (W/D and, after T,
// H/M/S); the year designator and the date-part month designator are rejected (datetime-spec §10.3). The
// result is bounded to ±MaxDurationMillis.
func ParseISODuration(s string) (int64, error) {
	neg := false
	if strings.HasPrefix(s, "-") {
		neg, s = true, s[1:]
	} else if strings.HasPrefix(s, "+") {
		s = s[1:]
	}
	if len(s) == 0 || s[0] != 'P' {
		return 0, fmt.Errorf("invalid ISO 8601 duration %q: must start with P", s)
	}
	s = s[1:]
	inTime := false
	seen := false
	total := 0.0
	i := 0
	for i < len(s) {
		if s[i] == 'T' {
			inTime = true
			i++
			continue
		}
		start := i
		for i < len(s) && ((s[i] >= '0' && s[i] <= '9') || s[i] == '.') {
			i++
		}
		if i == start || i >= len(s) {
			return 0, fmt.Errorf("invalid ISO 8601 duration %q", s)
		}
		val, err := strconv.ParseFloat(s[start:i], 64)
		if err != nil || math.IsNaN(val) || math.IsInf(val, 0) {
			return 0, fmt.Errorf("invalid ISO 8601 duration component %q", s[start:i])
		}
		des := s[i]
		i++
		var unit float64
		switch des {
		case 'W':
			if inTime {
				return 0, fmt.Errorf("W must precede T in %q", s)
			}
			unit = 604800000
		case 'D':
			if inTime {
				return 0, fmt.Errorf("D must precede T in %q", s)
			}
			unit = 86400000
		case 'H':
			if !inTime {
				return 0, fmt.Errorf("H must follow T in %q", s)
			}
			unit = 3600000
		case 'M':
			if !inTime {
				return 0, fmt.Errorf("month designator is not an exact duration in %q", s)
			}
			unit = 60000
		case 'S':
			if !inTime {
				return 0, fmt.Errorf("S must follow T in %q", s)
			}
			unit = 1000
		case 'Y':
			return 0, fmt.Errorf("year designator is not an exact duration in %q", s)
		default:
			return 0, fmt.Errorf("invalid ISO 8601 duration designator %q", string(des))
		}
		total += val * unit
		seen = true
	}
	if !seen {
		return 0, fmt.Errorf("empty ISO 8601 duration %q", s)
	}
	if neg {
		total = -total
	}
	if total > float64(MaxDurationMillis) || total < float64(-MaxDurationMillis) {
		return 0, fmt.Errorf("ISO 8601 duration %q is out of range", s)
	}
	return int64(total), nil
}
