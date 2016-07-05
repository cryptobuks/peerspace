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

	"github.com/stretchr/testify/suite"
)

// Test suite

type ConnectionStrategy_IsConnectionAcceptedWith_Test struct {
	ConnectionStrategy_TestBase
}

func TestConnectionStrategy_IsConnectionAcceptedWith(t *testing.T) {
	suite.Run(t, new(ConnectionStrategy_IsConnectionAcceptedWith_Test))
}

// Tests

func (suite *ConnectionStrategy_IsConnectionAcceptedWith_Test) TestConnectionIsRefusedIfMaxConnectedIsReached() {
	suite.configuration.MaxConnections = 3
	for i := 0; i < suite.configuration.MaxConnections; i++ {
		suite.addConnectedPeer()
	}
	suite.Assert().False(suite.sut.IsConnectionAcceptedWith(suite.createPeerConnector()))
}

func (suite *ConnectionStrategy_IsConnectionAcceptedWith_Test) TestConnectionIsRefusedIfMaxConnectingIsReached() {
	suite.configuration.MaxConnections = 3
	for i := 0; i < suite.configuration.MaxConnections; i++ {
		suite.addConnectingPeer()
	}
	suite.Assert().False(suite.sut.IsConnectionAcceptedWith(suite.createPeerConnector()))
}

func (suite *ConnectionStrategy_IsConnectionAcceptedWith_Test) TestConnectionIsRefusedIfMaxConnectionsIsReached() {
	suite.configuration.MaxConnections = 3
	suite.addConnectedPeer()
	for i := 0; i < suite.configuration.MaxConnections-1; i++ {
		suite.addConnectingPeer()
	}
	suite.Assert().False(suite.sut.IsConnectionAcceptedWith(suite.createPeerConnector()))
}

func (suite *ConnectionStrategy_IsConnectionAcceptedWith_Test) TestConnectionIsAcceptedIfMaxConnectionIsNotReached() {
	suite.configuration.MaxConnections = 3
	suite.Assert().True(suite.sut.IsConnectionAcceptedWith(suite.createPeerConnector()))
}
