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

package test

import (
	"testing"
	"time"

	"github.com/straightway/straightway/data"
	"github.com/straightway/straightway/mocked"
	"github.com/straightway/straightway/peer"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Test suite

type Node_Query_Test struct {
	suite.Suite
	*NodeContext
	queryPeer *mocked.PeerConnector
}

func TestNodeQuery(t *testing.T) {
	suite.Run(t, new(Node_Query_Test))
}

func (suite *Node_Query_Test) SetupTest() {
	suite.NodeContext = NewNodeContext()
	suite.AddKnownConnectedPeer(DoForward(true))
	suite.queryStrategy.On("TimeoutFor", mock.Anything).Return(time.Duration(30) * time.Second)
	suite.queryPeer = mocked.CreatePeerConnector()
	suite.node.Startup()
}

func (suite *Node_Query_Test) TearDownTest() {
	suite.ShutDownNode()
	suite.NodeContext = nil
	suite.queryPeer = nil
}

// Tests

func (suite *Node_Query_Test) Test_Query_LocallyStoredItemIsPushedToQueryNode() {
	suite.SetDataStorage(mocked.NewDataStorage(&untimedChunk))
	suite.Query(peer.Query{Id: untimedKey.Id})
	suite.assertQueryResult(&untimedChunk)
}

func (suite *Node_Query_Test) Test_Query_NotLocallyStoredItemIsNotDirectlyPushedBack() {
	suite.Query(peer.Query{Id: untimedKey.Id})
	suite.assertQueryResult( /*nothing*/ )
}

func (suite *Node_Query_Test) Test_Query_LocallyFailedQueryIsForwarded() {
	fwdPeer := suite.AddKnownConnectedPeer(DoForward(true))
	query := peer.Query{Id: untimedKey.Id}
	suite.Query(query)
	fwdPeer.AssertCalledOnce(suite.T(), "Query", query, suite.node)
}

func (suite *Node_Query_Test) Test_Query_ReceivedQueryResultIsForwardedOnce() {
	suite.AddKnownConnectedPeer(DoForward(true))
	suite.Query(peer.Query{Id: untimedKey.Id})
	suite.node.Push(&untimedChunk)
	suite.assertQueryResult(&untimedChunk)
	suite.node.Push(&untimedChunk)
	suite.assertQueryResult( /*nothing*/ )
}

func (suite *Node_Query_Test) Test_Query_ReceivedQueryResultIsForwardedToMultipleReceivers() {
	otherQueryPeer := mocked.CreatePeerConnector()
	suite.AddKnownConnectedPeer(DoForward(true))
	suite.Query(peer.Query{Id: untimedKey.Id})
	suite.node.Query(peer.Query{Id: untimedKey.Id}, otherQueryPeer)

	suite.node.Push(&untimedChunk)

	suite.assertQueryResult(&untimedChunk)
	AssertPushed(suite.T(), otherQueryPeer, &untimedChunk)
}

func (suite *Node_Query_Test) Test_Query_ResultIsSentOnceIfPeerIsAlsoForwardTarget() {
	suite.queryPeer = suite.AddKnownConnectedPeer(DoForward(true))
	suite.Query(peer.Query{Id: untimedKey.Id})
	suite.node.Push(&untimedChunk)
	suite.assertQueryResult(&untimedChunk)
}

func (suite *Node_Query_Test) Test_Query_IsDiscardedAfterTimeout() {
	suite.AddKnownConnectedPeer(DoForward(true))
	query := peer.Query{Id: untimedKey.Id}
	suite.Query(query)
	suite.advanceTimeByQueryTimeoutFactor(query, 1.1)
	suite.clearTimedOutQueries()

	suite.node.Push(&untimedChunk)

	suite.assertQueryResult( /*nothing*/ )
}

func (suite *Node_Query_Test) Test_Query_IsNotDiscardedBeforeTimeout() {
	suite.AddKnownConnectedPeer(DoForward(true))
	query := peer.Query{Id: untimedKey.Id}
	suite.Query(query)
	suite.advanceTimeByQueryTimeoutFactor(query, 0.5)
	suite.clearTimedOutQueries()

	suite.node.Push(&untimedChunk)

	suite.assertQueryResult(&untimedChunk)
}

func (suite *Node_Query_Test) Test_Query_ForTimeStampRangeYieldsMultipleResults() {
	suite.Query(peer.Query{Id: queryId, TimeFrom: 10, TimeTo: 20})
	suite.node.Push(&timedChunk10)
	suite.node.Push(&timedChunk20)
	suite.assertQueryResult(&timedChunk10, &timedChunk20)
}

func (suite *Node_Query_Test) Test_Query_ForTimeStampRangeIsDiscardedAfterTimeout() {
	query := peer.Query{Id: queryId, TimeFrom: 10, TimeTo: 20}
	suite.Query(query)
	suite.advanceTimeByQueryTimeoutFactor(query, 1.1)
	suite.node.Push(&timedChunk10)
	suite.node.Push(&timedChunk20)
	suite.assertQueryResult(&timedChunk10)
}

// Private

func (suite *Node_Query_Test) advanceTimeByQueryTimeoutFactor(query peer.Query, factor float32) {
	timeout := suite.queryStrategy.TimeoutFor(query)
	suite.AdvanceTimeBy(time.Duration(float32(timeout) * factor))
}

func (suite *Node_Query_Test) clearTimedOutQueries() {
	suite.node.Push(&data.Chunk{Key: data.Key{Id: "Other Key"}})
}

func (suite *Node_Query_Test) assertQueryResult(chunks ...*data.Chunk) {
	AssertPushed(suite.T(), suite.queryPeer, chunks...)
	suite.queryPeer.Calls = nil
}

func (suite *Node_Query_Test) Query(query peer.Query) {
	suite.node.Query(query, suite.queryPeer)
}
