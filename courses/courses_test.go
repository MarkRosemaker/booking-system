package courses

import (
	"math/rand"
	"testing"
	"time"

	"github.com/MarkRosemaker/booking-system/course"
	"golang.org/x/sync/errgroup"

	"cloud.google.com/go/civil"
)

func TestAdd(t *testing.T) {
	// create a bunch of dates
	k := 1000
	today := civil.DateOf(time.Now())
	dates := make([]civil.Date, k)
	for i := 0; i < k; i++ {
		dates[i] = today.AddDays(i - k/2)
	}

	// shuffle dates
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(k, func(i, j int) { dates[i], dates[j] = dates[j], dates[i] })

	createAndAdd := func(i int) error {
		if i >= k {
			t.Errorf("i is %d", i)
			return nil
		}

		start := dates[i]
		end := dates[i].AddDays(rand.Intn(100))
		// create a course
		c, err := course.NewHistoric("Test Course", start, end, 10)
		if err != nil {
			return err
		}

		Add(c)
		return nil
	}

	var eg errgroup.Group
	// TODO this doesn't work
	for i := 0; i < k; i++ {
		eg.Go(func() error { return createAndAdd(i) })
	}

	if err := eg.Wait(); err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	_, err := Get("fake ID")
	if err == nil {
		t.Errorf("didn't get error message for fake id")
	}

	// mux.Lock()
	// defer mux.Unlock()

	// if c, ok := byID[id]; ok {
	// 	return c, nil
	// }
	// return nil, fmt.Errorf("course with id '%s' does not exist", id)
}

// Add adds a course to the collection.
// If the course already exists, an error is returned. A course is considered a duplicate if the ID is the same or if there is another course with the same dates and the same names.
// func Add(c *course.Course) error {
// 	// later:
// 	// - consider passing context as input
// 	// - save to DB

// 	mux.Lock()
// 	defer mux.Unlock()

// 	// check if the ID exists already
// 	if _, ok := byID[c.ID()]; ok {
// 		return api.ErrBadRequest(fmt.Errorf(
// 			"a course with the ID '%s' has already been added", c.ID()))
// 	}

// 	// it's okay to add a course with the same name, but not if it's on the same dates
// 	if sameName, ok := byName[c.Name()]; ok {
// 		for _, o := range sameName {
// 			if c.Start() == o.Start() && c.End() == o.End() {
// 				return api.ErrBadRequest(fmt.Errorf(
// 					"a course '%s' with the same dates has already been added", c.Name()))
// 			}
// 		}
// 		// all other courses with that name have different dates, we can add the course to the list
// 		byName[c.Name()] = append(sameName, c)
// 	} else {
// 		// it's the first course with that name
// 		byName[c.Name()] = Courses{c}
// 	}

// 	// add to id map
// 	byID[c.ID()] = c

// 	// insert to sorted lists in the right place
// 	byStart = byStart.add(c, func(i int) bool {
// 		return c.Start().Before(byStart[i].Start())
// 	})
// 	byEnd = byEnd.add(c, func(i int) bool {
// 		return c.End().Before(byEnd[i].End())
// 	})

// 	return nil
// }
