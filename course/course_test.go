package course

import (
	"testing"

	"cloud.google.com/go/civil"
)

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
	c, err = NewCourse("pilates", start, end, 10)
	if err != nil {
		t.Fatal(err)
	}

	if num := c.NumClasses(); num != 20 {
		t.Fatalf("Pilates does not have 20 classes but %d", num)
	}
}
