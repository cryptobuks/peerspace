/*
 * Copyright 2016 github.com/straightway
 *
 *  Licensed under the Apache License, Version 2.0 (the &quot;License&quot;);
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an &quot;AS IS&quot; BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package straightway.data

import org.junit.jupiter.api.Test
import straightway.testing.flow.equal
import straightway.testing.flow.expect
import straightway.testing.flow.is_
import straightway.testing.flow.to_

class ChunkTest {

    @Test
    fun `key is as specified in construction`() =
            expect(Chunk(Key("1234"), "data").key is_ equal to_ Key("1234"))

    @Test
    fun `data is as specified in construction`() =
            expect(Chunk(Key("1234"), "data").data is_ equal to_ "data")

    @Test
    fun `Chunk is serializable`() {
        val sut = Chunk(Key("1234"), "data")
        val serialized = sut.serializeToByteArray()
        val deserialized = serialized.deserializeTo<Chunk>()
        expect(deserialized is_ equal to_ sut)
    }
}