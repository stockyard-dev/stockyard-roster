package server

import "github.com/stockyard-dev/stockyard-roster/internal/license"

type Limits struct {
	MaxContacts int
	DealPipeline bool
	Reminders    bool
}

var freeLimits = Limits{
	MaxContacts:  25,
	DealPipeline: true,
	Reminders:    true,
}

var proLimits = Limits{
	MaxContacts:  0,
	DealPipeline: true,
	Reminders:    true,
}

func LimitsFor(info *license.Info) Limits {
	if info != nil && info.IsPro() {
		return proLimits
	}
	return freeLimits
}

func LimitReached(limit, current int) bool {
	if limit == 0 {
		return false
	}
	return current >= limit
}
