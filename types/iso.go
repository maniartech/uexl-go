package types

import "fmt"

// Guaranteed portable range for datetime values (datetime-spec §5.2): 0001-01-01T00:00:00.000Z ..
// 9999-12-31T23:59:59.999Z, expressed as canonical milliseconds since the Unix epoch.
const (
	MinDateTimeMillis int64 = -62135596800000
	MaxDateTimeMillis int64 = 253402300799999

	// MaxDurationMillis bounds duration literals to the largest span between two in-range datetimes, so
	// the int64<->float64 storage (Phase A) stays lossless (it is < 2^53). datetime-spec §2, §3.3.
	MaxDurationMillis int64 = MaxDateTimeMillis - MinDateTimeMillis
)

// ParseISODateTime parses the accepted ISO 8601 / RFC 3339 subset (datetime-spec §4) and returns the
// canonical milliseconds since the Unix epoch (UTC). It is self-contained (no dependencies) and backs the
// d"..." literal; the datetime library reuses it.
//
// Accepted forms (§4.1): "YYYY-MM-DD", optionally "THH:MM[:SS[.frac]]", optionally an offset ("Z" or
// "±HH:MM"). Missing time defaults to midnight and a missing offset to UTC; a present offset is applied
// to compute the UTC instant and then discarded (§4.2). Sub-millisecond fractions are truncated toward
// zero (§4.3). Invalid fields and out-of-range instants return an error (§5.2, §5.3).
func ParseISODateTime(s string) (int64, error) {
	// Date: exactly YYYY-MM-DD.
	if len(s) < 10 || s[4] != '-' || s[7] != '-' {
		return 0, invalidISO(s, "expected YYYY-MM-DD date")
	}
	year, ok1 := parseDigitsN(s[0:4])
	month, ok2 := parseDigitsN(s[5:7])
	day, ok3 := parseDigitsN(s[8:10])
	if !ok1 || !ok2 || !ok3 {
		return 0, invalidISO(s, "non-numeric date field")
	}

	hour, minute, second, milli := 0, 0, 0, 0
	offsetMinutes := 0 // signed; applied to compute UTC, then discarded

	if rest := s[10:]; len(rest) > 0 {
		// Time: T HH:MM [:SS [.frac]] [offset]
		if rest[0] != 'T' || len(rest) < 6 || rest[3] != ':' {
			return 0, invalidISO(s, "expected 'T' followed by HH:MM")
		}
		var okh, okm bool
		hour, okh = parseDigitsN(rest[1:3])
		minute, okm = parseDigitsN(rest[4:6])
		if !okh || !okm {
			return 0, invalidISO(s, "non-numeric time field")
		}
		idx := 6
		if idx < len(rest) && rest[idx] == ':' { // optional seconds
			if idx+3 > len(rest) {
				return 0, invalidISO(s, "incomplete seconds field")
			}
			var oks bool
			second, oks = parseDigitsN(rest[idx+1 : idx+3])
			if !oks {
				return 0, invalidISO(s, "non-numeric seconds field")
			}
			idx += 3
			if idx < len(rest) && rest[idx] == '.' { // optional fractional seconds
				idx++
				start := idx
				for idx < len(rest) && rest[idx] >= '0' && rest[idx] <= '9' {
					idx++
				}
				if idx == start {
					return 0, invalidISO(s, "empty fractional seconds")
				}
				milli = fracToMillis(rest[start:idx]) // truncated toward zero (§4.3)
			}
		}
		if idx < len(rest) { // optional offset, consumes the remainder
			off, oko := parseOffsetMinutes(rest[idx:])
			if !oko {
				return 0, invalidISO(s, "invalid offset")
			}
			offsetMinutes = off
			idx = len(rest)
		}
		if idx != len(rest) {
			return 0, invalidISO(s, "trailing characters")
		}
	}

	// Field validation (§5.2, §5.3). Year is grammatically 4 digits, so 0000 / out-of-range fail here.
	switch {
	case year < 1 || year > 9999:
		return 0, invalidISO(s, "year out of range 0001..9999")
	case month < 1 || month > 12:
		return 0, invalidISO(s, "month out of range 1..12")
	case day < 1 || day > daysInMonth(year, month):
		return 0, invalidISO(s, "day out of range for month")
	case hour > 23:
		return 0, invalidISO(s, "hour out of range 0..23")
	case minute > 59:
		return 0, invalidISO(s, "minute out of range 0..59")
	case second > 59: // 60 = leap second, rejected (§5.3)
		return 0, invalidISO(s, "second out of range 0..59")
	}

	days := daysFromCivil(int64(year), month, day)
	ms := days*86400000 + int64(hour)*3600000 + int64(minute)*60000 + int64(second)*1000 + int64(milli)
	ms -= int64(offsetMinutes) * 60000 // UTC = local - offset

	if ms < MinDateTimeMillis || ms > MaxDateTimeMillis {
		return 0, invalidISO(s, "instant out of portable range 0001..9999")
	}
	return ms, nil
}

func invalidISO(s, reason string) error {
	return fmt.Errorf("invalid datetime %q: %s", s, reason)
}

// parseDigitsN parses a slice that must be entirely ASCII digits, returning its integer value.
func parseDigitsN(s string) (int, bool) {
	if len(s) == 0 {
		return 0, false
	}
	n := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return 0, false
		}
		n = n*10 + int(c-'0')
	}
	return n, true
}

// fracToMillis converts fractional-second digits to whole milliseconds, truncated toward zero (§4.3):
// it reads the first three digits, right-padding with zeros and discarding any beyond the third.
func fracToMillis(frac string) int {
	ms := 0
	for i := 0; i < 3; i++ {
		ms *= 10
		if i < len(frac) {
			ms += int(frac[i] - '0')
		}
	}
	return ms
}

// parseOffsetMinutes parses "Z" or "±HH:MM" into signed minutes east of UTC.
func parseOffsetMinutes(s string) (int, bool) {
	if s == "Z" {
		return 0, true
	}
	if len(s) != 6 || (s[0] != '+' && s[0] != '-') || s[3] != ':' {
		return 0, false
	}
	hh, ok1 := parseDigitsN(s[1:3])
	mm, ok2 := parseDigitsN(s[4:6])
	if !ok1 || !ok2 || hh > 23 || mm > 59 {
		return 0, false
	}
	total := hh*60 + mm
	if s[0] == '-' {
		return -total, true
	}
	return total, true
}

func isLeapYear(y int) bool {
	return (y%4 == 0 && y%100 != 0) || y%400 == 0
}

func daysInMonth(year, month int) int {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		if isLeapYear(year) {
			return 29
		}
		return 28
	default:
		return 0
	}
}

// daysFromCivil returns days since 1970-01-01 for a proleptic-Gregorian date (Howard Hinnant's algorithm).
// month is 1..12 and day is 1..31; the caller validates ranges beforehand.
func daysFromCivil(year int64, month, day int) int64 {
	y := year
	if month <= 2 {
		y--
	}
	var era int64
	if y >= 0 {
		era = y / 400
	} else {
		era = (y - 399) / 400
	}
	yoe := y - era*400 // [0, 399]
	mp := int64(month)
	if month > 2 {
		mp -= 3
	} else {
		mp += 9
	}
	doy := (153*mp+2)/5 + int64(day) - 1   // [0, 365]
	doe := yoe*365 + yoe/4 - yoe/100 + doy // [0, 146096]
	return era*146097 + doe - 719468
}
