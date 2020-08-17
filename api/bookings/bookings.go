package bookings

import (
	"net/http"
	"time"

	"github.com/MarkRosemaker/booking-system/courses"

	"cloud.google.com/go/civil"
	"github.com/MarkRosemaker/booking-system/course"
	"github.com/MarkRosemaker/go-server/server/api"
	"github.com/MarkRosemaker/go-server/server/context"
	"github.com/MarkRosemaker/go-server/server/form"
)

// Respond is the response function to an API request to '/bookings'.
//
// It parses the form input for a person's 'name', the 'date' of a class, and the 'id' of the course.
// Optionally, a 'timeout' parameter can be given.
//
// If any input does not make sense, an error is returned. Otherwise, the name is added to the attendees of the class on that date.
func Respond(req *http.Request) interface{} {
	ctx, cancel := context.WithUserTimeout(req)
	defer cancel()

	var (
		name string
		date civil.Date
		id   string
		c    *course.Course
		err  error
	)

	// get all the user input

	if name, err = form.GetStringE(req, "name"); err != nil {
		return api.ErrBadRequest(err)
	}

	if date, err = form.GetDateE(req, "date"); err != nil {
		return api.ErrBadRequest(err)
	}

	if id, err = form.GetStringE(req, "id"); err != nil {
		return api.ErrBadRequest(err)
	}

	errChan := make(chan error)
	go func() <-chan error {
		// for now: just get from map (quick)
		// later: get from database (potentially slow)

		if c, err = courses.Get(id); err != nil {
			errChan <- api.ErrBadRequest(err)
			return errChan
		}

		errChan <- c.BookClass(name, date)
		return errChan
	}()

	// timout if necessary
	select {
	case err = <-errChan:
		if err != nil {
			return api.ErrWrap(err)
		}
		return api.NewSuccessNow(0, nil, "Congratulations, %s! You are now registered for the %s class on %s.", name, c.Name(), date.In(time.Local).Format("Monday, 2. January 2006"))
	case <-ctx.Done():
		return api.ErrWrap(ctx.Err())
	}
}
