package tpl

import (
	"fmt"
	"net/http"

	"github.com/MarkRosemaker/booking-system/courses"
)

type Data struct {
	Request *http.Request
	Courses Courses
}

func DataFunc(req *http.Request) interface{} {
	return Data{Request: req}
}

type Courses int

func (c Courses) String() string {
	return fmt.Sprintf("%s", courses.All())
}

func (c Courses) All() courses.Courses {
	return courses.All()
}

func (c Courses) Past() courses.Courses {
	return courses.Past()
}

func (c Courses) Current() courses.Courses {
	return courses.Current()
}

func (c Courses) Upcoming() courses.Courses {
	return courses.Upcoming()
}
