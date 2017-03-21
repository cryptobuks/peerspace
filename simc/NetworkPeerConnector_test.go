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

package simc

import (
	"testing"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/straightway/straightway/data"
	"github.com/straightway/straightway/general"
	"github.com/straightway/straightway/general/duration"
	"github.com/straightway/straightway/general/id"
	"github.com/straightway/straightway/peer"
	slog "github.com/straightway/straightway/simc/log"
)

type NetworkPeerConnector_Test struct {
	suite.Suite
	sut        *NetworkPeerConnector
	properties *NetworkProperties
	wrapped    *peer.ConnectorMock
	scheduler  *EventScheduler
	others     []*peer.ConnectorMock
	rawStorage *data.RawStorageMock
}

func TestNetworkPeerConnector(t *testing.T) {
	suite.Run(t, new(NetworkPeerConnector_Test))
}

func (suite *NetworkPeerConnector_Test) SetupTest() {
	log.SetHandler(discard.New())
	suite.scheduler = NewEventScheduler()
	suite.scheduler.Schedule(duration.Parse("1000h"), func() { suite.scheduler.Stop() })
	suite.wrapped = peer.NewConnectorMock()
	suite.rawStorage = data.NewRawStorageMock(suite.scheduler)
	suite.createOtherPeers(4)
	suite.properties = &NetworkProperties{
		EventScheduler: suite.scheduler,
		SizeOfer:       suite.rawStorage,
		Latency:        duration.Parse("1s"),
		Bandwidth:      2.0}
	suite.sut = suite.wrap(suite.wrapped)
	suite.wrapped.Calls = nil
}

func (suite *NetworkPeerConnector_Test) TearDownTest() {
	suite.scheduler = nil
	suite.wrapped = nil
	suite.sut = nil
	suite.others = nil
}

// Tests

func (suite *NetworkPeerConnector_Test) Test_Id_IsTakenFromWrappedInstance() {
	id := suite.sut.Id()
	suite.wrapped.AssertCalledOnce(suite.T(), "Id")
	suite.Assert().Equal(suite.wrapped.Id(), id)
}

func (suite *NetworkPeerConnector_Test) Test_Equal_IsTakenFromWrappedInstance() {
	equalPeer := general.NewEqualerMock(false)
	result := suite.sut.Equal(equalPeer)
	suite.wrapped.AssertCalledOnce(suite.T(), "Equal", equalPeer)
	suite.Assert().False(result)
}

func (suite *NetworkPeerConnector_Test) Test_Push_IsDelayedByNetworkLatencyAndBandwidth() {
	size := float64(suite.rawStorage.SizeOf(&data.UntimedChunk))
	sendTimeSeconds := time.Duration(size/suite.sut.Bandwidth()) * time.Second
	expectedDuration := sendTimeSeconds + suite.sut.Latency()
	suite.sut.Push(&data.UntimedChunk, suite.others[0])
	suite.assertCallOnWrappedConnectorAfter(expectedDuration, "Push", &data.UntimedChunk, mock.Anything)
}

func (suite *NetworkPeerConnector_Test) Test_Push_IsDelayedByPreviousNetworkAction() {
	suite.sut.RequestPeers(suite.others[0])
	size := float64(suite.rawStorage.SizeOf(&data.UntimedChunk))
	sendTimeSeconds := time.Duration(size/suite.sut.Bandwidth()) * time.Second
	expectedDuration := sendTimeSeconds + 2*suite.sut.Latency()
	suite.sut.Push(&data.UntimedChunk, suite.others[0])
	suite.assertCallOnWrappedConnectorAfter(expectedDuration, "Push", &data.UntimedChunk, mock.Anything)
}

func (suite *NetworkPeerConnector_Test) Test_Push_IsNotDelayedByOtherPeersNetworkAction() {
	suite.sut.RequestPeers(suite.others[1])
	size := float64(suite.rawStorage.SizeOf(&data.UntimedChunk))
	sendTimeSeconds := time.Duration(size/suite.sut.Bandwidth()) * time.Second
	expectedDuration := sendTimeSeconds + suite.sut.Latency()
	suite.sut.Push(&data.UntimedChunk, suite.others[0])
	suite.assertCallOnWrappedConnectorAfter(expectedDuration, "Push", &data.UntimedChunk, mock.Anything)
}

