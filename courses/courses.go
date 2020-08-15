package courses

import (
	"github.com/MarkRosemaker/booking-system/course"
)

type Courses []*course.Course

var courseList Courses

func Add(c *course.Course) error {
	// for now it's simple, but later we might need to return an error
	courseList = append(courseList, c)
	return nil
}
