package types

import "fmt"

// This file is the proleptic-Gregorian calendar engine behind the datetime standard library: converting
// a canonical millisecond instant to/from civil components, calendar-aware add/diff, and fixed-offset
// helpers. It is pure and dependency-free (datetime-spec §5, §11). daysFromCivil/daysInMonth/isLeapYear
// live in iso.go.

const msPerDay = 86400000

// floorDivMod returns q, r such that a = q*b + r with 0 <= r < b (Euclidean/floored division), so it is
// correct for negative instants (pre-epoch dates).
func floorDivMod(a, b int64) (q, r int64) {
	q = a / b
	r = a % b
	if r < 0 {
		q--
		r += b
	}
	return q, r
}

func floorDivInt(a, b int) int {
	q := a / b
	if a%b < 0 {
		q--
	}
	return q
}

func floorModInt(a, b int) int {
	r := a % b
	if r < 0 {
		r += b
	}
	return r
}

// civilFromDays converts days-since-1970-01-01 to a proleptic-Gregorian (year, month, day) — Howard
// Hinnant's algorithm, the inverse of daysFromCivil.
func civilFromDays(z int64) (year int, month int, day int) {
	z += 719468
	var era int64
	if z >= 0 {
		era = z / 146097
	} else {
		era = (z - 146096) / 146097
	}
	doe := z - era*146097                                  // [0, 146096]
	yoe := (doe - doe/1460 + doe/36524 - doe/146096) / 365 // [0, 399]
	y := yoe + era*400
	doy := doe - (365*yoe + yoe/4 - yoe/100) // [0, 365]
	mp := (5*doy + 2) / 153                  // [0, 11]
	d := doy - (153*mp+2)/5 + 1              // [1, 31]
	m := mp + 3
	if mp >= 10 {
		m = mp - 9
	}
	if m <= 2 {
		y++
	}
	return int(y), int(m), int(d)
}

// isoWeekday returns the ISO-8601 day of week for a day count: 1 = Monday … 7 = Sunday. 1970-01-01 was a
// Thursday (day 0).
func isoWeekday(days int64) int {
	w := floorModInt(int(days%7)+4, 7) // 0 = Sunday … 6 = Saturday
	if w == 0 {
		return 7
	}
	return w
}

// ComponentsOf decomposes a canonical UTC instant into civil components (the time-of-day fields use
// floored division so pre-epoch instants decompose correctly). weekday is ISO (1=Mon … 7=Sun).
func ComponentsOf(ms int64) (year, month, day, hour, minute, second, milli, weekday int) {
	days, rem := floorDivMod(ms, msPerDay)
	year, month, day = civilFromDays(days)
	hour = int(rem / 3600000)
	minute = int((rem % 3600000) / 60000)
	second = int((rem % 60000) / 1000)
	milli = int(rem % 1000)
	weekday = isoWeekday(days)
	return
}

func clampToRange(ms int64) int64 {
	if ms < MinDateTimeMillis {
		return MinDateTimeMillis
	}
	if ms > MaxDateTimeMillis {
		return MaxDateTimeMillis
	}
	return ms
}

// MillisFromComponents builds a canonical instant from civil components. Structurally invalid fields
// (month, day-for-month, hour, minute, second, millisecond) are an error; a year outside 0001..9999 is
// clamped to the nearest boundary instant (datetime-spec §5.4).
func MillisFromComponents(year, month, day, hour, minute, second, milli int) (int64, error) {
	switch {
	case month < 1 || month > 12:
		return 0, fmt.Errorf("month %d out of range 1..12", month)
	case day < 1 || day > daysInMonth(year, month):
		return 0, fmt.Errorf("day %d out of range for %04d-%02d", day, year, month)
	case hour < 0 || hour > 23:
		return 0, fmt.Errorf("hour %d out of range 0..23", hour)
	case minute < 0 || minute > 59:
		return 0, fmt.Errorf("minute %d out of range 0..59", minute)
	case second < 0 || second > 59:
		return 0, fmt.Errorf("second %d out of range 0..59", second)
	case milli < 0 || milli > 999:
		return 0, fmt.Errorf("millisecond %d out of range 0..999", milli)
	}
	if year < 1 {
		return MinDateTimeMillis, nil
	}
	if year > 9999 {
		return MaxDateTimeMillis, nil
	}
	days := daysFromCivil(int64(year), month, day)
	ms := days*msPerDay + int64(hour)*3600000 + int64(minute)*60000 + int64(second)*1000 + int64(milli)
	return clampToRange(ms), nil
}

// AddMonths shifts an instant by a (signed) whole number of months, preserving time-of-day and clamping
// the day to the last valid day of the target month (end-of-month clamp, datetime-spec §5.3). The result
// is clamped to the portable range.
func AddMonths(ms int64, n int) int64 {
	days, rem := floorDivMod(ms, msPerDay)
	y, m, d := civilFromDays(days)
	total := y*12 + (m - 1) + n
	ny := floorDivInt(total, 12)
	nm := floorModInt(total, 12) + 1
	if maxD := daysInMonth(ny, nm); d > maxD {
		d = maxD
	}
	return clampToRange(daysFromCivil(int64(ny), nm, d)*msPerDay + rem)
}

// AddYears shifts an instant by a (signed) whole number of years (= AddMonths by 12n), with the same
// end-of-month clamp (e.g. 2024-02-29 + 1y -> 2025-02-28).
func AddYears(ms int64, n int) int64 {
	return AddMonths(ms, n*12)
}

// DiffMonths returns the signed count of complete calendar months from a to b, truncated toward zero —
// the inverse of AddMonths (datetime-spec §5.3): the largest k>=0 with AddMonths(a,k) <= b, negated when
// b precedes a.
func DiffMonths(a, b int64) int {
	if a == b {
		return 0
	}
	sign, lo, hi := 1, a, b
	if a > b {
		sign, lo, hi = -1, b, a
	}
	yl, ml, _ := civilFromDays(floorDiv64(lo, msPerDay))
	yh, mh, _ := civilFromDays(floorDiv64(hi, msPerDay))
	months := (yh-yl)*12 + (mh - ml)
	if AddMonths(lo, months) > hi {
		months--
	}
	return sign * months
}

// DiffYears returns the signed count of complete calendar years from a to b, truncated toward zero (the
// inverse of AddYears, datetime-spec §5.3).
func DiffYears(a, b int64) int {
	if a == b {
		return 0
	}
	sign, lo, hi := 1, a, b
	if a > b {
		sign, lo, hi = -1, b, a
	}
	yl, _, _ := civilFromDays(floorDiv64(lo, msPerDay))
	yh, _, _ := civilFromDays(floorDiv64(hi, msPerDay))
	years := yh - yl
	if AddYears(lo, years) > hi {
		years--
	}
	return sign * years
}

func floorDiv64(a, b int64) int64 {
	q, _ := floorDivMod(a, b)
	return q
}

// ParseFixedOffset parses a fixed ISO 8601 offset ("Z" or "±HH:MM") to signed minutes east of UTC; it is
// the exported form used by datePart/formatDate offset arguments (datetime-spec §10.4).
func ParseFixedOffset(s string) (int, bool) {
	return parseOffsetMinutes(s)
}
