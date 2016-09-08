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

package gui

import (
	"time"

	"github.com/andlabs/ui"
	"github.com/straightway/straightway/general/gui"
	sui "github.com/straightway/straightway/simc/ui"
)

type MainWindow struct {
	*ui.Window
	controller            sui.Controller
	startButton           *ui.Button
	stopButton            *ui.Button
	pauseButton           *ui.Button
	simulationTimeDisplay *gui.VCenteredLabel
}

func NewMainWindow(controller sui.Controller) *MainWindow {
	mainWindow := &MainWindow{controller: controller}
	mainWindow.init()
	return mainWindow
}

func (this *MainWindow) SetStartEnabled(enabled bool) {
	setEnabled(this.startButton, enabled)
}

func (this *MainWindow) SetStopEnabled(enabled bool) {
	setEnabled(this.stopButton, enabled)
}

func (this *MainWindow) SetPauseEnabled(enabled bool) {
	setEnabled(this.pauseButton, enabled)
}

func (this *MainWindow) SetSimulationTime(time time.Time) {
	this.simulationTimeDisplay.SetText(time.String())
}

// Event handlers

func (this *MainWindow) onStartClicked(*ui.Button) {
	this.controller.Start()
}

func (this *MainWindow) onStopClicked(*ui.Button) {
	this.controller.Stop()
}

func (this *MainWindow) onPauseClicked(*ui.Button) {
	this.controller.Pause()
}

// Private

func (this *MainWindow) init() {
	this.Window = ui.NewWindow("Straightway Simulation", 200, 100, false)

	mainLayout := ui.NewVerticalBox()

	commandBar := ui.NewHorizontalBox()
	mainLayout.Append(commandBar, false)

	this.stopButton = ui.NewButton("#")
	this.stopButton.OnClicked(this.onStopClicked)
	commandBar.Append(this.stopButton, false)

	this.pauseButton = ui.NewButton("||")
	this.pauseButton.OnClicked(this.onPauseClicked)
	commandBar.Append(this.pauseButton, false)

	this.startButton = ui.NewButton(">")
	this.startButton.OnClicked(this.onStartClicked)
	commandBar.Append(this.startButton, false)

	stretcher := ui.NewVerticalBox()
	commandBar.Append(stretcher, true)

	this.simulationTimeDisplay = gui.NewVCenteredLabel("01.01.0000 00:00:00.000")
	commandBar.Append(this.simulationTimeDisplay, false)

	mainArea := ui.NewHorizontalBox()
	mainLayout.Append(mainArea, true)

	this.SetChild(mainLayout)
	this.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})

	this.controller.SetUi(this)
}

func setEnabled(control ui.Control, enabled bool) {
	if enabled {
		control.Enable()
	} else {
		control.Disable()
	}
}