func (suite *NetworkPeerConnector_Test) Test_Push_SenderIsWrapped() {
	suite.sut.Push(&data.UntimedChunk, suite.others[0])
	suite.assertParameterWrappedOncePeer("Push", 2, 1)
}

func (suite *NetworkPeerConnector_Test) Test_Push_SenderIsNotReWrapped() {
	wrappedOrigin := suite.wrap(suite.others[0])
	suite.sut.Push(&data.UntimedChunk, wrappedOrigin)
	suite.assertParameterWrappedOncePeer("Push", 2, 1)
}

func (suite *NetworkPeerConnector_Test) Test_Push_SenderNotWrappedIfNotAPeerConnector() {
	origin := id.NewHolderMock(id.FromString("123"))
	suite.sut.Push(&data.UntimedChunk, origin)
	suite.scheduler.Run()
	suite.wrapped.AssertCalledOnce(suite.T(), "Push", &data.UntimedChunk, origin)
}

func (suite *NetworkPeerConnector_Test) Test_Query_IsDelayedByNetworkLatency() {
	query := data.Query{Id: data.QueryId}
	suite.sut.Query(query, suite.others[0])
	suite.assertCallOnWrappedConnectorAfter(suite.sut.Latency(), "Query", query, mock.Anything)
}

func (suite *NetworkPeerConnector_Test) Test_Query_IsDelayedByPreviousNetworkAction() {
	suite.sut.RequestPeers(suite.others[0])
	query := data.Query{Id: data.QueryId}
	suite.sut.Query(query, suite.others[0])
	suite.assertCallOnWrappedConnectorAfter(2*suite.sut.Latency(), "Query", query, mock.Anything)
}

func (suite *NetworkPeerConnector_Test) Test_Query_IsNotDelayedByOtherPeerNetworkAction() {
	suite.sut.RequestPeers(suite.others[1])
	query := data.Query{Id: data.QueryId}
	suite.sut.Query(query, suite.others[0])
	suite.assertCallOnWrappedConnectorAfter(suite.sut.Latency(), "Query", query, mock.Anything)
}

func (suite *NetworkPeerConnector_Test) Test_Query_ReceiverIsWrapped() {
	query := data.Query{Id: data.QueryId}
	suite.sut.Query(query, suite.others[0])
	suite.assertParameterWrappedOncePeer("Query", 2, 1)
}

func (suite *NetworkPeerConnector_Test) Test_Query_ReceiverIsNotReWrapped() {
	query := data.Query{Id: data.QueryId}
	suite.sut.Query(query, suite.wrap(suite.others[0]))
	suite.assertParameterWrappedOncePeer("Query", 2, 1)
}

func (suite *NetworkPeerConnector_Test) Test_Query_ReceiverNotWrappedIfNotAPeerConnector() {
	receiver := peer.NewPusherWithIdMock(id.FromString("123"))
	query := data.Query{Id: data.QueryId}
	suite.sut.Query(query, receiver)
	suite.scheduler.Run()
	suite.wrapped.AssertCalledOnce(suite.T(), "Query", query, receiver)
}

func (suite *NetworkPeerConnector_Test) Test_RequestConnectionWith_IsDelayedByNetworkLatency() {
	suite.sut.RequestConnectionWith(suite.others[0])
	suite.assertCallOnWrappedConnectorAfter(suite.sut.Latency(), "RequestConnectionWith", mock.Anything)
}

func (suite *NetworkPeerConnector_Test) Test_RequestConnectionWith_IsDelayedByPreviousNetworkAction() {
	suite.sut.RequestPeers(suite.others[0])
	suite.sut.RequestConnectionWith(suite.others[0])
	suite.assertCallOnWrappedConnectorAfter(2*suite.sut.Latency(), "RequestConnectionWith", mock.Anything)
}

