package main

import (
	"github.com/MarkRosemaker/booking-system/api/bookings"
	"github.com/MarkRosemaker/booking-system/api/classes"
	"github.com/MarkRosemaker/booking-system/tpl"
	"github.com/MarkRosemaker/go-server/server/api"

	"github.com/MarkRosemaker/go-server/server"
)

func main() {
	o := server.Options{
		ContentSource:    "site",
		TemplateDataFunc: tpl.DataFunc,
		Endpoints: api.Endpoints{
			api.BaseEndpoint{
				URL:          "/classes",
				ResponseFunc: classes.Respond},
			api.BaseEndpoint{
				URL:          "/bookings",
				ResponseFunc: bookings.Respond},
		},
		Verbose: true,
	}

	server.Run(o)
}
