package basics

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
)

func Print(app *newrelic.Application) {
	txn := app.StartTransaction("Print")
	defer txn.End()
	log.Println("Im inside print method doing nothing....")
}

func Hello(app *newrelic.Application, name string) (string, error) {

	// Monitor a transaction
	txn := app.StartTransaction("Hello")
	defer txn.End()

	// If no name was given, return an error with a message.
	if name == "" {
		return name, errors.New("empty name")
	}

	// Create a message using a random format.
	message := fmt.Sprintf(randomFormat(app), name)
	txn.AddAttribute("message", message)
	return message, nil
}

// randomFormat returns one of a set of greeting messages.
// The returned message is selected at random.
func randomFormat(app *newrelic.Application) string {

	// Workshop > Monitor a transaction
	// https://docs.newrelic.com/docs/apm/agents/go-agent/instrumentation/instrument-go-transactions/#go-txn
	txn := app.StartTransaction("randomFormat")
	defer txn.End()

	// Random sleep to simulate delays
	randomDelayOuter := rand.Intn(40)
	time.Sleep(time.Duration(randomDelayOuter) * time.Microsecond)

	// Workshop > Create a segment
	// https://docs.newrelic.com/docs/apm/agents/go-agent/instrumentation/instrument-go-segments
	tseg := txn.StartSegment("Formats")

	// Random sleep to simulate delays
	randomDelayInner := rand.Intn(80)
	time.Sleep(time.Duration(randomDelayInner) * time.Microsecond)

	// A slice of message formats.
	formats := []string{
		"Hi, %v. Welcome!",
		"Great to see you, %v!",
		"Good day, %v! Well met!",
		"%v! Hi there!",
		"Greetings %v!",
		"Hello there, %v!",
	}

	// Workshop > End a segment
	// https://docs.newrelic.com/docs/apm/agents/go-agent/instrumentation/instrument-go-segments
	tseg.End()

	// Return a randomly selected message format by specifying
	// a random index for the slice of formats.
	return formats[rand.Intn(len(formats))]
}
