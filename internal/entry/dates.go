package entry

import "time"

func StartOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func DayRange(t time.Time) (time.Time, time.Time) {
	start := StartOfDay(t)
	end := start.AddDate(0, 0, 1)
	return start, end
}

func PreviousWorkday(t time.Time) time.Time {
	days := 1
	switch t.Weekday() {
	case time.Monday:
		days = 3
	case time.Sunday:
		days = 2
	}
	return t.AddDate(0, 0, -days)
}