func (suite *NetworkPeerConnector_Test) Test_RequestConnectionWith_IsNotDelayedByOtherNodeNetworkAction() {
	suite.sut.RequestPeers(suite.others[1])
	suite.sut.RequestConnectionWith(suite.others[0])
	suite.assertCallOnWrappedConnectorAfter(suite.sut.Latency(), "RequestConnectionWith", mock.Anything)
}

func (suite *NetworkPeerConnector_Test) Test_RequestConnectionWith_ReceiverIsWrapped() {
	suite.sut.RequestConnectionWith(suite.others[0])
	suite.assertParameterWrappedOncePeer("RequestConnectionWith", 1, 0)
}

func (suite *NetworkPeerConnector_Test) Test_RequestConnectionWith_ReceiverIsNotReWrapped() {
	suite.sut.RequestConnectionWith(suite.wrap(suite.others[0]))
	suite.assertParameterWrappedOncePeer("RequestConnectionWith", 1, 0)
}

func (suite *NetworkPeerConnector_Test) Test_CloseConnectionWith_IsDelayedByNetworkLatency() {
	suite.sut.CloseConnectionWith(suite.others[0])
	suite.assertCallOnWrappedConnectorAfter(suite.sut.Latency(), "CloseConnectionWith", mock.Anything)
}

func (suite *NetworkPeerConnector_Test) Test_CloseConnectionWith_IsDelayedByPreviousNetworkAction() {
	suite.sut.RequestPeers(suite.others[0])
	suite.sut.CloseConnectionWith(suite.others[0])
	suite.assertCallOnWrappedConnectorAfter(2*suite.sut.Latency(), "CloseConnectionWith", mock.Anything)
}

func (suite *NetworkPeerConnector_Test) Test_CloseConnectionWith_IsNotDelayedByOtherPeerNetworkAction() {
	suite.sut.RequestPeers(suite.others[1])
	suite.sut.CloseConnectionWith(suite.others[0])
	suite.assertCallOnWrappedConnectorAfter(suite.sut.Latency(), "CloseConnectionWith", mock.Anything)
}

func (suite *NetworkPeerConnector_Test) Test_CloseConnectionWith_ReceiverIsWrapped() {
	suite.sut.CloseConnectionWith(suite.others[0])
	suite.assertParameterWrappedOncePeer("CloseConnectionWith", 1, 0)
}

func (suite *NetworkPeerConnector_Test) Test_CloseConnectionWith_ReceiverIsNotReWrapped() {
	suite.sut.CloseConnectionWith(suite.wrap(suite.others[0]))
	suite.assertParameterWrappedOncePeer("CloseConnectionWith", 1, 0)
}

func (suite *NetworkPeerConnector_Test) Test_RequestPeers_IsDelayedByNetworkLatency() {
	suite.sut.RequestPeers(suite.others[0])
	suite.assertCallOnWrappedConnectorAfter(suite.sut.Latency(), "RequestPeers", mock.Anything)
}

func (suite *NetworkPeerConnector_Test) Test_RequestPeers_IsDelayedByPreviousNetworkAction() {
	suite.sut.AnnouncePeersFrom(suite.others[0], nil)
	suite.sut.AnnouncePeersFrom(suite.others[0], nil)
	suite.sut.RequestPeers(suite.others[0])
	suite.assertCallOnWrappedConnectorAfter(3*suite.sut.Latency(), "RequestPeers", mock.Anything)
}

func (suite *NetworkPeerConnector_Test) Test_RequestPeers_IsNotDelayedByOtherPeerNetworkAction() {
	suite.sut.AnnouncePeersFrom(suite.others[1], nil)
	suite.sut.RequestPeers(suite.others[0])
	suite.assertCallOnWrappedConnectorAfter(suite.sut.Latency(), "RequestPeers", mock.Anything)
}

func (suite *NetworkPeerConnector_Test) Test_RequestPeers_ReceiverIsWrapped() {
	suite.sut.RequestPeers(suite.others[0])
	suite.assertParameterWrappedOncePeer("RequestPeers", 1, 0)
}

