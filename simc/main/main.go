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
	goui "github.com/andlabs/ui"
	"github.com/straightway/straightway/simc"
	"github.com/straightway/straightway/simc/gui"
	"github.com/straightway/straightway/simc/uic"
)

func main() {
	err := goui.Main(func() {
		eventScheduler := &simc.EventScheduler{}
		controller := &uic.Controller{
			SimulationController: eventScheduler}
		mainWindow := gui.NewMainWindow(controller)
		mainWindow.Show()
	})
	if err != nil {
		panic(err)
	}
}
