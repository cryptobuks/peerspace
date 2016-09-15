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

	"github.com/straightway/straightway/mocked"
	"github.com/straightway/straightway/simc/activity"
	"github.com/stretchr/testify/suite"
)

type SimulationCombinedActivity_Test struct {
	suite.Suite
}

func TestSimulationCombinedActivity(t *testing.T) {
	suite.Run(t, new(SimulationCombinedActivity_Test))
}

// Test

func (suite *SimulationCombinedActivity_Test) Test_ScheduleUntil_ForwardsToAllChildren() {
	child1 := mocked.NewSimulationUserActivity()
	child2 := mocked.NewSimulationUserActivity()
	sut := activity.NewCombined(child1, child2)
	t := time.Date(2000, 5, 7, 11, 37, 13, 7, time.UTC)
	sut.ScheduleUntil(t)
	child1.AssertCalledOnce(suite.T(), "ScheduleUntil", t)
	child2.AssertCalledOnce(suite.T(), "ScheduleUntil", t)
}