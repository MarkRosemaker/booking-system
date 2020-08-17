package course

import (
	"testing"
	"time"

	"cloud.google.com/go/civil"
)

func getPastDate(t *testing.T) civil.Date {
	pastDate, err := civil.ParseDate("2010-03-01")
	if err != nil {
		t.Fatalf("couldn't parse date")
	}
	return pastDate
}

func getPastCourse(t *testing.T) *Course {
	pastDate := getPastDate(t)
	c, err := NewHistoric("Karate", pastDate, pastDate.AddDays(10), 20)
	if err != nil {
		t.Fatalf("couldn't create course: %s", err)
	}

	return c
}

var today civil.Date = civil.DateOf(time.Now())

func getTestCourse(t *testing.T) *Course {
	c, err := New("Karate", today.AddDays(-10), today.AddDays(10), 20)
	if err != nil {
		t.Fatalf("couldn't create course: %s", err)
	}

	return c
}

func TestNewHistoric(t *testing.T) {

	pastDate := getPastDate(t)

	_, err := NewHistoric("", pastDate, pastDate, 10)
	if err == nil {
		t.Errorf("created course even though it doesn't have a name")
	}

	_, err = NewHistoric("Karate", pastDate, pastDate.AddDays(-1), 10)
	if err == nil {
		t.Errorf("start date should be before end date")
	}

	_, err = NewHistoric("Karate", pastDate, pastDate, 0)
	if err == nil {
		t.Errorf("don't allow zero capacity")
	}

	_, err = NewHistoric("Karate", pastDate, pastDate, -1)
	if err == nil {
		t.Errorf("don't allow negative capacity")
	}

	getPastCourse(t)

}

func TestNew(t *testing.T) {

	pastDate := getPastDate(t)

	_, err := New("Yoga", pastDate, pastDate, 10)
	if err == nil {
		t.Errorf("don't allow course in the past to be created")
	}

	getTestCourse(t)
}

// NOTE: For more tests of course creation, see api/classes/classes_test.go

func TestGetClassOn(t *testing.T) {
	pastDate := getPastDate(t)
	c := getPastCourse(t)

	cl, err := c.getClassOn(pastDate.AddDays(-1))
	if err == nil {
		t.Errorf("no error even though date too early")
	}

	cl, err = c.getClassOn(pastDate.AddDays(11))
	if err == nil {
		t.Errorf("no error even though date too late")
	}

	for i := 0; i <= 10; i++ {
		cl, err = c.getClassOn(pastDate.AddDays(i))
		if err != nil || cl == nil {
			t.Errorf("failed to initialize classes")
		}
	}
}

// // BookClass registers a customer for a class on the given day.
// // That day must be during the course duration and be in the future.
// // A customer can only book a class once.
func TestBookClass(t *testing.T) {

	pastDate := getPastDate(t)
	if getPastCourse(t).BookClass("Arnold", pastDate) == nil {
		t.Errorf("could book course that was in the past")
	}

	c := getTestCourse(t)

	if c.BookClass("Arnold", today.AddDays(-1)) == nil {
		t.Errorf("could book class for yesterday")
	}

	if err := c.BookClass("Arnold", today.AddDays(1)); err != nil {
		t.Errorf("could book class for tomorrow: %s", err)
	}

	if err := c.BookClass("Arnold", today.AddDays(1)); err == nil {
		t.Errorf("could book class for tomorrow twice")
	}
}

func TestNumClasses(t *testing.T) {
	var (
		start, end civil.Date
		c          *Course
		err        error
	)

	// "Ex: If a class by name pilates starts on 1st Dec and ends on 20th Dec, with capacity 10, that means Pilates has 20 classes and for each class the maximum capacity of attendance is 10."
	start, err = civil.ParseDate("2020-12-01")
	if err != nil {
		t.Fatal(err)
	}
	end, err = civil.ParseDate("2020-12-20")
	if err != nil {
		t.Fatal(err)
	}
	c, err = New("pilates", start, end, 10)
	if err != nil {
		t.Fatal(err)
	}

	if num := c.NumClasses(); num != 20 {
		t.Fatalf("Pilates does not have 20 classes but %d", num)
	}
}
