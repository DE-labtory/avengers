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

type NetworkManager struct {
	ChannelMap    map[string]chan interface{}
	CallbackQueue map[string]chan interface{}
}

func (n *NetworkManager) GrpcCall(queue string, params interface{}, callback interface{}) error {

	//find receiver process and deliver command through channel
	for _, v := range params.(command.DeliverGrpc).RecipientList {

		//convert grpc deliver message to grpc receive message
		extracted := command.ReceiveGrpc{
			Body:         params.(command.DeliverGrpc).Body,
			Protocol:     params.(command.DeliverGrpc).Protocol,
			ConnectionID: v,
		}

		n.ChannelMap[v] <- extracted
	}

	//run go routine for receive callback
	//go func() {
	//	task := <-n.CallbackQueue[corrId]
	//	handleResponse(task.(amqp.Delivery).Body, callback)
	//}()

	return nil
}

func (n *NetworkManager) GrpcConsume(queue string, handler func(a interface{}) error) error {

	//start command distributer

	go func() {

		end := true

		for end {
			select {
			case message := <-n.ChannelMap[queue]:
				handler(message)

			case <-time.After(2 * time.Second):
				end = false
			}
		}
	}()

	return nil
}

func (m *NetworkManager) AddProcess(process Process) {

	m.ChannelMap[process.Id] = process.grpcCommandReceiver
	process.Init(process.Id)
}