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

package env

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/straightway/straightway/general/id"
	"github.com/straightway/straightway/peer"
	"github.com/straightway/straightway/simc/ui"
)

type NodeModelTest struct {
	suite.Suite
	sut                    *NodeModel
	nodeModels             *ui.NodeModelRepositoryMock
	connectionInfoProvider *peer.ConnectionInfoProviderMock
}

type peerConnectionType int

const (
	unconnected = peerConnectionType(iota)
	connecting
	connected
)

var (
	defaultNodeModelId = id.FromString("defaultModelId")
	otherNodeModelId   = id.FromString("otherModelId")
)

func TestModeModel(t *testing.T) {
	suite.Run(t, new(NodeModelTest))
}

func (suite *NodeModelTest) SetupTest() {
	suite.nodeModels = ui.NewNodeModelRepositoryMock()
	suite.connectionInfoProvider = peer.NewConnectionInfoProviderMock()
	suite.sut = suite.createSut(defaultNodeModelId)
}

func (suite *NodeModelTest) TearDownTest() {
	suite.sut = nil
	suite.nodeModels = nil
	suite.connectionInfoProvider = nil
}

// Tests

func (suite *NodeModelTest) Test_Construction_HasId() {
	suite.Assert().Equal(defaultNodeModelId, suite.sut.Id())
}

func (suite *NodeModelTest) Test_Equal_TrueIfHasSameId() {
	sut2 := suite.createSut(defaultNodeModelId)
	suite.Assert().True(suite.sut.Equal(sut2))
}

func (suite *NodeModelTest) Test_Equal_FalseIfHasDifferentId() {
	sut2 := suite.createSut(otherNodeModelId)
	suite.Assert().False(suite.sut.Equal(sut2))
}

func (suite *NodeModelTest) Test_Equal_FalseIfHasDifferentType() {
	sut2 := ui.NewNodeModelMock(defaultNodeModelId, 0.0, 0.0)
	suite.Assert().False(suite.sut.Equal(sut2))
}

func (suite *NodeModelTest) Test_Position_IsInitially00() {
	x, y := suite.sut.Position()
	suite.Assert().Equal(0.0, x)
	suite.Assert().Equal(0.0, y)
}

func (suite *NodeModelTest) Test_Position_CanBeSet() {
	suite.sut.SetPosition(1.0, 1.0)
	x, y := suite.sut.Position()
	suite.Assert().Equal(1.0, x)
	suite.Assert().Equal(1.0, y)
}

func (suite *NodeModelTest) Test_Connections_InitiallyYieldsNoConnectedNodes() {
	result := suite.sut.Connections()
	suite.Assert().Empty(result)
}

func (suite *NodeModelTest) Test_Connections_YieldsConnectedNodesFromConnectionInfoProvider() {
	peers := suite.createPeers(connected, 2)
	result := suite.sut.Connections()

	suite.Assert().Equal(peers, result)
}

func (suite *NodeModelTest) Test_Connections_YieldsConnectingNodesFromConnectionInfoProvider() {
	peers := suite.createPeers(connecting, 2)
	result := suite.sut.Connections()

	suite.Assert().Equal(peers, result)
}
func (suite *NodeModelTest) Test_Connections_YieldsConnectingAndConnectedNodesFromConnectionInfoProvider() {
	peers := suite.createPeers(connected, 2)
	peers = append(peers, suite.createPeers(connecting, 2)...)
	result := suite.sut.Connections()

	suite.Assert().Equal(peers, result)
}

// Private

func (suite *NodeModelTest) createSut(id id.Type) *NodeModel {
	return NewNodeModel(id, suite.nodeModels, suite.connectionInfoProvider)
}

func (suite *NodeModelTest) createPeers(
	connectionType peerConnectionType,
	num int) []ui.NodeModel {

	result := make([]ui.NodeModel, num)
	for i := range result {
		result[i] = suite.createPeer(connectionType)
	}

	return result
}

func (suite *NodeModelTest) createPeer(connectionType peerConnectionType) *NodeModel {
	peerNode := peer.NewConnectorMock()
	switch connectionType {
	case connected:
		suite.connectionInfoProvider.AllConnectedPeers = append(
			suite.connectionInfoProvider.AllConnectedPeers,
			peerNode)
	case connecting:
		suite.connectionInfoProvider.AllConnectingPeers = append(
			suite.connectionInfoProvider.AllConnectingPeers,
			peerNode)
	}

	peer := suite.createSut(peerNode.Id())
	suite.nodeModels.Nodes[peer.Id()] = peer

	return peer
}