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

	"github.com/it-chain/engine/common/command"
	"time"
)

//network manager builds environment for communication between multiple nodes in network

type NetworkManager struct {
	mutex      sync.Mutex
	ChannelMap map[string]map[string]chan interface{} // channel for receive deliverGrpc command
}

func NewNetworkManager() NetworkManager {

	return NetworkManager{
		ChannelMap: make(map[string]map[string]chan interface{}),
	}
}

func (n *NetworkManager) Push(processId string, queue string, command interface{}) error {
	if n.ChannelMap[processId][queue] == nil{
		n.mutex.Lock()
		defer n.mutex.Unlock()
		n.ChannelMap[processId] = make(map[string]chan interface{})
		n.ChannelMap[processId][queue] = make(chan interface{})
	}
	n.ChannelMap[processId][queue]<-command
	return nil
}

//Grpc call function would be injected to rpc client
func (n *NetworkManager) GrpcCall(queue string, params interface{}, callback interface{}) {

	//find receiver process and deliver command through channel
	for _, v := range params.(command.DeliverGrpc).RecipientList {
		fmt.Println("calling process:", v)
		//convert grpc deliver message to grpc receive message
		extracted := command.ReceiveGrpc{
			Body:         params.(command.DeliverGrpc).Body,
			Protocol:     params.(command.DeliverGrpc).Protocol,
			ConnectionID: v,
		}

		queue = "message.receive"
		go func(v string, queue string) {
			fmt.Println("insert into channel:", v)

			n.Push(v, queue, extracted)

		}(v, queue)
	}

	//run go routine for receive callback
	//go func() {
	//	task := <-n.CallbackQueue[corrId]
	//	handleResponse(task.(amqp.Delivery).Body, callback)
	//}()

}

//GrpcConsume would be injected to rpc server
func (n *NetworkManager) GrpcConsume(processId string, queue string, handler func(a interface{}) error) error {
	if n.ChannelMap[processId][queue] == nil{
		n.mutex.Lock()
		defer n.mutex.Unlock()
		n.ChannelMap[processId] = make(map[string]chan interface{})
		n.ChannelMap[processId][queue] = make(chan interface{})
	}
	fmt.Println("consume from processId:", processId, "queue:", queue, "map:", n.ChannelMap[processId][queue])
	//start command distributer
	go func() {
		//fmt.Println(n.ChannelMap[processId][queue])
		//fmt.Println("asdf",<-n.ChannelMap[processId][queue])
		end := true

		for end {
			select {
			case message := <-n.ChannelMap[processId][queue]:
				fmt.Println("receive message: ", message)
				handler(message)

			case <-time.After(5 * time.Second):
				end = false
			}
		}
	}()

	return nil
}

func (m *NetworkManager) AddProcess(process Process) {

	m.ChannelMap[process.Id]["message.receive"] = process.grpcCommandReceiver
	process.Init(process.Id)
}