func (suite *NetworkPeerConnector_Test) Test_RequestPeers_ReceiverIsNotReWrapped() {
	suite.sut.RequestPeers(suite.wrap(suite.others[0]))
	suite.assertParameterWrappedOncePeer("RequestPeers", 1, 0)
}

func (suite *NetworkPeerConnector_Test) Test_AnnouncePeersFrom_IsDelayedByNetworkLatencyAndBandwidth() {
	peers := []peer.Connector{suite.others[1], suite.others[2], suite.others[3]}
	sizePerPeer := len(suite.others[1].Id())
	allPeersSize := float64(len(peers) * sizePerPeer)
	sendTimeSeconds := time.Duration(allPeersSize/suite.sut.Bandwidth()) * time.Second
	suite.sut.AnnouncePeersFrom(suite.others[0], peers)
	suite.assertCallOnWrappedConnectorAfter(
		suite.sut.Latency()+sendTimeSeconds,
		"AnnouncePeersFrom",
		mock.Anything,
		mock.Anything)
}

func (suite *NetworkPeerConnector_Test) Test_AnnouncePeersFrom_IsDelayedByPreviousNetworkAction() {
	suite.sut.RequestPeers(suite.others[0])
	peers := []peer.Connector{suite.others[1], suite.others[2], suite.others[3]}
	sizePerPeer := len(suite.others[1].Id())
	allPeersSize := float64(len(peers) * sizePerPeer)
	sendTimeSeconds := time.Duration(allPeersSize/suite.sut.Bandwidth()) * time.Second
	suite.sut.AnnouncePeersFrom(suite.others[0], peers)
	suite.assertCallOnWrappedConnectorAfter(
		2*suite.sut.Latency()+sendTimeSeconds,
		"AnnouncePeersFrom",
		mock.Anything,
		mock.Anything)
}

func (suite *NetworkPeerConnector_Test) Test_AnnouncePeersFrom_IsNotDelayedByOtherPeerNetworkAction() {
	suite.sut.RequestPeers(suite.others[1])
	peers := []peer.Connector{suite.others[1], suite.others[2], suite.others[3]}
	sizePerPeer := len(suite.others[1].Id())
	allPeersSize := float64(len(peers) * sizePerPeer)
	sendTimeSeconds := time.Duration(allPeersSize/suite.sut.Bandwidth()) * time.Second
	suite.sut.AnnouncePeersFrom(suite.others[0], peers)
	suite.assertCallOnWrappedConnectorAfter(
		suite.sut.Latency()+sendTimeSeconds,
		"AnnouncePeersFrom",
		mock.Anything,
		mock.Anything)
}

func (suite *NetworkPeerConnector_Test) Test_AnnouncePeersFrom_ReceiverIsWrapped() {
	peers := []peer.Connector{suite.others[1], suite.others[2], suite.others[3]}
	suite.sut.AnnouncePeersFrom(suite.others[0], peers)
	suite.assertParameterWrappedOncePeer("AnnouncePeersFrom", 2, 0)
}

func (suite *NetworkPeerConnector_Test) Test_AnnouncePeersFrom_ReceiverIsNotReWrapped() {
	peers := []peer.Connector{suite.others[1], suite.others[2], suite.others[3]}
	suite.sut.AnnouncePeersFrom(suite.others[0], peers)
	suite.assertParameterWrappedOncePeer("AnnouncePeersFrom", 2, 0)
}

func (suite *NetworkPeerConnector_Test) Test_AnnouncePeersFrom_PeersAreWrapped() {
	peers := []peer.Connector{suite.others[1], suite.others[2], suite.others[3]}
	suite.sut.AnnouncePeersFrom(suite.others[0], peers)
	suite.assertParameterWrappedOncePeer("AnnouncePeersFrom", 2, 1)
}

func (suite *NetworkPeerConnector_Test) Test_AnnouncePeersFrom_PeersAreNotReWrapped() {
	peers := []peer.Connector{suite.others[1], suite.others[2], suite.others[3]}
	suite.sut.AnnouncePeersFrom(suite.others[0], peers)
	suite.assertParameterWrappedOncePeer("AnnouncePeersFrom", 2, 1)
}

