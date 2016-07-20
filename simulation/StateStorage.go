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

package simulation

import (
	"github.com/straightway/straightway/general"
	"github.com/straightway/straightway/peer"
)

type StateStorage struct {
	connectors []peer.Connector
}

func (this *StateStorage) GetAllKnownPeers() []peer.Connector {
	return this.connectors
}

func (this *StateStorage) IsKnownPeer(peer peer.Connector) bool {
	return general.Contains(this.connectors, peer)
}

func (this *StateStorage) AddKnownPeer(peer peer.Connector) {
	if this.IsKnownPeer(peer) == false {
		this.connectors = append(this.connectors, peer)
	}
}
