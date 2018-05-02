// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
