package database

import (
	"time"
)

// TimeRange time calculator from, to val int
func TimeRange(from, to time.Duration) (gte time.Time, lte time.Time) {
	fromVal := dayHours * from
	toVal := dayHours * to
	gte = time.Now().Add(-time.Hour * fromVal)
	lte = time.Now().Add(-time.Hour * toVal)
	return gte, lte
}
