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

const DateLayout = "2006-01-02"

func SprintRange(anchor, t time.Time, days int) (time.Time, time.Time) {
	anchor = StartOfDay(anchor)
	diff := daysBetween(anchor, t)
	n := diff / days
	if diff < 0 && diff%days != 0 {
		n--
	}
	start := anchor.AddDate(0, 0, n*days)
	return start, start.AddDate(0, 0, days)
}

func daysBetween(a, b time.Time) int {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	ua := time.Date(ay, am, ad, 0, 0, 0, 0, time.UTC)
	ub := time.Date(by, bm, bd, 0, 0, 0, 0, time.UTC)
	return int(ub.Sub(ua).Hours() / 24)
}