func (suite *NetworkPeerConnector_Test) Test_Wrapping_AlwaysYieldsTheSameWrappedInstance() {
	otherWrapped := suite.wrap(suite.wrapped)
	otherWrapped.AnnouncePeersFrom(suite.others[0], nil)
	suite.sut.RequestPeers(suite.others[0])
	suite.assertCallOnWrappedConnectorAfter(2*suite.sut.Latency(), "RequestPeers", mock.Anything)
}

func (suite *NetworkPeerConnector_Test) Test_NothingIsLogged_WhenLogIsDisabled() {
	logHandler := slog.NewHandlerMock()
	defer log.SetHandler(discard.New())
	log.SetHandler(logHandler)
	suite.sut.CloseConnectionWith(suite.others[0])
	logHandler.AssertNotCalled(suite.T(), "HandleLog", mock.Anything)
	suite.scheduler.Run()
	logHandler.AssertNotCalled(suite.T(), "HandleLog", mock.Anything)
}

func (suite *NetworkPeerConnector_Test) Test_SomethingIsLogged_WhenLogIsEnabled() {
	logHandler := slog.NewHandlerMock()
	defer log.SetHandler(discard.New())
	log.SetHandler(logHandler)
	slog.SetEnabled(true)
	defer slog.SetEnabled(false)
	suite.sut.CloseConnectionWith(suite.others[0])
	logHandler.AssertCalled(suite.T(), "HandleLog", mock.Anything)
	logHandler.Calls = nil
	suite.scheduler.Run()
	logHandler.AssertCalled(suite.T(), "HandleLog", mock.Anything)
}

func (suite *NetworkPeerConnector_Test) Test_SubLogsHaveNoSubSubLogs() {
	logHandler := &slog.HandlerMock{}
	logHandler.On("HandleLog", mock.Anything).Run(func(args mock.Arguments) {
		logEntry := args.Get(0).(*log.Entry)
		suite.Assert().False(hasSubSubFields(logEntry))
	})
	defer log.SetHandler(discard.New())
	log.SetHandler(logHandler)
	slog.SetEnabled(true)
	defer slog.SetEnabled(false)

	suite.sut.RequestConnectionWith(suite.others[0])
	suite.sut.CloseConnectionWith(suite.others[0])
}

// Private

func hasSubSubFields(entry *log.Entry) bool {
	for _, logField := range entry.Fields {
		subFields, _ := logField.([]log.Fields)
		for _, subField := range subFields {
			for _, subFieldEntry := range subField {
				_, hasSubSubFields := subFieldEntry.([]log.Fields)
				if hasSubSubFields {
					return true
				}
			}
		}
	}

	return false
}

func (suite *NetworkPeerConnector_Test) wrap(toWrap peer.Connector) *NetworkPeerConnector {
	return NewNetworkPeerConnector(toWrap, suite.properties)
}

func (suite *NetworkPeerConnector_Test) assertParameterWrappedOncePeer(
	method string, numParams int, checkedParamIndex int) {
	params := make([]interface{}, numParams, numParams)
	for i := range params {
		params[i] = mock.Anything
	}
	suite.wrapped.OnNew(method, params...).Run(
		func(args mock.Arguments) {
			for _, p := range suite.args(args.Get(checkedParamIndex)) {
				wrapped, ok := p.(*NetworkPeerConnector)
				suite.Assert().True(ok)
				_, ok = wrapped.Wrapped().(*NetworkPeerConnector)
				suite.Assert().False(ok)
			}
		})
	suite.scheduler.Run()
	suite.wrapped.AssertCalledOnce(suite.T(), method, mock.Anything, mock.Anything)
}

func (suite *NetworkPeerConnector_Test) args(args interface{}) []interface{} {
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

func (suite *NetworkPeerConnector_Test) assertCallOnWrappedConnectorAfter(
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

func (suite *NetworkPeerConnector_Test) createOtherPeers(num int) {
	suite.others = make([]*peer.ConnectorMock, num)
	for i := range suite.others {
		suite.others[i] = peer.NewConnectorMock()
	}
}
