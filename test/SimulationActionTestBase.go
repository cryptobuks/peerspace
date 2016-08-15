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
	"time"

	"github.com/straightway/straightway/general/duration"
	"github.com/straightway/straightway/mocked"
	"github.com/straightway/straightway/simc"
	"github.com/stretchr/testify/suite"
)

type SimulationActionTestBase struct {
	suite.Suite
	scheduler   *simc.EventScheduler
	user        *simc.User
	node        *mocked.Node
	offlineTime time.Time
}

func (suite *SimulationActionTestBase) SetupTest() {
	suite.scheduler = &simc.EventScheduler{}
	suite.node = mocked.NewNode("1")
	suite.user = &simc.User{
		Scheduler: suite.scheduler,
		Node:      suite.node}
	suite.scheduler.Schedule(duration.Parse("1000h"), func() {
		suite.scheduler.Stop()
	})
	now := suite.scheduler.Time()
	suite.offlineTime = now.Add(onlineDuration)
}

func (suite *SimulationActionTestBase) TearDownTest() {
	suite.scheduler = nil
	suite.node = nil
	suite.user = nil
}
