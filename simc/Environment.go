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
	"fmt"
	"hash/crc64"
	"math/rand"
	"time"

	"github.com/straightway/straightway/app"
	"github.com/straightway/straightway/data"
	"github.com/straightway/straightway/datac"
	"github.com/straightway/straightway/general/duration"
	"github.com/straightway/straightway/mocked"
	"github.com/straightway/straightway/peer"
	"github.com/straightway/straightway/peerc"
	"github.com/straightway/straightway/sim"
	"github.com/straightway/straightway/simc/activity"
	"github.com/straightway/straightway/simc/measure"
	"github.com/straightway/straightway/simc/randvar"
	"github.com/straightway/straightway/simc/ui"
	"github.com/straightway/straightway/strategy"
)

const (
	gb = 1024 * 1024 * 1024
)

type Environment struct {
	scheduler            sim.EventScheduler
	users                []*User
	uiNodeForId          map[string]ui.NodeModel
	uiNodes              []ui.NodeModel
	nextNodeId           uint
	randSource           rand.Source
	initialUser          *User
	queryDurationMeasure *measure.Discrete
	querySuccessMeasure  *measure.Discrete
}

func NewSimulationEnvironment(
	scheduler sim.EventScheduler,
	numberOfUsers int) *Environment {
	result := &Environment{
		scheduler:            scheduler,
		randSource:           rand.NewSource(12345),
		queryDurationMeasure: &measure.Discrete{},
		querySuccessMeasure:  &measure.Discrete{},
		uiNodeForId:          make(map[string]ui.NodeModel)}
	result.createSeedNode()
	for i := 0; i < numberOfUsers; i++ {
		result.addNewUser()
	}
	return result
}

func (this *Environment) QueryDurationMeasure() *measure.Discrete {
	return this.queryDurationMeasure
}

func (this *Environment) QuerySuccessMeasure() *measure.Discrete {
	return this.querySuccessMeasure
}

func (this *Environment) Audience() []sim.DataConsumer {
	result := make([]sim.DataConsumer, len(this.users))
	for i, u := range this.users {
		result[i] = u
	}
	return result
}

func (this *Environment) Nodes() []ui.NodeModel {
	return this.uiNodes
}

func (this *Environment) NodeModelForId(id string) ui.NodeModel {
	result := this.uiNodeForId[id]
	if result == nil {
		panic(fmt.Sprintf("Cannot get node model for %v", id))
	}
	return result
}

// Private

func (this *Environment) addNewUser() *User {
	newUser := this.createUser()
	this.users = append(this.users, newUser)
	return newUser
}

func (this *Environment) createSeedNode() {
	node, _, _ := this.createNode()
	node.PostConnectAction = func(node peer.Node, peer peer.Connector) {
		node.CloseConnectionWith(peer)
		peer.CloseConnectionWith(node)
	}

	this.initialUser = &User{
		SchedulerInstance:      this.scheduler,
		NodeInstance:           node,
		StartupDuration:        randvar.NewNormalDuration(this.randSource, time.Duration(0), time.Duration(0)),
		OnlineDuration:         randvar.NewNormalDuration(this.randSource, 2000000*time.Hour, time.Duration(0)),
		OnlineActivity:         mocked.NewSimulationUserActivity(),
		QuerySelectionSelector: rand.New(this.randSource)}
	this.initialUser.Activate()
}

func (this *Environment) createUser() *User {
	node, configuration, rawStorage := this.createNode()
	newUser := &User{
		SchedulerInstance:            this.scheduler,
		NodeInstance:                 node,
		StartupDuration:              this.newPositiveNormalDuration(duration.Parse("8h"), duration.Parse("2h")),
		OnlineDuration:               this.newPositiveNormalDuration(duration.Parse("2h"), duration.Parse("2h")),
		QueryDurationSampleCollector: this.queryDurationMeasure,
		QuerySuccessSampleCollector:  this.querySuccessMeasure,
		QueryWaitingTimeout:          duration.Parse("5m"),
		QuerySelectionSelector:       rand.New(this.randSource)}
	newUser.OnlineActivity = this.createActivity(newUser, configuration, rawStorage)
	newUser.Activate()
	return newUser
}

