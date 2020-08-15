package course

import (
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/civil"
	"github.com/MarkRosemaker/go-server/server/api"
	"github.com/google/uuid"
)

// For now, assume that there will be only one class per given day.
// "Ex: If a class by name pilates starts on 1st Dec and ends on 20th Dec, with capacity 10, that means Pilates has 20 classes and for each class the maximum capacity of attendance is 10."
type Course struct {
	id       string
	name     string
	start    civil.Date
	end      civil.Date
	capacity int
	classes  []*class
}

type class struct {
	// course    *Course
	attendees []string // later: slice of Member
}

// getter methods

func (c Course) ID() string {
	return c.id
}

func (c Course) Name() string {
	return c.name
}

func (c Course) Start() civil.Date {
	return c.start
}

func (c Course) End() civil.Date {
	return c.end
}

func (c Course) Capacity() int {
	return c.capacity
}

// initializers

func NewHistoric(name string, start, end civil.Date, capacity int) (*Course, error) {

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

func New(name string, start, end civil.Date, capacity int) (*Course, error) {
	// check if the course is in the past because then it might be a faulty input
	today := civil.DateOf(time.Now())
	if end.Before(today) {
		return nil, fmt.Errorf("invalid course parameters: course is in the past (ended %d days ago)", today.DaysSince(end))
	}

	return NewHistoric(name, start, end, capacity)
}

func (c Course) getClassOn(date civil.Date) (*class, error) {
	if date.Before(c.start) || date.After(c.end) {
		return nil, api.ErrBadRequest(fmt.Errorf("the chosen date is not within the timeframe of the course"))
	}

	idx := date.DaysSince(c.start)
	return c.classes[idx], nil // safe because of the above check
}

func (c Course) BookClass(name string, date civil.Date) error {
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
		if att == name {
			return api.ErrBadRequest(fmt.Errorf("you are already attending this class"))
		}
	}

	class.attendees = append(class.attendees, name)
	if over := len(class.attendees) - c.capacity; over > 0 {
		log.Printf("course %s (%s) over capacity by %d on %s", c.name, c.id, over, date)
	}

	return nil
}

func (c Course) NumClasses() int {
	return c.end.DaysSince(c.start) + 1
}
