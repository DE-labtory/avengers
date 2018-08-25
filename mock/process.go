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
	"time"

	"reflect"
	"github.com/it-chain/engine/common/command"
	"sync"
)

type Process struct {
	mutex sync.Mutex
	Id                  string
	GrpcCommandHandlers []func(command command.ReceiveGrpc) error
	GrpcCommandReceiver chan command.ReceiveGrpc //should be register to network's channel map
	Services            map[string]interface{}   // register service or api for testing which has injected mock client
}

func NewProcess() Process {
	return Process{}
}

func (p *Process) Init(id string) {
	p.Services = make(map[string]interface{})
	p.Id = id
	p.GrpcListen()
}

//every grpc command handler listens command and deal with it
func (p *Process) GrpcListen() {

	go func() {

		end := true

		for end {
			select {
			case message := <-p.GrpcCommandReceiver:
				for _, handler := range p.GrpcCommandHandlers {
					handler(message)
				}

			case <-time.After(5 * time.Second):
				end = false
			}
		}
	}()
}

func (p *Process) Register(service interface{}) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.Services[reflect.ValueOf(service).Elem().Type().Name()] = service
}

//func (p *Process) RegisterHandler(handler func(command command.ReceiveGrpc) error) error {
//
//	p.GrpcCommandHandlers = append(p.GrpcCommandHandlers, handler)
//
//	return nil
//}
