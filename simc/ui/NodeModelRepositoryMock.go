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

package ui

import (
	"github.com/stretchr/testify/mock"

	"github.com/straightway/straightway/general/id"
	"github.com/straightway/straightway/general/mocked"
)

type NodeModelRepositoryMock struct {
	mocked.Base
	Nodes map[id.Type]NodeModel
}

func NewNodeModelRepositoryMock(nodes ...NodeModel) *NodeModelRepositoryMock {
	result := &NodeModelRepositoryMock{Nodes: make(map[id.Type]NodeModel)}
	result.On("NodeModelForId", mock.Anything).Return()
	return result
}

func (m *NodeModelRepositoryMock) NodeModelForId(id id.Type) NodeModel {
	m.Called(id)
	return m.Nodes[id]
}
