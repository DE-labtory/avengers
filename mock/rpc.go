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

type Client struct {
	ProcessId string
	CallFunc  func(queue string, params interface{}, callback interface{}) error // network manager grpc call
}

func NewClient(callFunc func(queue string, params interface{}, callback interface{}) error) Client {
	client := Client{
		CallFunc: callFunc,
	}
	return client
}

func (c *Client) Call(queue string, params interface{}, callback interface{}) error {
	return c.CallFunc(queue, params, callback)
}

type Server struct {
	ProcessId    string
	RegisterFunc func(processId string, queue string, handler func(a interface{}) error) error // network manager grpc consume
}

func NewServer(processId string, registerFunc func(processId string, queue string, handler func(a interface{}) error) error) Server {
	server := Server{
		ProcessId:    processId,
		RegisterFunc: registerFunc,
	}
	return server
}

func (s Server) Register(queue string, handler func(a interface{}) error) error {
	return s.RegisterFunc(s.ProcessId, queue, handler)
}
