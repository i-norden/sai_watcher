package test_helpers

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type MockWatchedEventsRepository struct {
	EventNames   []string
	ReturnEvents []*core.WatchedEvent
}

func (wer *MockWatchedEventsRepository) GetWatchedEvents(name string) ([]*core.WatchedEvent, error) {
	wer.EventNames = append(wer.EventNames, name)
	returnVal := wer.ReturnEvents
	// wipe return val after returning to only return once if called from loop
	wer.ReturnEvents = []*core.WatchedEvent{}
	return returnVal, nil
}
