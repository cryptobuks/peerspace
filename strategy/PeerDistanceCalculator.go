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

package strategy

import (
	"github.com/straightway/straightway/data"
	"github.com/straightway/straightway/general/id"
)

// Calculate distances between an ID holder and a data key or a query.
type PeerDistanceCalculator interface {
	Distance(id.Holder, data.Key) uint64
	Distances(id.Holder, data.Query) []uint64
}
