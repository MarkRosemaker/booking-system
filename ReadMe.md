# Booking System via API

- [Booking System via API](#booking-system-via-api)
	- [Usage](#usage)
	- [Coding Challenge Implementation](#coding-challenge-implementation)
		- [Terminology](#terminology)
		- [Summary](#summary)
			- [Classes](#classes)
			- [Bookings](#bookings)
		- [Historic Flag: Safeguard Against Invalid Dates](#historic-flag-safeguard-against-invalid-dates)
		- [Timeout Parameter](#timeout-parameter)
		- [Unique IDs](#unique-ids)
	- [Additions](#additions)
		- [go-server](#go-server)
		- [Pages for Your Convenience](#pages-for-your-convenience)

## Usage

Install the files as usual:

`go get github.com/MarkRosemaker/booking-system`

Compile and run in the repository folder.

The API routes are then available at http://localhost:8080/classes/ and http://localhost:8080/bookings/.

## Coding Challenge Implementation

### Terminology

To avoid confusion, I decided to use the word `course` for something that happens over several days from a certain start date until an end date.

A `class` is then one day in a course. There is one class for every day of a course.

### Summary

`main.go` starts a server (using my repository [`go-server`](https://github.com/markrosemaker/go-server)) that connects a handler with the API route.

#### Classes

In [`booking-system/api/classes/classes.go`](https://github.com/MarkRosemaker/booking-system/blob/master/api/classes/classes.go), `func Respond(req *http.Request) interface{}` calculates the response to the request to `/classes`:

- The form is parsed to get the 'name' of the course, the 'start' and 'end' dates (as [civil.Date](https://pkg.go.dev/cloud.google.com/go/civil?tab=doc)), the 'capacity', and the 'historic' flag (more about that later).
- From that, a [`course.Course`](https://github.com/MarkRosemaker/booking-system/blob/master/course/course.go) with a unique ID is created.
- That Course is then added to our [list of courses](https://github.com/MarkRosemaker/booking-system/blob/master/courses/courses.go).
- A success or error message is returned by the function.

This `interface{}` is encoded into JSON and returned to the user. An HTTP status code is stored in the object and the handler will write that code in the header. This implementation of the handler can be viewed at [`go-server/server/api/endpoint_base.go`](https://github.com/MarkRosemaker/go-server/blob/master/server/api/endpoint_base.go).

If an error occurred, `Respond` will return an [`api.Error` (from `go-server`)](https://github.com/MarkRosemaker/go-server/blob/master/server/api/error.go).

Otherwise, `Respond` will return an [`api.Success`](https://github.com/MarkRosemaker/go-server/blob/master/server/api/success.go).

#### Bookings

In [`booking-system/api/bookings/bookings.go`](https://github.com/MarkRosemaker/booking-system/blob/master/api/bookings/bookings.go), `func Respond(req *http.Request) interface{}` calculates the response to the request to `/bookings`:

- The form is parsed to get the 'name' of the customer, the 'date' on which they want to attend the class, and the 'id' of the course the class is part of.
- From that, the right [`course.Course`](https://github.com/MarkRosemaker/booking-system/blob/master/course/course.go) is fetched from the ID.
- Via the function [`BookClass`](https://github.com/MarkRosemaker/booking-system/blob/master/course/course.go), the customer's name is then added to the list of names of the class that they are attending.
- A success or error message is returned by the function.

### Historic Flag: Safeguard Against Invalid Dates

When creating courses, we most likely don't want to add courses that are already in the past.

Therefore, such an input is rejected unless we add a parameter 'historic' and set it to 'true' or similar.

### Timeout Parameter

Optionally, you can set a 'timeout' duration. For now, the program is very fast and a timeout is not needed.

However, in a real-world application we need to consider a delay; say, because the connection to the database is slow.

### Unique IDs

As an ID, we have an `uint64`. A new ID is created by simple incrementation of a counter.

When we store the courses in a database and restart the program, we need to remember to continue counting where we left off.

An alternative way to create unique IDs is the package "[github.com/google/uuid](https://github.com/google/uuid)", which creates unique, albeit long, ID strings.

## Additions

For various reasons, like my enjoyment of the project and a desire to learn, I've done a bit more than what was required.

### go-server

Since my implementation of the server is out of the scope of the project, I decided to use this part of the code from a separate repository called [`go-server`](https://github.com/markrosemaker/go-server).

### Pages for Your Convenience

The server hosts templates and files from the folder `site`.

- At http://localhost:8080/create-courses, you can test the course creation with a form.
- At http://localhost:8080/courses, you can see all courses and test the booking process.
- At http://localhost:8080/invalid, you can see what happens if the parameters are invalid since it the form won't restrict input values.
- At http://localhost:8080/too-slow, you can test out what happens if the request takes too long.
