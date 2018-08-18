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

package mock_test

import (
	"testing"

	"github.com/it-chain/engine/common/command"
	"github.com/magiconair/properties/assert"
	"time"
	"github.com/it-chain/avengers/mock"
)

func TestNewNetworkManager(t *testing.T) {
}

func TestNetworkManager_AddProcess(t *testing.T) {

}

func TestNetworkManager_GrpcCall(t *testing.T) {
	tests := map[string]struct {
		input struct {
			RecipientList []string
			Protocol      string
		}
	}{
		"success": {input: struct {
			RecipientList []string
			Protocol      string
		}{RecipientList: []string{"1", "2"}, Protocol: "test"}},
		"no receiver test": {input: struct {
			RecipientList []string
			Protocol      string
		}{RecipientList: []string{}, Protocol: "test"}},
	}

	for testName, test := range tests {
		networkManager := mock.NewNetworkManager()
		t.Logf("running test case %s", testName)

		deliverGrpc := &command.DeliverGrpc{
			RecipientList: test.input.RecipientList,
			Protocol:      test.input.Protocol,
		}
		networkManager.GrpcCall("1","message.deliver", *deliverGrpc, func() {})
		t.Logf("end of test calling")
		for _, processId := range test.input.RecipientList {
			t.Logf("processId:%s is receiving", processId)
			go func(processId string) {
				a := <-networkManager.ChannelMap[processId]["message.receive"]
				assert.Equal(t, a.(command.ReceiveGrpc).Protocol, test.input.Protocol)
			}(processId)
		}
	}
}

func TestNetworkManager_GrpcConsume(t *testing.T) {
	callbackIndex := 1

	tests := map[string]struct {
		input struct {
			RecipientList []string
			ProcessId     string
			handler       func(a interface{}) error
		}
	}{
		"success": {input: struct {
			RecipientList []string
			ProcessId     string
			handler       func(a interface{}) error
		}{RecipientList: []string{"1", "2", "3"},
		ProcessId: "1",
		handler: func(a interface{}) error {callbackIndex=2; t.Logf("handler!"); return nil }}},
		//"do not receive": {input: struct {
		//	RecipientList []string
		//	ProcessId     string
		//	handler       func(a interface{}) error
		//}{RecipientList: []string{"2", "3"},
		//ProcessId: "1",
		//handler: func(a interface{}) error {callbackIndex=2; t.Logf("handler!"); return nil }}},
	}

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		networkManager := mock.NewNetworkManager()

		deliverGrpc := &command.DeliverGrpc{
			RecipientList: test.input.RecipientList,
		}
		networkManager.GrpcCall("1","message.deliver", *deliverGrpc, func() {})
		t.Logf("end of calling!")
		networkManager.GrpcConsume(test.input.ProcessId, "message.receive", test.input.handler)

	}

	time.Sleep(4 * time.Second)
	assert.Equal(t, callbackIndex, 2)
}
