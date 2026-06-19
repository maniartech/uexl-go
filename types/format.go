package types

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// This file implements datetime/duration rendering and interchange for the datetime standard library:
// NITES datetime formatting (datetime-spec §10.1) and ISO 8601 duration format/parse (§10.3). All
// English/invariant (the core, locale-independent profile). Pure, dependency-free.

var monthsShort = [...]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
var monthsFull = [...]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
var weekdaysShort = [...]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
var weekdaysFull = [...]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

// namedLayouts are the NITES named layouts recognized by the core profile. "iso" is the default.
var namedLayouts = map[string]string{
	"iso":     "yyyy-mm-ddThhh:ii:ss",
	"isodate": "yyyy-mm-dd",
	"date":    "yyyy-mm-dd",
	"time":    "hhh:ii:ss",
	"sql":     "yyyy-mm-dd hhh:ii:ss",
	"rfc":     "www, dd mmm yyyy hhh:ii:ss",
}

// DefaultLayout is the layout used by formatDate when no pattern is supplied.
const DefaultLayout = "iso"

func nitesLetter(c byte) bool {
	switch c {
	case 'y', 'm', 'd', 'h', 'i', 's', 'a', 'w', 'o', 'z':
		return true
	}
	return false
}

func lowerASCII(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + 32
	}
	return c
}

// FormatNITES renders a canonical instant as seen at offsetMinutes east of UTC, using a NITES pattern or
// named layout (case-insensitive). Name-based fields render English/invariant names (§10.1).
func FormatNITES(msUTC int64, offsetMinutes int, pattern string) string {
	if exp, ok := namedLayouts[strings.ToLower(strings.TrimSpace(pattern))]; ok {
		pattern = exp
	}
	local := msUTC + int64(offsetMinutes)*60000
	y, mo, d, h, mi, se, _, wd := ComponentsOf(local)

	var b strings.Builder
	for i := 0; i < len(pattern); {
		lc := lowerASCII(pattern[i])
		if !nitesLetter(lc) {
			b.WriteByte(pattern[i])
			i++
			continue
		}
		j := i + 1
		for j < len(pattern) && lowerASCII(pattern[j]) == lc {
			j++
		}
		b.WriteString(nitesField(lc, j-i, y, mo, d, h, mi, se, wd, offsetMinutes))
		i = j
	}
	return b.String()
}

func nitesField(letter byte, n, y, mo, d, h, mi, se, wd, offMin int) string {
	switch letter {
	case 'y':
		switch {
		case n >= 4:
			return fmt.Sprintf("%04d", y)
		case n == 2:
			return fmt.Sprintf("%02d", y%100)
		default:
			return strconv.Itoa(y)
		}
	case 'm':
		switch {
		case n >= 4:
			return monthsFull[mo-1]
		case n == 3:
			return monthsShort[mo-1]
		case n == 2:
			return fmt.Sprintf("%02d", mo)
		default:
			return strconv.Itoa(mo)
		}
	case 'd':
		if n >= 2 {
			return fmt.Sprintf("%02d", d)
		}
		return strconv.Itoa(d)
	case 'h':
		if n >= 3 { // hhh = 24-hour
			return fmt.Sprintf("%02d", h)
		}
		h12 := h % 12
		if h12 == 0 {
			h12 = 12
		}
		if n == 2 {
			return fmt.Sprintf("%02d", h12)
		}
		return strconv.Itoa(h12)
	case 'i':
		if n >= 2 {
			return fmt.Sprintf("%02d", mi)
		}
		return strconv.Itoa(mi)
	case 's':
		if n >= 2 {
			return fmt.Sprintf("%02d", se)
		}
		return strconv.Itoa(se)
	case 'a':
		up := h < 12
		if n >= 2 {
			if up {
				return "AM"
			}
			return "PM"
		}
		if up {
			return "am"
		}
		return "pm"
	case 'w':
		if n >= 4 {
			return weekdaysFull[wd-1]
		}
		return weekdaysShort[wd-1]
	case 'z':
		if offMin == 0 {
			return "Z"
		}
		return formatOffset(offMin, true)
	case 'o':
		colon := n >= 3
		return formatOffset(offMin, colon)
	}
	return ""
}

func formatOffset(offMin int, colon bool) string {
	sign := "+"
	if offMin < 0 {
		sign = "-"
		offMin = -offMin
	}
	if colon {
		return fmt.Sprintf("%s%02d:%02d", sign, offMin/60, offMin%60)
	}
	return fmt.Sprintf("%s%02d%02d", sign, offMin/60, offMin%60)
}

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
