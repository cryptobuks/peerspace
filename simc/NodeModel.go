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
	"github.com/straightway/straightway/simc/ui"
	"github.com/straightway/straightway/strategy"
)

type NodeModel struct {
	nodeModels     NodeModelRepository
	connectionInfo strategy.ConnectionInfoProvider
	x, y           float64
}

func NewNodeModel(nodeModels NodeModelRepository, connectionInfo strategy.ConnectionInfoProvider) *NodeModel {
	result := &NodeModel{
		nodeModels:     nodeModels,
		connectionInfo: connectionInfo}
	return result
}

func (this *NodeModel) Position() (x, y float64) {
	return this.x, this.y
}

func (this *NodeModel) SetPosition(x, y float64) {
	this.x = x
	this.y = y
}

func (this *NodeModel) Connections() []ui.NodeModel {
	peers := this.connectionInfo.ConnectedPeers()
	result := make([]ui.NodeModel, len(peers))
	for i, peer := range peers {
		result[i] = this.nodeModels.NodeModelForId(peer.Id())
	}

	return result
}
