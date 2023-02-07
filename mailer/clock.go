package mailer

import "time"

// Clock returns the current time.
type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (rc realClock) Now() time.Time { return time.Now() }
