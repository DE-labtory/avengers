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
	"fmt"

	"sync"

	"time"

	"github.com/it-chain/engine/common/command"
)

//network manager builds environment for communication between multiple nodes in network

type NetworkManager struct {
	mutex      sync.Mutex
	ChannelMap map[string]map[string]chan command.ReceiveGrpc // channel for receive deliverGrpc command
}

func NewNetworkManager() NetworkManager {

	return NetworkManager{
		ChannelMap: make(map[string]map[string]chan command.ReceiveGrpc),
	}
}

//receiver => processId
//queue name => queue
func (n *NetworkManager) Push(processId string, queue string, c command.ReceiveGrpc) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if n.ChannelMap[processId][queue] == nil {
		n.ChannelMap[processId] = make(map[string]chan command.ReceiveGrpc)
		n.ChannelMap[processId][queue] = make(chan command.ReceiveGrpc)
	}

	n.ChannelMap[processId][queue] <- c
	return nil
}

//Grpc call function would be injected to rpc client
//processId => sender
//params.RecipientList => receiver
func (n *NetworkManager) GrpcCall(processId string, queue string, params interface{}, callback interface{}) error {

	//find receiver process and deliver command through channel
	for _, v := range params.(command.DeliverGrpc).RecipientList {
		//convert grpc deliver message to grpc receive message
		extracted := command.ReceiveGrpc{
			Body:         params.(command.DeliverGrpc).Body,
			Protocol:     params.(command.DeliverGrpc).Protocol,
			ConnectionID: processId, // sender of this message
		}

		go func(v string, queue string) {
			fmt.Println("insert into channel:", v, extracted)
			queue = "message.receive"
			n.Push(v, queue, extracted)

		}(v, queue)
	}

	//run go routine for receive callback
	//go func() {
	//	task := <-n.CallbackQueue[corrId]
	//	handleResponse(task.(amqp.Delivery).Body, callback)
	//}()
	return nil
}

//GrpcConsume would be injected to rpc server
//processId => receiver
func (n *NetworkManager) GrpcConsume(processId string, queue string, handler func(command command.ReceiveGrpc) error) error {


	//defer

	//start command distributer
	go func(processId string, queue string) {

		if n.ChannelMap[processId][queue] == nil {
			n.mutex.Lock()
			n.ChannelMap[processId] = make(map[string]chan command.ReceiveGrpc)
			n.ChannelMap[processId][queue] = make(chan command.ReceiveGrpc)
			n.mutex.Unlock()
		}
		//fmt.Println("asdf",<-n.ChannelMap[processId][queue])
		end := true
		fmt.Println("listening!", "process: ", processId, "queue:", queue)
		for end {
			select {
			case message := <-n.ChannelMap[processId][queue]:
				fmt.Println("receive message from : ",processId," message:", message)
				handler(message)

			case <-time.After(4 * time.Second):
				fmt.Println("failed to consume, timed out!")
				end = false
			}
		}
	}(processId, queue)

	return nil
}

func (n *NetworkManager) AddProcess(process Process) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	if n.ChannelMap[process.Id]["message.receive"] == nil {
		n.ChannelMap[process.Id] = make(map[string]chan command.ReceiveGrpc)
		n.ChannelMap[process.Id]["message.receive"] = make(chan command.ReceiveGrpc)
	}
	n.ChannelMap[process.Id]["message.receive"] = process.GrpcCommandReceiver
}
