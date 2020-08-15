package classes

import (
	"net/http"

	"cloud.google.com/go/civil"
	"github.com/MarkRosemaker/booking-system/course"
	"github.com/MarkRosemaker/booking-system/courses"
	"github.com/MarkRosemaker/go-server/server/form"

	"github.com/MarkRosemaker/go-server/server/api"
	"github.com/MarkRosemaker/go-server/server/context"
)

func Respond(req *http.Request) interface{} {
	ctx, cancel := context.WithUserTimeout(req)
	defer cancel()

	var (
		name       string
		start, end civil.Date
		capacity   int
		historic   bool
		err        error
	)

	// get all the user input

	if name, err = form.GetStringE(req, "name"); err != nil {
		return api.ErrBadRequest(err)
	}

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

	errChan := make(chan error)
	go func() <-chan error {
		// for now: just add to array (quick)
		// later: add to database (potentially slow)

		var (
			c   *course.Course
			err error
		)

		if historic {
			c, err = course.NewHistoric(name, start, end, capacity)
		} else {
			c, err = course.New(name, start, end, capacity)
		}
		if err != nil {
			errChan <- api.ErrBadRequest(err)
			return errChan
		}

		errChan <- courses.Add(c)
		return errChan
	}()

	// timout if necessary
	select {
	case err = <-errChan:
		if err != nil {
			return api.ErrWrap(err)
		}
		// return api.NewSuccessNow(http.StatusCreated, "created course '%s' with given parameters", name)
		return api.NewSuccessNow(0, "created course '%s' with given parameters", name)
	case <-ctx.Done():
		return api.ErrWrap(ctx.Err())
	}
}
