package courses

import (
	"fmt"
	"sort"
	"sync"

	"github.com/MarkRosemaker/go-server/server/api"

	"github.com/MarkRosemaker/booking-system/course"
)

// Courses is a slice of courses.
type Courses []*course.Course

var (
	// maps for quick access
	byID   map[string]*course.Course = make(map[string]*course.Course)
	byName map[string]Courses        = make(map[string]Courses)

	// sorted lists
	byStart Courses = make(Courses, 0)
	byEnd   Courses = make(Courses, 0)

	// protect maps and lists with mutex
	mux *sync.Mutex = &sync.Mutex{}
)

// Get returns the course with the given id or an error, if no course with the ID exists.
func Get(id string) (*course.Course, error) {
	mux.Lock()
	defer mux.Unlock()

	if c, ok := byID[id]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("course with id '%s' does not exist", id)
}

// Add adds a course to the collection.
// If the course already exists, an error is returned. A course is considered a duplicate if the ID is the same or if there is another course with the same dates and the same names.
func Add(c *course.Course) error {
	// later:
	// - consider passing context as input
	// - save to DB

	mux.Lock()
	defer mux.Unlock()

	// check if the ID exists already
	if _, ok := byID[c.ID()]; ok {
		return api.ErrBadRequest(fmt.Errorf(
			"a course with the ID '%s' has already been added", c.ID()))
	}

	// it's okay to add a course with the same name, but not if it's on the same dates
	if sameName, ok := byName[c.Name()]; ok {
		for _, o := range sameName {
			if c.Start() == o.Start() && c.End() == o.End() {
				return api.ErrBadRequest(fmt.Errorf(
					"a course '%s' with the same dates has already been added", c.Name()))
			}
		}
		// all other courses with that name have different dates, we can add the course to the list
		byName[c.Name()] = append(sameName, c)
	} else {
		// it's the first course with that name
		byName[c.Name()] = Courses{c}
	}

	// add to id map
	byID[c.ID()] = c

	// insert to sorted lists in the right place
	byStart = byStart.add(c, func(i int) bool {
		return c.Start().Before(byStart[i].Start())
	})
	byEnd = byEnd.add(c, func(i int) bool {
		return c.End().Before(byEnd[i].End())
	})

	return nil
}

// add adds the course to the list, given a search function that determines where.
func (cs Courses) add(c *course.Course, f func(int) bool) Courses {
	k := len(cs)
	if k == 0 {
		return Courses{c}
	}

	switch idx := sort.Search(k, f); idx {
	case k, k - 1: // new value is the biggest
		return append(cs, c)
	default: // 0 to k-1
		// most likely the most efficient way to insert
		// see: https://github.com/golang/go/wiki/SliceTricks
		cs = append(cs, nil)
		copy(cs[idx+1:], cs[idx:])
		cs[idx] = c
		return cs
	}
}
