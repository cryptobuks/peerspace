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

package strategy

import (
	"math/rand"

	"github.com/straightway/straightway/app"
	"github.com/straightway/straightway/data"
	"github.com/straightway/straightway/general"
	"github.com/straightway/straightway/general/id"
	"github.com/straightway/straightway/general/slice"
	"github.com/straightway/straightway/peer"
)

// The default connections strategy.
type Connection struct {
	Configuration          *app.Configuration
	ConnectionInfoProvider peer.ConnectionInfoProvider
	RandSource             rand.Source
	PeerDistanceCalculator PeerDistanceCalculator
}

var _ peer.ConnectionStrategy = (*Connection)(nil)

type peersSortedByDistance struct {
	peers                  []peer.Connector
	local                  peer.Connector
	peerDistanceCalculator PeerDistanceCalculator
}

func (this peersSortedByDistance) Len() int { return len(this.peers) }
func (this peersSortedByDistance) Swap(i, j int) {
	this.peers[i], this.peers[j] = this.peers[j], this.peers[i]
}
func (this peersSortedByDistance) Less(i, j int) bool {
	a := this.peers[i]
	distanceToA := this.peerDistanceCalculator.Distance(this.local, data.Key{a.Id(), 0})
	b := this.peers[j]
	distanceToB := this.peerDistanceCalculator.Distance(this.local, data.Key{b.Id(), 0})
	return distanceToA < distanceToB
}

func (this *Connection) PeersToConnect(allPeers []peer.Connector) []peer.Connector {
	var maxConn = this.Configuration.MaxConnections
	filteredPeers := this.existingConnectionsFilteredFrom(allPeers)
	if len(filteredPeers) <= maxConn {
		return filteredPeers
	} else {
		permutation := rand.New(this.RandSource).Perm(len(filteredPeers))
		result := make([]peer.Connector, maxConn, maxConn)
		for i := 0; i < maxConn; i++ {
			result[i] = filteredPeers[permutation[i]]
		}

		return result
	}
}

func (this *Connection) IsConnectionAcceptedWith(peer id.Holder) bool {
	return len(this.existingConnections()) < this.Configuration.MaxConnections
}

func (this *Connection) existingConnectionsFilteredFrom(allPeers []peer.Connector) []peer.Connector {
	filteredPeers := append([]peer.Connector(nil), allPeers...)
	omittedPeers := this.existingConnections()
	return slice.RemoveItemsIf(filteredPeers, func(item interface{}) bool {
		return slice.Contains(omittedPeers, item.(general.Equaler))
	}).([]peer.Connector)
}

func (this *Connection) existingConnections() []peer.Connector {
	connected := this.ConnectionInfoProvider.ConnectedPeers()
	return append(connected, this.ConnectionInfoProvider.ConnectingPeers()...)
}
