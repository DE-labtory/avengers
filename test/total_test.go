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

package test

import (
	"testing"

	"github.com/it-chain/avengers/mock"
	"github.com/it-chain/engine/common/command"
	"time"
)

func TestTotal(t *testing.T) {
	tests := map[string]struct {
		Input struct {
			ProcessList []string
		}
	}{
		"three node test": {
			Input: struct{ ProcessList []string }{ProcessList: []string{"1", "2", "3"}},
		},
	}

	//make virtual network

	for testName, test := range tests {
		t.Logf("running test case %s", testName)
		networkManager := mock.NewNetworkManager()
		clientList := make([]mock.Client,0)
		serverList := make([]mock.Server,0)

		for _, processId := range test.Input.ProcessList {
			process := mock.NewProcess()
			//process.Init(processId)

			client := mock.NewClient(processId, networkManager.GrpcCall)
			server := mock.NewServer(processId, networkManager.GrpcConsume)

			server.Register("message.receive", func(a interface{}) error {
				t.Logf("consumed process: %s", process.Id)
				return nil
			})

			clientList = append(clientList, client)
			serverList = append(serverList, server)


			networkManager.AddProcess(process)
		}

		clientList[0].Call("message.deliver",command.DeliverGrpc{
			RecipientList:[]string{"2","3"},
		}, func() {})
	}

	time.Sleep(10*time.Second)

}
