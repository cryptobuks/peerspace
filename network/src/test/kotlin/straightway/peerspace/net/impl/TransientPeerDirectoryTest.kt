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
package straightway.peerspace.net.impl

import org.junit.jupiter.api.Test
import straightway.peerspace.data.Id
import straightway.testing.bdd.Given
import straightway.testing.flow.Empty
import straightway.testing.flow.Equal
import straightway.testing.flow.Values
import straightway.testing.flow.expect
import straightway.testing.flow.is_
import straightway.testing.flow.to_

class TransientPeerDirectoryTest {

    private val test get() = Given {
        object {
            val sut = TransientPeerDirectory()
            val id = Id("id")
        }
    }

    @Test
    fun `initially the peer directory is empty`() =
            test when_ { sut.allKnownPeersIds } then {
                expect(it.result is_ Empty)
            }

    @Test
    fun `adding to the empty peer directory yields a one element directory`() =
            test when_ { sut add id } then {
                expect(sut.allKnownPeersIds is_ Equal to_ Values(id))
            }

    @Test
    fun `adding the same id twice ignores the second add`() =
            test while_ {
                sut add id
            } when_ {
                sut add id
            } then {
                expect(sut.allKnownPeersIds is_ Equal to_ Values(id))
            }
}