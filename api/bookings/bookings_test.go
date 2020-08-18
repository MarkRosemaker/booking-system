package bookings

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/MarkRosemaker/booking-system/course"
	"github.com/MarkRosemaker/booking-system/courses"
	"github.com/MarkRosemaker/go-server/server/api"
)

func TestRespond(t *testing.T) {

	// populate course list with test courses

	var (
		cPast, cTest *course.Course
		err          error
	)

	today := civil.DateOf(time.Now())
	cPast, err = course.NewHistoric("Yoga", today.AddDays(-20), today.AddDays(-10), 10)
	if err != nil {
		t.Fatalf("couldn't create past course")
	}
	cTest, err = course.New("Pilates", today, today.AddDays(3), 10)
	if err != nil {
		t.Fatalf("couldn't create test course")
	}

	if err = courses.Add(cPast); err != nil {
		t.Fatalf("couldn't add past course: %s", err)
	}
	if err = courses.Add(cTest); err != nil {
		t.Fatalf("couldn't add test course: %s", err)
	}

	// the output on successful bookings
	successOn := func(d civil.Date) string {
		return fmt.Sprintf("Congratulations, Arnold! You are now registered for the Pilates class on %s.", d.In(time.Local).Format("Monday, 2. January 2006"))
	}

	// create table

	tables := []struct {
		params string
		res    string
	}{
		// test all errors
		{"",
			"400 Bad Request: name value not provided"},
		{"?name=Arnold",
			"400 Bad Request: date value not provided"},
		{"?name=Arnold&date=now",
			"400 Bad Request: date value 'now' could not be parsed to date"},
		{"?name=Arnold&date=2010-02-30",
			"400 Bad Request: date value '2010-02-30' could not be parsed to date"},
		{"?name=Arnold&date=2010-02-01",
			"400 Bad Request: id value not provided"},
		{"?name=Arnold&date=2010-02-01&id=fake_ID",
			"400 Bad Request: id value 'fake_ID' could not be parsed to uint64"},
		{"?name=Arnold&date=2010-02-01&id=-1",
			"400 Bad Request: id value '-1' could not be parsed to uint64"},
		{"?name=Arnold&date=2010-02-01&id=0",
			"400 Bad Request: course with id 0 does not exist"},
		{fmt.Sprintf("?name=Arnold&date=2010-02-01&id=%d", cPast.ID()),
			"400 Bad Request: the course is in the past"},
		{fmt.Sprintf("?name=Arnold&date=%s&id=%d", today.AddDays(-15), cPast.ID()),
			"400 Bad Request: the course is in the past"},
		{fmt.Sprintf("?name=Arnold&date=%s&id=%d", cTest.Start(), cTest.ID()),
			successOn(cTest.Start())},
		{fmt.Sprintf("?name=Arnold&date=%s&id=%d", cTest.Start(), cTest.ID()),
			"400 Bad Request: you are already attending this class"},
		{fmt.Sprintf("?name=Arnold&date=%s&id=%d", today.AddDays(2), cTest.ID()),
			successOn(today.AddDays(2))},
	}

	// test

	for _, table := range tables {
		url := fmt.Sprintf("/classes%s", table.params)
		resp := Respond(httptest.NewRequest("GET", url, nil))

		switch v := resp.(type) {
		case api.Error:
			if s := v.Error(); s != table.res {
				t.Errorf("Result of %s was incorrect, got: %q, want: %q.",
					url, s, table.res)
			}
		case api.Success:
			if v.Message != table.res {
				t.Errorf("Result of %s was incorrect, got: '%s', want: '%s'.",
					url, v.Message, table.res)
			}
		default:
			t.Errorf("Result of %s has wrong type, expected: api.Error or api.Success, got: %T", url, resp)
		}
	}
}
