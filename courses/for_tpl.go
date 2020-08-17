package courses

import (
	"sort"
	"time"

	"cloud.google.com/go/civil"
)

func All() Courses {
	return byStart
}

// Upcoming returns all upcoming courses, i.e. courses which start date is after today.
func Upcoming() Courses {
	today := civil.DateOf(time.Now())

	idx := sort.Search(len(byStart), func(i int) bool {
		return byStart[i].Start().After(today)
	})

	return byStart[idx:]
}

// Current returns all current courses, i.e. courses which start date is today or before,
// and which end date is today or after.
func Current() Courses {
	today := civil.DateOf(time.Now())

	idx := sort.Search(len(byEnd), func(i int) bool {
		return !byEnd[i].End().Before(today)
	})

	// since we ignore all past courses, the list is relatively small
	curr := make(Courses, 0)
	for _, c := range byEnd[idx:] {
		// filter upcoming courses
		if !c.Start().After(today) {
			curr = append(curr, c)
		}
	}

	return curr
}

// Past returns all past courses, i.e. courses which end date is before today.
func Past() Courses {
	today := civil.DateOf(time.Now())

	idx := sort.Search(len(byEnd), func(i int) bool {
		return !byEnd[i].End().Before(today)
	})

	return byStart[:idx]
}
