package calev

import (
	"time"
)

// SpecOpt is a Spec configuration option.
// Most of them are constraints imposed on a Spec.
type SpecOpt func(s *Spec)

// Month adds one or more months to the month constraint.
func Month(months ...time.Month) SpecOpt {
	var monthSet uint16
	for _, month := range months {
		monthSet |= 1 << posMod(month-1, 12)
	}
	return func(s *Spec) {
		s.months |= monthSet
	}
}

// MonthRange adds months to the month constraint that fall in the specified range.
// The range starts at low, ends at high (inclusive, or maximum value if high is negative),
// and increments by step (or only includes low if step is 0).
func MonthRange(low, high time.Month, step int) SpecOpt {
	low = posMod(low-1, 12)
	high = posMod(negDefault(high, 12)-1, 12)
	high = max(low, high)
	step = zeroClamp(step, int(high-low))

	monthSet := rangeToSet[uint16](int(low), int(high), step)
	return func(s *Spec) {
		s.months |= monthSet
	}
}

// Day adds one or more days of month to the days of month constraint.
func Day(days ...int) SpecOpt {
	var daySet uint32
	for _, day := range days {
		daySet |= 1 << posMod(day-1, 31)
	}
	return func(s *Spec) {
		s.days |= daySet
	}
}

// DayRange adds days to the days of month constraint that fall in the specified range.
// The range starts at low, ends at high (inclusive, or maximum value if high is negative),
// and increments by step (or only includes low if step is 0).
func DayRange(low, high, step int) SpecOpt {
	low = posMod(low-1, 31)
	high = posMod(negDefault(high, 31)-1, 31)
	high = max(low, high)
	step = zeroClamp(step, high-low)

	daySet := rangeToSet[uint32](low, high, step)
	return func(s *Spec) {
		s.days |= daySet
	}
}

// LastDay adds one or more last days of month to the last days of month constraint.
func LastDay(lastDays ...int) SpecOpt {
	var lastDaySet uint32
	for _, lastDay := range lastDays {
		lastDaySet |= 1 << (30 - posMod(lastDay-1, 31))
	}
	return func(s *Spec) {
		s.ldays |= lastDaySet
	}
}

// LastDayRange adds last days of month to the last days of month constraint that fall in the specified range.
// The range starts at low, ends at high (inclusive, or maximum value if high is negative),
// and increments by step (or only includes low if step is 0).
func LastDayRange(low, high, step int) SpecOpt {
	low = posMod(low-1, 31)
	high = posMod(negDefault(high, 31)-1, 31)
	high = max(low, high)
	step = zeroClamp(step, high-low)

	lastDaySet := rangeToSetReverse[uint32](low, high, step, 30)
	return func(s *Spec) {
		s.ldays |= lastDaySet
	}
}

// Weekday adds one or more days of week to the days of week constraint.
func Weekday(weekdays ...time.Weekday) SpecOpt {
	var weekdaySet uint8
	for _, weekday := range weekdays {
		weekdaySet |= 1 << posMod(weekday, 7)
	}
	return func(s *Spec) {
		s.weekdays |= weekdaySet
	}
}

// WeekdayRange adds days of week to the days of week constraint that fall in the specified range.
// The range starts at low, ends at high (inclusive, or maximum value if high is negative),
// and increments by step (or only includes low if step is 0).
func WeekdayRange(low, high time.Weekday, step int) SpecOpt {
	low = posMod(low, 7)
	high = posMod(negDefault(high, 6), 7)
	high = max(low, high)
	step = zeroClamp(step, int(high-low))

	weekdaySet := rangeToSet[uint8](int(low), int(high), step)
	return func(s *Spec) {
		s.weekdays |= weekdaySet
	}
}

// WeekdayStrict adds one or more days of week to the restricted days of week constraint.
func WeekdayStrict(weekdays ...time.Weekday) SpecOpt {
	var weekdayStrictSet uint8
	for _, weekday := range weekdays {
		weekdayStrictSet |= 1 << posMod(weekday, 7)
	}
	return func(s *Spec) {
		s.weekdaysStrict |= weekdayStrictSet
	}
}

// WeekdayStrictRange adds days of week to the restricted days of week constraint that fall in the specified range.
// The range starts at low, ends at high (inclusive, or maximum value if high is negative),
// and increments by step (or only includes low if step is 0).
func WeekdayStrictRange(low, high time.Weekday, step int) SpecOpt {
	low = posMod(low, 7)
	high = posMod(negDefault(high, 6), 7)
	high = max(low, high)
	step = zeroClamp(step, int(high-low))

	weekdayStrictSet := rangeToSet[uint8](int(low), int(high), step)
	return func(s *Spec) {
		s.weekdaysStrict |= weekdayStrictSet
	}
}

// Hour adds one or more hours to the hours constraint.
func Hour(hours ...int) SpecOpt {
	var hourSet uint32
	for _, hour := range hours {
		hourSet |= 1 << posMod(hour, 24)
	}
	return func(s *Spec) {
		s.hours |= hourSet
	}
}

// HourRange adds hours to the hours constraint that fall in the specified range.
// The range starts at low, ends at high (inclusive, or maximum value if high is negative),
// and increments by step (or only includes low if step is 0).
func HourRange(low, high, step int) SpecOpt {
	low = posMod(low, 24)
	high = posMod(negDefault(high, 23), 24)
	high = max(low, high)
	step = zeroClamp(step, high-low)

	hourSet := rangeToSet[uint32](low, high, step)
	return func(s *Spec) {
		s.hours |= hourSet
	}
}

// Minute adds one or more minutes to the minutes constraint.
func Minute(minutes ...int) SpecOpt {
	var minuteSet uint64
	for _, minute := range minutes {
		minuteSet |= 1 << posMod(minute, 60)
	}
	return func(s *Spec) {
		s.minutes |= minuteSet
	}
}

// MinuteRange adds minutes to the minutes constraint that fall in the specified range.
// The range starts at low, ends at high (inclusive, or maximum value if high is negative),
// and increments by step (or only includes low if step is 0).
func MinuteRange(low, high, step int) SpecOpt {
	low = posMod(low, 60)
	high = posMod(negDefault(high, 59), 60)
	high = max(low, high)
	step = zeroClamp(step, high-low)

	minuteSet := rangeToSet[uint64](low, high, step)
	return func(s *Spec) {
		s.minutes |= minuteSet
	}
}

// Second adds one or more seconds to the seconds constraint.
func Second(seconds ...int) SpecOpt {
	var secondSet uint64
	for _, second := range seconds {
		secondSet |= 1 << posMod(second, 60)
	}
	return func(s *Spec) {
		s.seconds |= secondSet
	}
}

// SecondRange adds seconds to the seconds constraint that fall in the specified range.
// The range starts at low, ends at high (inclusive, or maximum value if high is negative),
// and increments by step (or only includes low if step is 0).
func SecondRange(low, high, step int) SpecOpt {
	low = posMod(low, 60)
	high = posMod(negDefault(high, 59), 60)
	high = max(low, high)
	step = zeroClamp(step, high-low)

	secondSet := rangeToSet[uint64](low, high, step)
	return func(s *Spec) {
		s.seconds |= secondSet
	}
}

func posMod[T ~int](value, mod T) T {
	if value < 0 {
		return 0
	}
	return value % mod
}

func zeroClamp[T ~int](value, max T) T {
	if value < 0 {
		return 0
	}
	if value > max {
		return max
	}
	return value
}

func negDefault[T ~int](value, def T) T {
	if value < 0 {
		return def
	}
	return value
}
