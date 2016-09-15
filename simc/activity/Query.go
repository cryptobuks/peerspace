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

package activity

import (
	"log"
	"time"

	"github.com/straightway/straightway/sim"
	"github.com/straightway/straightway/sim/randvar"
)

type Query struct {
	Scheduler          sim.EventScheduler
	User               sim.User
	QueryPauseDuration randvar.Duration
}

func (this *Query) ScheduleUntil(maxTime time.Time) {
	currSimTime := this.Scheduler.Time()
	for {
		currSimTime = currSimTime.Add(this.QueryPauseDuration.NextSample())
		if maxTime.Before(currSimTime) {
			return
		}

		this.Scheduler.ScheduleAbsolute(currSimTime, this.doQuery)
	}
}

// Private

func (this *Query) doQuery() {
	query, isQueryFound := this.User.PopAttractiveQuery()
	if isQueryFound {
		log.Printf("%v: %v queries %v", this.Scheduler.Time(), this.User.Id(), query)
		this.User.Node().Query(query, this.User)
	} else {
		log.Printf("%v: %v has no more queries", this.Scheduler.Time(), this.User.Id())
	}
}