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
	"github.com/straightway/straightway/general/duration"
	"github.com/straightway/straightway/mocked"
	"github.com/straightway/straightway/peer"
	"github.com/straightway/straightway/simc"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SimulationNetworkPeerConnector_Test struct {
	suite.Suite
	sut        *simc.NetworkPeerConnector
	wrapped    *mocked.PeerConnector
	scheduler  *simc.EventScheduler
	other      *mocked.PeerConnector
	rawStorage *mocked.RawStorage
}

func TestSimulationNetworkPeerConnector(t *testing.T) {
	suite.Run(t, new(SimulationNetworkPeerConnector_Test))
}

func (suite *SimulationNetworkPeerConnector_Test) SetupTest() {
	suite.scheduler = &simc.EventScheduler{}
	suite.scheduler.Schedule(duration.Parse("1000h"), func() { suite.scheduler.Stop() })
	suite.wrapped = mocked.NewPeerConnector()
	suite.rawStorage = mocked.NewRawStorage(suite.scheduler)
	suite.other = mocked.NewPeerConnector()
	suite.sut = suite.wrap(suite.wrapped)
}

func (suite *SimulationNetworkPeerConnector_Test) TearDownTest() {
	suite.scheduler = nil
	suite.wrapped = nil
	suite.sut = nil
	suite.other = nil
}

// Tests

func (suite *SimulationNetworkPeerConnector_Test) Test_Id_IsTakenFromWrappedInstance() {
	id := suite.sut.Id()
	suite.wrapped.AssertCalledOnce(suite.T(), "Id")
	suite.Assert().Equal(suite.wrapped.Id(), id)
}

func (suite *SimulationNetworkPeerConnector_Test) Test_Equal_IsTakenFromWrappedInstance() {
	equalPeer := mocked.NewEqualer(false)
	result := suite.sut.Equal(equalPeer)
	suite.wrapped.AssertCalledOnce(suite.T(), "Equal", equalPeer)
	suite.Assert().False(result)
}

func (suite *SimulationNetworkPeerConnector_Test) Test_Push_IsDelayedByNetworkLatencyAndBandwidth() {
	size := float64(suite.rawStorage.SizeOf(&untimedChunk))
	sendTimeSeconds := time.Duration(size/suite.sut.Bandwidth) * time.Second
	expectedDuration := sendTimeSeconds + suite.sut.Latency
	suite.sut.Push(&untimedChunk, suite.other)
	suite.assertCallOnWrappedConnectorAfter(expectedDuration, "Push", &untimedChunk, mock.Anything)
}

func (suite *SimulationNetworkPeerConnector_Test) Test_Push_SenderIsWrapped() {
	suite.sut.Push(&untimedChunk, suite.other)
	suite.assertParameterWrappedOncePeer("Push", 2, 1)
}

func (suite *SimulationNetworkPeerConnector_Test) Test_Push_SenderIsNotReWrapped() {
	wrappedOrigin := suite.wrap(suite.other)
	suite.sut.Push(&untimedChunk, wrappedOrigin)
	suite.assertParameterWrappedOncePeer("Push", 2, 1)
}

func (suite *SimulationNetworkPeerConnector_Test) Test_Push_SenderNotWrappedIfNotAPeerConnector() {
	origin := mocked.NewIdHolder("123")
	suite.sut.Push(&untimedChunk, origin)
	suite.scheduler.Run()
	suite.wrapped.AssertCalledOnce(suite.T(), "Push", &untimedChunk, origin)
}

func (suite *SimulationNetworkPeerConnector_Test) Test_Query_IsDelayedByNetworkLatency() {
	query := data.Query{Id: queryId}
	suite.sut.Query(query, suite.other)
	suite.assertCallOnWrappedConnectorAfter(suite.sut.Latency, "Query", query, mock.Anything)
}

func (suite *SimulationNetworkPeerConnector_Test) Test_Query_ReceiverIsWrapped() {
	query := data.Query{Id: queryId}
	suite.sut.Query(query, suite.other)
	suite.assertParameterWrappedOncePeer("Query", 2, 1)
}

func (suite *SimulationNetworkPeerConnector_Test) Test_Query_ReceiverIsNotReWrapped() {
	query := data.Query{Id: queryId}
	suite.sut.Query(query, suite.wrap(suite.other))
	suite.assertParameterWrappedOncePeer("Query", 2, 1)
}

func (suite *SimulationNetworkPeerConnector_Test) Test_Query_ReceiverNotWrappedIfNotAPeerConnector() {
	receiver := mocked.NewPusherWithId("123")
	query := data.Query{Id: queryId}
	suite.sut.Query(query, receiver)
	suite.scheduler.Run()
	suite.wrapped.AssertCalledOnce(suite.T(), "Query", query, receiver)
}

func (suite *SimulationNetworkPeerConnector_Test) Test_RequestConnectionWith_IsDelayedByNetworkLatency() {
	suite.sut.RequestConnectionWith(suite.other)
	suite.assertCallOnWrappedConnectorAfter(suite.sut.Latency, "RequestConnectionWith", mock.Anything)
}

func (suite *SimulationNetworkPeerConnector_Test) Test_RequestConnectionWith_ReceiverIsWrapped() {
	suite.sut.RequestConnectionWith(suite.other)
	suite.assertParameterWrappedOncePeer("RequestConnectionWith", 1, 0)
}

func (suite *SimulationNetworkPeerConnector_Test) Test_RequestConnectionWith_ReceiverIsNotReWrapped() {
	suite.sut.RequestConnectionWith(suite.wrap(suite.other))
	suite.assertParameterWrappedOncePeer("RequestConnectionWith", 1, 0)
}

