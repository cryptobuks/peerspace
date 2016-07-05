/****************************************************************************
   Copyright 2016 github.com/straightway

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
****************************************************************************/

package mocked

import (
	"time"

	"github.com/straightway/straightway/general"
	"github.com/straightway/straightway/peer"
	"github.com/stretchr/testify/mock"
)

type QueryStrategy struct {
	Base
}

func NewQueryForwardStrategy(resultPeers []peer.Connector) *QueryStrategy {
	result := &QueryStrategy{}
	result.On("ForwardTargetsFor", mock.Anything, mock.Anything, mock.Anything).Return(resultPeers)
	result.On("TimeoutFor", mock.Anything).Return(general.ParseDuration("2h"))
	result.On("IsQueryAccepted", mock.Anything, mock.Anything).Return(true)
	return result
}

func (m *QueryStrategy) IsQueryAccepted(query peer.Query, receiver peer.Pusher) bool {
	args := m.Called(query, receiver)
	return args.Get(0).(bool)
}

func (m *QueryStrategy) ForwardTargetsFor(allPeers []peer.Connector, receiver peer.Pusher, query peer.Query) []peer.Connector {
	args := m.Called(allPeers, receiver, query)
	return args.Get(0).([]peer.Connector)
}

func (m *QueryStrategy) TimeoutFor(query peer.Query) time.Duration {
	args := m.Called(query)
	return args.Get(0).(time.Duration)
}
