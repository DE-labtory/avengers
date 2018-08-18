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

	"github.com/it-chain/engine/common/command"
)

type Process struct {
	Id                  string
	GrpcCommandHandlers []func(command command.ReceiveGrpc) error
	GrpcCommandReceiver chan interface{} //should be register to network's channel map
}

func NewProcess() Process {
	return Process{}
}

func (p *Process) Init(id string) {
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
					handler(message.(command.ReceiveGrpc))
				}

			case <-time.After(5 * time.Second):
				end = false
			}
		}
	}()
}

//func (p *Process) RegisterHandler(handler func(command command.ReceiveGrpc) error) error {
//
//	p.GrpcCommandHandlers = append(p.GrpcCommandHandlers, handler)
//
//	return nil
//}
