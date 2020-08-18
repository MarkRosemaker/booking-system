package courses

import (
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/MarkRosemaker/booking-system/course"
	"golang.org/x/sync/errgroup"

	"cloud.google.com/go/civil"
)

const addedCourses = 1000

var today civil.Date = civil.DateOf(time.Now())

func addCourses(t *testing.T) {
	createAndAdd := func(start, end civil.Date) error {
		// create the course
		c, err := course.NewHistoric(
			uuid.New().String(), // unique name for each course
			start, end, 10)
		if err != nil {
			return err
		}

		return Add(c)
	}

	var eg errgroup.Group
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < addedCourses; i++ {

		// start before and after the current day
		start := today.AddDays(rand.Intn(100) - 50)
		end := start.AddDays(rand.Intn(100))

		eg.Go(func() error {
			return createAndAdd(start, end)
		})
	}

	if err := eg.Wait(); err != nil {
		t.Error(err)
	}
}

func TestAdd(t *testing.T) {
	c, err := course.NewHistoric("Test Course", today, today.AddDays(3), 10)
	if err != nil {
		t.Fatal(err)
	}

	if err = Add(c); err != nil {
		t.Errorf("couldn't add course: %s", err)
	}

	// don't add same course twice

	if err = Add(c); err == nil {
		t.Errorf("could add same course twice")
	}

	// don't add course with same name and dates (even if capacity is different)
	// idea: update capacity?

	var duplicate *course.Course
	duplicate, err = course.NewHistoric("Test Course", today, today.AddDays(3), 20)
	if err != nil {
		t.Fatal(err)
	}
	if err = Add(duplicate); err == nil {
		t.Errorf("could add course with same name and dates")
	}

	// but add course with same name and different dates

	var previous *course.Course
	previous, err = course.NewHistoric("Test Course", today.AddDays(-7), today.AddDays(-4), 10)
	if err != nil {
		t.Fatal(err)
	}
	if err = Add(previous); err != nil {
		t.Errorf("couldn't add previous course with same name: %s", err)
	}

	// create and add a bunch of courses
	addCourses(t)

	// all length values correct?

	count := 0
	for _, list := range byName {
		count += len(list)
	}

	if count != len(byID) {
		t.Errorf("byID doesn't have correct lenght, want: %d, have: %d", count, len(byID))
	}

	if count != len(byStart) {
		t.Errorf("byStart doesn't have correct lenght, want: %d, have: %d", count, len(byStart))
	}

	if count != len(byEnd) {
		t.Errorf("byEnd doesn't have correct lenght, want: %d, have: %d", count, len(byEnd))
	}

	// all sorted?

	sorted := sort.SliceIsSorted(byStart, func(i, j int) bool {
		return byStart[i].Start().Before(byStart[j].Start())
	})
	if !sorted {
		t.Errorf("byStart is not sorted")
	}

	sorted = sort.SliceIsSorted(byEnd, func(i, j int) bool {
		return byEnd[i].End().Before(byEnd[j].End())
	})
	if !sorted {
		t.Errorf("byEnd is not sorted")
	}
}

func TestGet(t *testing.T) {
	_, err := Get("fake ID")
	if err == nil {
		t.Errorf("didn't get error message for fake id")
	}

	rand.Seed(time.Now().UnixNano())
	start := today.AddDays(rand.Intn(100) - 50)
	end := start.AddDays(rand.Intn(100))

	c, err := course.NewHistoric("Test Course", start, end, 10)
	if err != nil {
		t.Fatalf("couldn't create test course: %s", err)
	}
	Add(c)

	var c2 *course.Course
	c2, err = Get(c.ID())
	if err != nil || c != c2 {
		t.Errorf("couldn't retrieve test course")
	}
}

func TestAll(t *testing.T) {
	addCourses(t)

	total := len(All())
	lenP := len(Past())
	lenC := len(Current())
	lenU := len(Upcoming())
	if total != lenP+lenC+lenU {
		t.Errorf("sum of courses not correct (%d past + %d current + %d upcoming != %d total)", lenP, lenC, lenU, total)
	}
}

func TestUpcoming(t *testing.T) {

	addCourses(t)

	up := Upcoming()

	if len(up) == 0 {
		t.Errorf("very unlikely occurance of 0 upcoming courses after randomly generating %d courses", addedCourses)
	}

	for _, c := range up {
		if !c.Start().After(today) {
			t.Errorf("course starting on %s was erronously given as upcoming course", c.Start())
		}
	}
}

func TestCurrent(t *testing.T) {

	addCourses(t)

	curr := Current()

	if len(curr) == 0 {
		t.Errorf("very unlikely occurance of 0 current courses after randomly generating %d courses", addedCourses)
	}

	for _, c := range curr {
		if c.Start().After(today) || c.End().Before(today) {
			t.Errorf("course from %s to %s was erronously given as current course", c.Start(), c.End())
		}
	}
}

func TestPast(t *testing.T) {

	addCourses(t)

	past := Past()

	if len(past) == 0 {
		t.Errorf("very unlikely occurance of 0 past courses after randomly generating %d courses", addedCourses)
	}

	for _, c := range past {
		if !c.End().Before(today) {
			t.Errorf("course ending on %s was erronously given as past course", c.End())
		}
	}
}
