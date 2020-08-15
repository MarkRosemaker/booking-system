package courses

import (
	"fmt"
	"sort"
	"sync"

	"github.com/MarkRosemaker/go-server/server/api"

	"github.com/MarkRosemaker/booking-system/course"
)

type Courses []*course.Course

var (
	byID    map[string]*course.Course = make(map[string]*course.Course)
	byName  map[string]Courses        = make(map[string]Courses)
	byStart Courses                   = make(Courses, 0)
	byEnd   Courses                   = make(Courses, 0)
	mux     *sync.Mutex               = &sync.Mutex{}
)

func Get(id string) (*course.Course, error) {
	if c, ok := byID[id]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("course with id '%s' does not exist", id)
}

func Add(c *course.Course) error {
	// later:
	// - consider passing context as input
	// - save to DB

	mux.Lock()
	defer mux.Unlock()

	fmt.Println("TODO continue debugging here")

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
		byName[c.Name()] = append(sameName, c)
	} else {
		byName[c.Name()] = Courses{c}
	}

	byID[c.ID()] = c

	// insert to sorted lists in the right place
	byStart.add(c, func(i int) bool {
		return c.Start().Before(byStart[i].Start())
	})
	byEnd.add(c, func(i int) bool {
		return c.End().Before(byEnd[i].End())
	})

	return nil
}

func (cs *Courses) add(c *course.Course, f func(int) bool) {
	k := len(*cs)
	if k == 0 {
		cs = &Courses{c}
		return
	}

	switch idx := sort.Search(k, f); idx {
	case k, k - 1: // new value is the biggest
		list := append(*cs, c)
		cs = &list
	default:
		// most likely the most efficient way to insert
		// see: https://github.com/golang/go/wiki/SliceTricks
		list := append(*cs, nil)
		copy(list[idx+1:], list[idx:])
		list[idx] = c
		cs = &list
	}
}
