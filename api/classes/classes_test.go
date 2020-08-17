package classes

import (
	"fmt"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/MarkRosemaker/go-server/server/api"
)

func TestRespond(t *testing.T) {
	today := civil.DateOf(time.Now())

	tables := []struct {
		params string
		res    string
	}{
		// test all errors
		{"",
			"400 Bad Request: name value not provided"},
		{"?name=Karate",
			"400 Bad Request: start value not provided"},
		{"?name=Karate&start=now",
			"400 Bad Request: start value 'now' could not be parsed to date"},
		{"?name=Karate&start=2010-30-02",
			"400 Bad Request: start value '2010-30-02' could not be parsed to date"},
		{"?name=Karate&start=2010-01-02",
			"400 Bad Request: end value not provided"},
		{"?name=Karate&start=2010-01-01&end=2010-02-01",
			"400 Bad Request: capacity value not provided"},
		{"?name=Mesoamerican_Ballgame&start=2010-01-01&end=2010-02-01&capacity=10",
			"400 Bad Request: invalid course parameters: course is in the past"},
		{"?name=Time_Travelling&start=2015-10-21&end=1985-10-26&capacity=10&historic=true",
			"400 Bad Request: invalid course parameters: start date (2015-10-21) after end date (1985-10-26)"},
		{fmt.Sprintf("?name=Karate&start=%s&end=%s&capacity=10", today.AddDays(100), today.AddDays(10)),
			fmt.Sprintf("400 Bad Request: invalid course parameters: start date (%s) after end date (%s)", today.AddDays(100), today.AddDays(10))},

		// successful course creation
		{fmt.Sprintf("?name=Karate&start=%s&end=%s&capacity=10", today, today),
			"course created"},
		{fmt.Sprintf("?name=Karate&start=%s&end=%s&capacity=10", today.AddDays(100), today.AddDays(230)),
			"course created"},
	}

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

	// "Ex: If a class by name pilates starts on 1st Dec and ends on 20th Dec, with capacity 10, that means Pilates has 20 classes and for each class the maximum capacity of attendance is 10."
	resp := Respond(httptest.NewRequest("GET", "/classes?name=pilates&start=2019-12-01&end=2019-12-20&capacity=10&historic=true", nil))
	if s, ok := resp.(api.Success); ok {
		v := reflect.ValueOf(s.Object)
		if classes := v.FieldByName("Classes").Int(); classes != 20 {
			t.Errorf("Pilates from the example in the challenge specification doesn't have 20 classes, it has: %d", classes)
		}
		if capacity := v.FieldByName("Capacity").Int(); capacity != 10 {
			t.Errorf("Pilates from the example in the challenge specification doesn't have capacity 10, it has: %d", capacity)
		}
	} else {
		t.Errorf("Didn't receive api.Success when testing the example in the challenge specification, got: %T", resp)
	}
}
