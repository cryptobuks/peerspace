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

package mocked

import "github.com/stretchr/testify/mock"

type SimulationRandVarIntner struct {
	Base
	values       []int
	currentIndex int
}

func NewSimulationRandVarIntner(firstValue int, nextValues ...int) *SimulationRandVarIntner {
	result := &SimulationRandVarIntner{}
	result.SetValues(firstValue, nextValues...)
	result.On("Intn", mock.Anything).Return()
	return result
}

func (m *SimulationRandVarIntner) SetValues(firstValue int, nextValues ...int) {
	m.values = append([]int{firstValue}, nextValues...)
	m.currentIndex = 0
}

func (m *SimulationRandVarIntner) Intn(n int) int {
	m.Called(n)
	nextValue := m.values[m.currentIndex]
	m.currentIndex = (m.currentIndex + 1) % len(m.values)
	return nextValue
}
