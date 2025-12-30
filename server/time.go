package main

import (
	"fmt"

	"cloud.google.com/go/civil"
	"github.com/jackc/pgx/v5/pgtype"
)

type Date struct{ civil.Date }

func (d *Date) ScanDate(v pgtype.Date) error {
	*d = Date{civil.DateOf(v.Time)}
	return nil
}

// https://github.com/jackc/pgx/blob/master/pgtype/interval.go#L13
const (
	microsecondsPerSecond  = 1_000_000
	microsecondsPerMinute  = 60 * microsecondsPerSecond
	microsecondsPerHour    = 60 * microsecondsPerMinute
	microsecondsPerDay     = 24 * microsecondsPerHour
	microsecondsPerMonth   = 30 * microsecondsPerDay

	maxRepresentableByTime = 24*60*60*1000000 - 1
)

type Time struct{ civil.Time }

// https://github.com/jackc/pgx/blob/master/pgtype/builtin_wrappers.go#L466
func (t *Time) ScanTime(v pgtype.Time) error {
	if !v.Valid {
		return fmt.Errorf("cannot scan NULL into *civil.Time")
	}

	if v.Microseconds > maxRepresentableByTime {
		return fmt.Errorf("%d microseconds cannot be represented as civil.Time", v.Microseconds)
	}

	usec := v.Microseconds
	hours := usec / microsecondsPerHour
	usec -= hours * microsecondsPerHour
	minutes := usec / microsecondsPerMinute
	usec -= minutes * microsecondsPerMinute
	seconds := usec / microsecondsPerSecond
	usec -= seconds * microsecondsPerSecond
	ns := usec * 1000

	*t = Time{civil.Time{
		Hour: int(hours),
		Minute: int(minutes),
		Second: int(seconds),
		Nanosecond: int(ns),
	}}

	return nil
}
