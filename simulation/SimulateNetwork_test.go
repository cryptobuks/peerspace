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
	"math/rand"
	"testing"

	"github.com/straightway/straightway/general"
	"github.com/straightway/straightway/peer"
	"github.com/straightway/straightway/simulation/randvar"
	"github.com/straightway/straightway/storage"
	"github.com/straightway/straightway/strategy"
)

func TestSimulatedNetwork(t *testing.T) {
	env := newSimulationEnvironment()
	numberOfUsers := 2
	for i := 0; i < numberOfUsers; i++ {
		env.addNewUser()
	}

	env.eventScheduler.Schedule(general.ParseDuration("24h"), func() { env.eventScheduler.Stop() })
	env.eventScheduler.Run()
}

const (
	gb = 1024 * 1024 * 1024
)

type simulationEnvironment struct {
	randSource     rand.Source
	eventScheduler EventScheduler
	users          []*User
}

func newSimulationEnvironment() *simulationEnvironment {
	return &simulationEnvironment{
		randSource: rand.NewSource(12345),
	}
}

func (this *simulationEnvironment) addNewUser() *User {
	newUser := this.createUser()
	this.users = append(this.users, newUser)
	return newUser
}

func (this *simulationEnvironment) createUser() *User {
	newUser := &User{
		Scheduler:       &this.eventScheduler,
		Node:            this.createNode(),
		StartupDuration: randvar.NewNormalDuration(this.randSource, general.ParseDuration("8h"), general.ParseDuration("2h")),
		OnlineDuration:  randvar.NewNormalDuration(this.randSource, general.ParseDuration("2h"), general.ParseDuration("2h")),
		OnlineAction:    func(*User) {},
		ActionDuration:  randvar.NewNormalDuration(this.randSource, general.ParseDuration("10m"), general.ParseDuration("30m"))}
	newUser.Activate()
	return newUser
}

func (this *simulationEnvironment) createNode() peer.Node {
	configuration := peer.DefaultConfiguration()
	peerDistanceRelated := &strategy.PeerDistanceRelated{}
	newNode := &peer.NodeImpl{
		StateStorage:  this.createStateStorage(),
		DataStorage:   this.createDataStorage(peerDistanceRelated),
		DataStrategy:  this.createDataStrategy(configuration, peerDistanceRelated),
		Timer:         &this.eventScheduler,
		Configuration: configuration}
	newNode.ConnectionStrategy = this.createConnecionStrategy(configuration, newNode)
	return newNode
}

func (this *simulationEnvironment) createStateStorage() peer.StateStorage {
	stateStorage := &StateStorage{}
	if 0 < len(this.users) {
		stateStorage.AddKnownPeer(this.users[0].Node)
	}

	return stateStorage
}

func (this *simulationEnvironment) createDataStorage(priorityGenerator storage.PriorityGenerator) peer.DataStorage {
	rawStorage := &RawStorage{
		FreeStorageValue: 2 * gb,
		Timer:            &this.eventScheduler}
	return &storage.DataImpl{
		PriorityGenerator: priorityGenerator,
		RawStorage:        rawStorage}
}

func (this *simulationEnvironment) createDataStrategy(
	configuration *peer.Configuration,
	peerDistanceCalculator strategy.PeerDistanceCalculator) peer.DataStrategy {
	return &strategy.Data{
		Configuration:          configuration,
		PeerDistanceCalculator: peerDistanceCalculator}
}

func (this *simulationEnvironment) createConnecionStrategy(
	configuration *peer.Configuration,
	connectionInfoProvider strategy.ConnectionInfoProvider) peer.ConnectionStrategy {
	return &strategy.Connection{
		Configuration:          configuration,
		ConnectionInfoProvider: connectionInfoProvider,
		RandSource:             this.randSource}
}
