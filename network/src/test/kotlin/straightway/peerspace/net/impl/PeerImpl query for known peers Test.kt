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

import com.nhaarman.mockito_kotlin.any
import com.nhaarman.mockito_kotlin.argThat
import com.nhaarman.mockito_kotlin.verify
import org.junit.jupiter.api.Test
import straightway.peerspace.data.Id
import straightway.peerspace.data.Key
import straightway.peerspace.net.Administrative
import straightway.peerspace.net.QueryRequest
import straightway.testing.bdd.Given
import straightway.utils.deserializeTo

class `PeerImpl query for known peers Test` {

    private companion object {
        val peerId = Id("PeerId")
        val knownPeerId = Id("knownPeerId")
        val queryingPeerId = Id("QueryingPeerId")
        val knownPeersRequest = QueryRequest(queryingPeerId, Administrative.KnownPeers)
    }

    private val test get() = Given {
        object : PeerTestEnvironment by PeerTestEnvironmentImpl(
                peerId,
                knownPeersIds = listOf(knownPeerId, queryingPeerId)
        ) {
            val queryingPeer = getPeer(queryingPeerId)
        }
    }

    @Test
    fun `a query for known peers is answered immediately`() =
            test when_ {
                sut.query(knownPeersRequest)
            } then {
                verify(queryingPeer).push(any())
            }

    @Test
    fun `a query for known peers is answered with Administrative_KnownPeers_id`() =
            test when_ {
                sut.query(knownPeersRequest)
            } then {
                verify(queryingPeer).push(argThat {
                    chunk.key == Key(Administrative.KnownPeers.id)
                })
            }

    @Test
    fun `a query for known peers is answered with the list of known peers`() =
            test when_ {
                sut.query(knownPeersRequest)
            } then {
                verify(queryingPeer).push(argThat {
                    chunk.data.deserializeTo<List<Id>>() == knownPeersIds
                })
            }

    @Test
    fun `a query for known peers yields not more peers than specified`() =
            test while_ {
                configuration = configuration.copy(maxKnownPeersAnswers = 1)
                knownPeerAnswerChooser = createChooser { knownPeersIds.slice(0..0) }
            } when_ {
                sut.query(knownPeersRequest)
            } then {
                verify(queryingPeer).push(argThat {
                    chunk.data.deserializeTo<List<Id>>() == knownPeersIds.slice(0..0)
                })
            }
}