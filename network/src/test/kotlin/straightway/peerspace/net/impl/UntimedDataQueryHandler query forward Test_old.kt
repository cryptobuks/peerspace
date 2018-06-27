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
import com.nhaarman.mockito_kotlin.eq
import com.nhaarman.mockito_kotlin.mock
import com.nhaarman.mockito_kotlin.never
import com.nhaarman.mockito_kotlin.verify
import org.junit.jupiter.api.Test
import straightway.peerspace.data.Chunk
import straightway.peerspace.data.Id
import straightway.peerspace.data.Key
import straightway.peerspace.net.DataChunkStore
import straightway.peerspace.net.DataQueryHandler
import straightway.peerspace.net.ForwardState
import straightway.peerspace.net.ForwardStrategy
import straightway.peerspace.net.QueryRequest
import straightway.testing.bdd.Given

class `UntimedDataQueryHandler query forward Test_old` : KoinTestBase() {

    companion object {
        val peerId = Id("peerId")
        val queryingPeerId = Id("queryingPeerId")
        val queriedChunkId = Id("queriedChunkId")
        val untimedQueriedChunk = Chunk(Key(queriedChunkId), byteArrayOf(1, 2, 3))
        val knownPeersIds = ids("1", "2", "3") + queryingPeerId
        val receivedUntimedQueryRequest = QueryRequest(queryingPeerId, queriedChunkId)
        val forwardedUntimedQueryRequest = QueryRequest(peerId, queriedChunkId)
        val forwardedPeers = 1..2
    }

    private val test get() = Given {
        PeerTestEnvironment(
                peerId,
                knownPeersIds = knownPeersIds,
                forwardStrategyFactory = {
                    mock {
                        on {
                            getQueryForwardPeerIdsFor(any(), any())
                        }.thenReturn(knownPeersIds.slice(forwardedPeers))
                    }
                },
                dataQueryHandlerFactory = { UntimedDataQueryHandler() },
                queryForwarderFactory = { QueryForwarder() },
                queryForwardTrackerFactory = {
                    ForwardStateTrackerImpl(get("queryForwarder"))
                },
                pushForwarderFactory = { PushForwarder() })
    }

    @Test
    fun `query request is forwarded according to forward strategy`() =
            test when_ {
                get<DataQueryHandler>().handle(receivedUntimedQueryRequest)
            } then {
                verify(get<ForwardStrategy>()).getQueryForwardPeerIdsFor(
                        receivedUntimedQueryRequest, ForwardState())
            }

    @Test
    fun `query request is forwarded to peers returned by forward strategy`() =
            test when_ {
                get<DataQueryHandler>().handle(receivedUntimedQueryRequest)
            } then {
                forwardedPeers.forEach {
                    verify(knownPeers[it]).query(eq(forwardedUntimedQueryRequest), any())
                }
            }

    @Test
    fun `query is not forwarded if local result found`() =
            test while_ {
                get<DataChunkStore>().store(untimedQueriedChunk)
            } when_ {
                get<DataQueryHandler>().handle(receivedUntimedQueryRequest)
            } then {
                forwardedPeers.forEach {
                    verify(knownPeers[it], never()).query(any(), any())
                }
            }
}