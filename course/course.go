package course

import (
	"fmt"
	"time"

	"cloud.google.com/go/civil"
)

// For now, assume that there will be only one class per given day.
// "Ex: If a class by name pilates starts on 1st Dec and ends on 20th Dec, with capacity 10, that means Pilates has 20 classes and for each class the maximum capacity of attendance is 10."
type Course struct {
	name       string
	start, end civil.Date
	capacity   int
}

func NewHistoric(name string, start, end civil.Date, capacity int) (*Course, error) {

	if start.After(end) {
		return nil, fmt.Errorf("invalid course parameters: start date (%s) after end date (%s)", start, end)
	}

	if capacity < 1 {
		return nil, fmt.Errorf("invalid course parameters: capacity (%d) must be positive", capacity)
	}

	return &Course{name, start, end, capacity}, nil
}

func New(name string, start, end civil.Date, capacity int) (*Course, error) {
	c, err := NewHistoric(name, start, end, capacity)
	if err != nil {
		return nil, err
	}

	// also check if the course is in the past because then it might be a faulty input
	today := civil.DateOf(time.Now())
	if end.Before(today) {
		return nil, fmt.Errorf("invalid course parameters: course is in the past", start, end)
	}

	return c, nil
}

func (c Course) NumClasses() int {
	return c.end.DaysSince(c.start) + 1
}
