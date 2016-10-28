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

package main

import (
	"github.com/andlabs/ui"

	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"

	"github.com/pkg/profile"

	ggui "github.com/straightway/straightway/general/gui"
	"github.com/straightway/straightway/simc"
	"github.com/straightway/straightway/simc/env"
	//simlog "github.com/straightway/straightway/simc/log"
	"github.com/straightway/straightway/simc/uic"
	"github.com/straightway/straightway/simc/uic/gui"
)

func main() {
	defer profile.Start().Stop()

	err := ui.Main(func() {
		scheduler := simc.NewEventScheduler()
		//actionLogger := simlog.NewActionHandler(simlog.DefaultBasicHandler)
		//simTimeLogHandler := simlog.NewSimulationTimeHandler(actionLogger, scheduler)
		//log.SetHandler(simTimeLogHandler)
		log.SetHandler(discard.New())
		toolkitAdapter := &ggui.ToolkitAdapter{}
		eventControllerAdapter := &uic.SimulationControllerAdapter{
			SimulationController: scheduler,
			ToolkitAdapter:       toolkitAdapter,
			TimeProvider:         scheduler,
			EnvironmentFactory:   func() *env.Environment { return env.New(scheduler, 1000) }}
		controller := uic.NewController(scheduler, eventControllerAdapter, eventControllerAdapter, toolkitAdapter)
		controller.MeasurementUpdateRatio = 10
		mainWindow := gui.NewMainWindow(controller, eventControllerAdapter)
		mainWindow.Show()
	})
	if err != nil {
		panic(err)
	}
}
