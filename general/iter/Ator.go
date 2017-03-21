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

package iter

import (
	"reflect"

	"github.com/straightway/straightway/general/loop"
)

type Ator func() (interface{}, bool)

func (this Ator) Do(body func(interface{}) loop.Control) {
	iterFunc := (func() (interface{}, bool))(this)
	for i, isFound := iterFunc(); isFound; i, isFound = iterFunc() {
		if body(i) == loop.Break {
			break
		}
	}
}

func OnSlice(slice interface{}) Ator {
	index := 0
	itemsSlice := reflect.ValueOf(slice)
	return Ator(func() (result interface{}, isFound bool) {
		if index < itemsSlice.Len() {
			isFound = true
			result = itemsSlice.Index(index).Interface()
			index++
		}
		return
	})
}
