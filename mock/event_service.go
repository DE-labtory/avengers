/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mock

import (
	"errors"
	"reflect"
)

var ErrEventType = errors.New("Error type of event is not struct")

type EventService struct {
	ProcessId   string
	PublishFunc func(processId string, topic string, event interface{}) error
}

func NewEventService(processId string, publishFunc func(processId string, topic string, event interface{}) error) *EventService {
	return &EventService{
		ProcessId:   processId,
		PublishFunc: publishFunc,
	}
}

func (s *EventService) Publish(topic string, event interface{}) error {

	if !eventIsStruct(event) {
		return ErrEventType
	}

	if reflect.ValueOf(event).Elem().Type().Name() == "DeliverGrpc" {
		return s.PublishFunc(s.ProcessId, topic, event)
	}

	return nil
}

func eventIsStruct(event interface{}) bool {
	return reflect.TypeOf(event).Kind() == reflect.Struct
}
