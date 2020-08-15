package bookings

import (
	"net/http"

	"github.com/MarkRosemaker/booking-system/courses"

	"cloud.google.com/go/civil"
	"github.com/MarkRosemaker/booking-system/course"
	"github.com/MarkRosemaker/go-server/server/api"
	"github.com/MarkRosemaker/go-server/server/context"
	"github.com/MarkRosemaker/go-server/server/form"
)

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
		return api.NewSuccessNow(0, "congratulations, you are now registered for the course '%s'", c.Name())
	case <-ctx.Done():
		return api.ErrWrap(ctx.Err())
	}
	return nil
}
