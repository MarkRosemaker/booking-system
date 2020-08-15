package bookings

import (
	"net/http"
)

func Respond(req *http.Request) interface{} {
	// ctx, cancel := context.WithUserTimeout(req)
	// defer cancel()

	// origin, err := point.FromRequest(req)
	// if err != nil {
	// 	return api.ErrBadRequest(err)
	// }

	// var r int
	// if r, err = radius.FromRequest(req); err != nil {
	// 	return api.ErrBadRequest(err)
	// }

	// success := make(chan bool)
	// go func() <-chan bool {
	// 	// nb = pts.Neighbors(*origin, r)
	// 	success <- true
	// 	return success
	// }()

	// timout if necessary
	// select {
	// case <-success:
	// 	return nb
	// case <-ctx.Done():
	// 	return api.ErrWrap(ctx.Err())
	// }
	return nil
}
