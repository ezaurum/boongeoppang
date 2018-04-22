package boongeoppang

import "time"

func asDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func asDate12HMinute(t time.Time) string {
	return t.Format("2006-01-02 03:04 PM")
}

func asDate24HMinute(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

func asTime24H(t time.Time) string {
	return t.Format("15:04")
}

func asTime12H(t time.Time) string {
	return t.Format("PM 03:04")
}