func (this *Environment) createNode() (*peerc.Node, *app.Configuration, *RawStorage) {
	this.nextNodeId++
	nodeId := fmt.Sprintf("%v", this.nextNodeId)
	configuration := app.DefaultConfiguration()
	peerDistanceRelated := &strategy.PeerDistanceRelated{
		LocalPeerId: nodeId,
		Timer:       this.scheduler,
		Hasher:      crc64.New(crc64.MakeTable(crc64.ECMA))}
	dataStorage, rawStorage := this.createDataStorage(peerDistanceRelated)
	stateStorage := this.createStateStorage(rawStorage)
	newNode := &peerc.Node{
		Identifier:           nodeId,
		StateStorage:         stateStorage,
		DataStorage:          dataStorage,
		AnnouncementStrategy: this.createAnnouncementStrategy(configuration, stateStorage),
		Timer:                this.scheduler,
		Configuration:        configuration}
	newNode.DataStrategy = this.createDataStrategy(configuration, peerDistanceRelated, newNode)
	newNode.ConnectionStrategy = this.createConnecionStrategy(configuration, newNode)
	newNode.QueryStrategy = this.createQueryStrategy(configuration, peerDistanceRelated, newNode)
	nodeModel := NewNodeModel(nodeId, this, newNode)
	this.uiNodeForId[nodeId] = nodeModel
	this.uiNodes = append(this.uiNodes, nodeModel)
	return newNode, configuration, rawStorage
}

func (this *Environment) createStateStorage(rawStorage data.RawStorage) peer.StateStorage {
	stateStorage := &StateStorage{}
	if 0 < len(this.users) {
		networkAccessedNode := &NetworkPeerConnector{
			Wrapped:        this.initialUser.Node(),
			EventScheduler: this.scheduler,
			RawStorage:     rawStorage,
			Latency:        duration.Parse("50ms"),
			Bandwidth:      1024 * 1024}
		stateStorage.AddKnownPeer(networkAccessedNode)
	}

	return stateStorage
}

func (this *Environment) createDataStorage(
	priorityGenerator data.PriorityGenerator) (dataStorage data.Storage, rawStorage *RawStorage) {
	rawStorage = &RawStorage{
		FreeStorageValue: 2 * gb,
		Timer:            this.scheduler}
	dataStorage = &datac.Storage{
		PriorityGenerator: priorityGenerator,
		RawStorage:        rawStorage}
	return
}

func (this *Environment) createDataStrategy(
	configuration *app.Configuration,
	peerDistanceCalculator strategy.PeerDistanceCalculator,
	connectionInfoProvider strategy.ConnectionInfoProvider) peer.DataStrategy {
	return &strategy.Data{
		Configuration:          configuration,
		PeerDistanceCalculator: peerDistanceCalculator,
		ConnectionInfoProvider: connectionInfoProvider}
}

func (this *Environment) createAnnouncementStrategy(
	configuration *app.Configuration,
	stateStorage peer.StateStorage) peer.AnnouncementStrategy {
	return &strategy.Announcement{
		Configuration: configuration,
		RandomSource:  this.randSource,
		StateStorage:  stateStorage}
}

func (this *Environment) createConnecionStrategy(
	configuration *app.Configuration,
	connectionInfoProvider strategy.ConnectionInfoProvider) peer.ConnectionStrategy {
	return &strategy.Connection{
		Configuration:          configuration,
		ConnectionInfoProvider: connectionInfoProvider,
		RandSource:             this.randSource}
}

func (this *Environment) createQueryStrategy(
	configuration *app.Configuration,
	peerDistanceCalculator strategy.PeerDistanceCalculator,
	connectionInfoProvider strategy.ConnectionInfoProvider) peer.QueryStrategy {
	return &strategy.Query{
		ConnectionInfoProvider: connectionInfoProvider,
		PeerDistanceCalculator: peerDistanceCalculator,
		Configuration:          configuration}
}

func (this *Environment) createActivity(
	user *User,
	configuration *app.Configuration,
	chunkCreator sim.ChunkCreator) sim.UserActivity {
	return activity.NewCombined(
		this.createUploadActivity(user, configuration, chunkCreator),
		this.createQueryActivity(user))
}

func (this *Environment) createUploadActivity(
	user *User,
	configuration *app.Configuration,
	chunkCreator sim.ChunkCreator) sim.UserActivity {
	return &activity.Upload{
		User:               user,
		Configuration:      configuration,
		Delay:              this.newPositiveNormalDuration(duration.Parse("15m"), duration.Parse("30m")),
		DataSize:           randvar.NewNormalFloat64(this.randSource, 32000, 32000),
		IdGenerator:        &IdGenerator{RandSource: this.randSource},
		ChunkCreator:       chunkCreator,
		AudienceProvider:   this,
		AttractionRatio:    randvar.NewNormalFloat64(this.randSource, 0.3, 0.1),
		AudiencePermutator: rand.New(this.randSource)}
}

func (this *Environment) createQueryActivity(
	user *User) sim.UserActivity {
	return &activity.Query{
		Scheduler:          user.Scheduler(),
		User:               user,
		QueryPauseDuration: this.newPositiveNormalDuration(duration.Parse("5m"), duration.Parse("10m"))}
}

func (this *Environment) newPositiveNormalDuration(mean, stdDev time.Duration) *randvar.PositiveDuration {
	return randvar.NewPositiveDuration(randvar.NewNormalDuration(this.randSource, mean, stdDev))
}
