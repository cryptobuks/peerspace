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

package storage

import (
	"github.com/straightway/straightway/data"
	"github.com/straightway/straightway/general"
	"github.com/straightway/straightway/peer"
)

type ChunkOrder func(a, b *data.Key) bool

type Raw interface {
	Query(peer.Query) []DataRecord
	Store(chunk *data.Chunk, priority float32)
	Delete(data.Key) int
	GetSizeOf(*data.Chunk) int
	GetFreeStorage() int
	GetLeastImportantData() general.Iterator
}
