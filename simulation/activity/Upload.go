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
	"time"

	"github.com/straightway/straightway/data"
	"github.com/straightway/straightway/general"
	"github.com/straightway/straightway/peer"
	"github.com/straightway/straightway/simulation"
	"github.com/straightway/straightway/simulation/randvar"
)

type Upload struct {
	User         *simulation.User
	Delay        randvar.Duration
	DataSize     randvar.Float64
	IdGenerator  general.IdGenerator
	ChunkCreator simulation.ChunkCreator
	Audience     []simulation.DataConsumer
}

func (this *Upload) ScheduleUntil(maxTime time.Time) {
	scheduler := this.User.Scheduler
	now := scheduler.Time()
	nextActionTime := now.Add(this.Delay.NextSample())
	for !maxTime.Before(nextActionTime) {
		scheduler.ScheduleAbsolute(nextActionTime, this.doPush)
		nextActionTime = nextActionTime.Add(this.Delay.NextSample())
	}
}

// Private

func (this *Upload) doPush() {
	newDataChunk := this.ChunkCreator.CreateChunk(
		data.Key{Id: data.Id(this.IdGenerator.NextId())},
		uint64(this.DataSize.NextSample()))
	this.User.Node.Push(newDataChunk, this.User)
	for _, consumer := range this.Audience {
		consumer.AttractTo(peer.Query{Id: newDataChunk.Key.Id})
	}
}
