// Package tpl contains structs and methods to pass data to the html templates.
package tpl

import (
	"net/http"
	"time"

	"cloud.google.com/go/civil"

	"github.com/MarkRosemaker/booking-system/courses"
)

// DataFunc returns data to be given to a template, given the request.
func DataFunc(req *http.Request) interface{} {
	return Data{Request: req}
}

// Data is a struct holding the data we want to pass to a template.
type Data struct {
	Request *http.Request
	Courses Courses
}

// Today returns the current date.
func (d Data) Today() civil.Date {
	// later: add this to the function map instead
	return civil.DateOf(time.Now())
}

// Courses is a dummy type to attach methods to.
//
// This implemenation was chosen so that we can use the intuitive notation {{ .Courses.All }}, {{ .Courses.Past }}, {{ .Courses.Current }}, and {{ .Courses.Upcoming }} in our templates.
type Courses int

// All returns all courses to be accessed by the template.
func (c Courses) All() courses.Courses {
	return courses.All()
}

// Past returns all past courses to be accessed by the template.
func (c Courses) Past() courses.Courses {
	return courses.Past()
}

// Current returns all current courses to be accessed by the template.
func (c Courses) Current() courses.Courses {
	return courses.Current()
}

// Upcoming returns all upcoming courses to be accessed by the template.
func (c Courses) Upcoming() courses.Courses {
	return courses.Upcoming()
}
