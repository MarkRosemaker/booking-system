package classes

import (
	"net/http"

	"cloud.google.com/go/civil"
	"github.com/MarkRosemaker/booking-system/course"
	"github.com/MarkRosemaker/booking-system/courses"
	"github.com/MarkRosemaker/go-server/server/form"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/MarkRosemaker/go-server/server/api"
	"github.com/MarkRosemaker/go-server/server/context"
)

var toTitleCase cases.Caser = cases.Title(language.English)

// Respond is the response function to an API request to '/classes'.
//
// It parses the form input for the course 'name', the 'start' and 'end' dates of the course, and the 'capacity' of the course.
// Optionally, a 'timeout' and 'historic' parameter can be given. The latter signifies whether or not we want to allow the course to be in the past.
//
// If any input does not make sense, an error is returned. Otherwise, the course is added to the list of courses.
func Respond(req *http.Request) interface{} {
	ctx, cancel := context.WithUserTimeout(req)
	defer cancel()

	var (
		name       string
		start, end civil.Date
		capacity   int
		historic   bool

		c *course.Course

		err error
	)

	// get all the user input

	if name, err = form.GetStringE(req, "name"); err != nil {
		return api.ErrBadRequest(err)
	}
	name = toTitleCase.String(name)

	if start, err = form.GetDateE(req, "start"); err != nil {
		return api.ErrBadRequest(err)
	}

	if end, err = form.GetDateE(req, "end"); err != nil {
		return api.ErrBadRequest(err)
	}

	if capacity, err = form.GetIntE(req, "capacity"); err != nil {
		return api.ErrBadRequest(err)
	}

	if historic, err = form.GetBoolE(req, "historic"); err != nil {
		return api.ErrBadRequest(err)
	}

	if historic {
		c, err = course.NewHistoric(name, start, end, capacity)
	} else {
		c, err = course.New(name, start, end, capacity)
	}
	if err != nil {
		return api.ErrBadRequest(err)
	}

	errChan := make(chan error)
	go func() <-chan error {
		// for now: add to arrays and map, check for duplicates (quick)
		// later: add to database (potentially slow)

		errChan <- courses.Add(c)
		return errChan
	}()

	// timout if necessary
	select {
	case err = <-errChan:
		if err != nil {
			return api.ErrWrap(err)
		}
		// return new information about the course, such as ID and number of classes
		return api.NewSuccessNow(http.StatusCreated, struct {
			ID       string
			Name     string
			Start    civil.Date
			End      civil.Date
			Capacity int
			Classes  int
		}{
			c.ID(),
			c.Name(),
			c.Start(),
			c.End(),
			c.Capacity(),
			c.NumClasses(),
		}, "course created")
	case <-ctx.Done():
		return api.ErrWrap(ctx.Err())
	}
}