func (suite *SimulationNetworkPeerConnector_Test) Test_CloseConnectionWith_IsDelayedByNetworkLatency() {
	suite.sut.CloseConnectionWith(suite.other)
	suite.assertCallOnWrappedConnectorAfter(suite.sut.Latency, "CloseConnectionWith", mock.Anything)
}

func (suite *SimulationNetworkPeerConnector_Test) Test_CloseConnectionWith_ReceiverIsWrapped() {
	suite.sut.CloseConnectionWith(suite.other)
	suite.assertParameterWrappedOncePeer("CloseConnectionWith", 1, 0)
}

func (suite *SimulationNetworkPeerConnector_Test) Test_CloseConnectionWith_ReceiverIsNotReWrapped() {
	suite.sut.CloseConnectionWith(suite.wrap(suite.other))
	suite.assertParameterWrappedOncePeer("CloseConnectionWith", 1, 0)
}

func (suite *SimulationNetworkPeerConnector_Test) Test_RequestPeers_IsDelayedByNetworkLatency() {
	suite.sut.RequestPeers(suite.other)
	suite.assertCallOnWrappedConnectorAfter(suite.sut.Latency, "RequestPeers", mock.Anything)
}

func (suite *SimulationNetworkPeerConnector_Test) Test_RequestPeers_ReceiverIsWrapped() {
	suite.sut.RequestPeers(suite.other)
	suite.assertParameterWrappedOncePeer("RequestPeers", 1, 0)
}

func (suite *SimulationNetworkPeerConnector_Test) Test_RequestPeers_ReceiverIsNotReWrapped() {
	suite.sut.RequestPeers(suite.wrap(suite.other))
	suite.assertParameterWrappedOncePeer("RequestPeers", 1, 0)
}

func (suite *SimulationNetworkPeerConnector_Test) Test_AnnouncePeers_IsDelayedByNetworkLatencyAndBandwidth() {
	peers := []peer.Connector{suite.other, suite.other, suite.other}
	sizePerPeer := len(suite.other.Id())
	allPeersSize := float64(len(peers) * sizePerPeer)
	sendTimeSeconds := time.Duration(allPeersSize/suite.sut.Bandwidth) * time.Second
	suite.sut.AnnouncePeers(peers)
	suite.assertCallOnWrappedConnectorAfter(suite.sut.Latency+sendTimeSeconds, "AnnouncePeers", mock.Anything)
}

func (suite *SimulationNetworkPeerConnector_Test) Test_AnnouncePeers_ReceiverIsWrapped() {
	peers := []peer.Connector{suite.other, suite.other, suite.other}
	suite.sut.AnnouncePeers(peers)
	suite.assertParameterWrappedOncePeer("AnnouncePeers", 1, 0)
}

func (suite *SimulationNetworkPeerConnector_Test) Test_AnnouncePeers_ReceiverIsNotReWrapped() {
	peers := []peer.Connector{suite.other, suite.other, suite.other}
	suite.sut.AnnouncePeers(peers)
	suite.assertParameterWrappedOncePeer("AnnouncePeers", 1, 0)
}

// Private

func (suite *SimulationNetworkPeerConnector_Test) wrap(toWrap peer.Connector) *simc.NetworkPeerConnector {
	return &simc.NetworkPeerConnector{
		Wrapped:        toWrap,
		EventScheduler: suite.scheduler,
		RawStorage:     suite.rawStorage,
		Latency:        duration.Parse("1s"),
		Bandwidth:      2.0}
}

func (suite *SimulationNetworkPeerConnector_Test) assertParameterWrappedOncePeer(
	method string, numParams int, checkedParamIndex int) {
	params := make([]interface{}, numParams, numParams)
	for i := range params {
		params[i] = mock.Anything
	}
	suite.wrapped.OnNew(method, params...).Run(
		func(args mock.Arguments) {
			for _, p := range suite.args(args.Get(checkedParamIndex)) {
				wrapped, ok := p.(*simc.NetworkPeerConnector)
				suite.Assert().True(ok)
				_, ok = wrapped.Wrapped.(*simc.NetworkPeerConnector)
				suite.Assert().False(ok)
			}
		})
	suite.scheduler.Run()
	suite.wrapped.AssertCalledOnce(suite.T(), method, mock.Anything, mock.Anything)
}

func (suite *SimulationNetworkPeerConnector_Test) args(args interface{}) []interface{} {
	var result []interface{}
	peerArray, isPeerArray := args.([]peer.Connector)
	if isPeerArray {
		for _, p := range peerArray {
			result = append(result, p)
		}
	} else {
		result = append(result, args)
	}

	return result
}

func (suite *SimulationNetworkPeerConnector_Test) assertCallOnWrappedConnectorAfter(
	delay time.Duration,
	method string,
	parameters ...interface{}) {
	anyParameters := make([]interface{}, len(parameters), len(parameters))
	for i := range anyParameters {
		anyParameters[i] = mock.Anything
	}
	suite.wrapped.AssertNotCalled(suite.T(), method, anyParameters...)
	expectedExecutionTime := suite.scheduler.Time().Add(delay)
	var actualExecutionTime time.Time
	suite.wrapped.OnNew(method, anyParameters...).Run(func(mock.Arguments) {
		suite.Assert().Zero(actualExecutionTime)
		actualExecutionTime = suite.scheduler.Time()
	})
	suite.scheduler.Run()
	suite.wrapped.AssertCalledOnce(suite.T(), method, parameters...)
	suite.Assert().Equal(expectedExecutionTime, actualExecutionTime)
}