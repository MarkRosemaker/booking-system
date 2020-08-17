package course

import (
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/civil"
	"github.com/MarkRosemaker/go-server/server/api"
	"github.com/google/uuid"
)

// A Course represents an online course and consists of one or more classes, one for each day of the duration of the course.
// All fields are private so they cannot be changed from another package.
type Course struct {
	id       string
	name     string
	start    civil.Date
	end      civil.Date
	capacity int
	classes  []*class // len(classes) == end-start +1 == NumClasses()
}

// A class represents one day of a course.
type class struct {
	// course *Course // as the project grows more complex, we might want to consider a pointer back to the original course

	// list of the names of the attendees
	// later, this could be a slice of a struct 'Member' with not just the name, but also member ID to avoid mix-ups
	attendees []string
}

// getter methods

// ID returns the course ID.
func (c Course) ID() string {
	return c.id
}

// Name returns the course name.
func (c Course) Name() string {
	return c.name
}

// Start returns the start date of the course.
func (c Course) Start() civil.Date {
	return c.start
}

// End returns the end date of the course.
func (c Course) End() civil.Date {
	return c.end
}

// Capacity returns the capacity of the course.
func (c Course) Capacity() int {
	return c.capacity
}

// initializers

// NewHistoric creates a new course, if the input passes some checks or an error, if not.
// The course may be in the past (i.e. be 'historic').
func NewHistoric(name string, start, end civil.Date, capacity int) (*Course, error) {

	if name == "" {
		return nil, fmt.Errorf("please provide a course name")
	}

	if start.After(end) {
		return nil, fmt.Errorf("invalid course parameters: start date (%s) after end date (%s)", start, end)
	}

	if capacity < 1 {
		return nil, fmt.Errorf("invalid course parameters: capacity (%d) must be positive", capacity)
	}

	c := &Course{
		id:       uuid.New().String(),
		name:     name,
		start:    start,
		end:      end,
		capacity: capacity}

	k := c.NumClasses()
	c.classes = make([]*class, k)
	for i := 0; i < k; i++ {
		c.classes[i] = &class{
			// course:    c,
			attendees: make([]string, 0)}
	}

	return c, nil
}

// New creates a new course, if the input passes some checks or an error, if not.
//
// The course may not be in the past, i.e. the end date is today or in the future.
// Otherwise, we assume that we have faulty input data.
func New(name string, start, end civil.Date, capacity int) (*Course, error) {
	// check if the course is in the past because then it might be a faulty input
	today := civil.DateOf(time.Now())
	if end.Before(today) {
		return nil, fmt.Errorf("invalid course parameters: course is in the past")
	}

	return NewHistoric(name, start, end, capacity)
}

// getClassOn returns the class of the course that is happening on a certain day.
func (c Course) getClassOn(date civil.Date) (*class, error) {
	if date.Before(c.start) || date.After(c.end) {
		return nil, api.ErrBadRequest(fmt.Errorf("the chosen date is not within the timeframe of the course"))
	}

	idx := date.DaysSince(c.start)
	return c.classes[idx], nil // safe because of the above check
}

// BookClass registers a customer for a class on the given day.
// That day must be during the course duration and be in the future.
// A customer can only book a class once.
func (c Course) BookClass(customer string, date civil.Date) error {
	today := civil.DateOf(time.Now())
	if today.After(c.end) {
		return api.ErrBadRequest(fmt.Errorf("the course is in the past"))
	}

	if date.Before(today) {
		return api.ErrBadRequest(fmt.Errorf("please pick a future date"))
	}

	class, err := c.getClassOn(date)
	if err != nil {
		return err
	}

	// obviously some people have the same names
	// in the future, attendees can be a slice of a 'Member' struct that contains member id etc.
	for _, att := range class.attendees {
		if att == customer {
			return api.ErrBadRequest(fmt.Errorf("you are already attending this class"))
		}
	}

	class.attendees = append(class.attendees, customer)
	if over := len(class.attendees) - c.capacity; over > 0 {
		// per specification, it is possible to overbook
		// so we simply log the overbooking
		log.Printf("course %s (%s) over capacity by %d on %s", c.name, c.id, over, date)
	}

	return nil
}

// NumClasses returns the number of classes for the course.
// For now, there is a class on every day of the duration of the course.
func (c Course) NumClasses() int {
	return c.end.DaysSince(c.start) + 1
}
