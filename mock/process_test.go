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
	"github.com/it-chain/avengers/mock"
	"github.com/it-chain/engine/common/command"
	"github.com/magiconair/properties/assert"
)

func TestProcess_GrpcListen(t *testing.T) {

}

func TestProcess_RegisterHandler(t *testing.T) {
	process:= mock.NewProcess()
	process.Init("1")
	handler := func(command command.ReceiveGrpc) error{
		return nil
	}
	handler2 := func(command command.ReceiveGrpc) error{
		return nil
	}
	process.RegisterHandler(handler)
	process.RegisterHandler(handler2)

	assert.Equal(t, len(process.GrpcCommandHandlers), 2)
}

func TestProcess_Init(t *testing.T) {

}