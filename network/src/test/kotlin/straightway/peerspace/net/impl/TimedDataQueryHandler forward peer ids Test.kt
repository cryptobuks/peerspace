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
import com.nhaarman.mockito_kotlin.mock
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import straightway.peerspace.data.Chunk
import straightway.peerspace.data.Id
import straightway.peerspace.data.Key
import straightway.peerspace.net.Configuration
import straightway.peerspace.net.DataQueryHandler
import straightway.peerspace.net.PushRequest
import straightway.peerspace.net.QueryRequest
import straightway.testing.bdd.Given
import straightway.testing.flow.Empty
import straightway.testing.flow.Equal
import straightway.testing.flow.Values
import straightway.testing.flow.expect
import straightway.testing.flow.is_
import straightway.testing.flow.to_
import straightway.units.get
import straightway.units.plus
import straightway.units.second
import straightway.units.toDuration
import straightway.units.year
import java.time.LocalDateTime

class `TimedDataQueryHandler forward peer ids Test` : KoinTestBase() {

    private companion object {
        val peerId = Id("peerId")
        val timedQueryingPeerId = Id("timedQueryingPeerId")
        val queriedChunkId = Id("queriedChunkId")
        val knownPeersIds = ids("1") + timedQueryingPeerId
        val timedQueryRequest = QueryRequest(timedQueryingPeerId, queriedChunkId, 1L..2L)
        val timedQueryResult = Chunk(Key(queriedChunkId, 1L), byteArrayOf(1, 2, 3))
        val timedResultPushRequest = PushRequest(peerId, timedQueryResult)
        val otherTimedQueryResult = Chunk(Key(queriedChunkId, 2L), byteArrayOf(1, 2, 3))
        val forwardedPeers = 0..0
    }

    private var currentTime = LocalDateTime.of(2001, 1, 1, 14, 30)
    private val test get() = Given {
        PeerTestEnvironment(
                peerId,
                knownPeersIds = knownPeersIds,
                forwardStrategyFactory = {
                    mock {
                        on { getQueryForwardPeerIdsFor(any(), any()) }
                                .thenReturn(knownPeersIds.slice(forwardedPeers))
                    }
                },
                configurationFactory = {
                    Configuration(
                        untimedDataQueryTimeout = 10[second],
                        timedDataQueryTimeout = 10[second])
                },
                dataQueryHandlerFactory = { TimedDataQueryHandler() },
                timeProviderFactory = {
                    mock {
                        on { currentTime }.thenAnswer { currentTime }
                    }
                })
    }

    @BeforeEach
    fun setup() {
        currentTime = LocalDateTime.of(2001, 1, 1, 14, 30)
    }

    @Test
    fun `not matching chunk is not forwarded`() =
            test while_ {
                get<DataQueryHandler>().handle(timedQueryRequest)
            } when_ {
                get<DataQueryHandler>().getForwardPeerIdsFor(Key(Id("otherId")))
            } then {
                expect(it.result is_ Empty)
            }

    @Test
    fun `timed result being received immediately is forwarded`() =
            test while_ {
                get<DataQueryHandler>().handle(timedQueryRequest)
            } when_ {
                get<DataQueryHandler>().getForwardPeerIdsFor(timedResultPushRequest.chunk.key)
            } then {
                expect(it.result is_ Equal to_ Values(timedQueryingPeerId))
            }

    @Test
    fun `timed result not forwarded after timeout expired`() =
            test andGiven {
                it.delayForwardingOfUntimedQueries()
            } while_ {
                get<DataQueryHandler>().handle(timedQueryRequest)
                currentTime += (get<Configuration>().timedDataQueryTimeout +
                        1[second]).toDuration()
            } when_ {
                get<DataQueryHandler>().getForwardPeerIdsFor(timedResultPushRequest.chunk.key)
            } then {
                expect(it.result is_ Empty)
            }

    @Test
    fun `multiple timed result being received are all forwarded`() =
            test while_ {
                get<DataQueryHandler>().handle(timedQueryRequest)
            } when_ {
                get<DataQueryHandler>().getForwardPeerIdsFor(timedResultPushRequest.chunk.key) +
                get<DataQueryHandler>().getForwardPeerIdsFor(otherTimedQueryResult.key)
            } then {
                expect(it.result is_ Equal to_ Values(timedQueryingPeerId, timedQueryingPeerId))
            }

    @Test
    fun `timed result being received twice is not forwarded after again`() =
            test while_ {
                get<DataQueryHandler>().handle(timedQueryRequest)
                get<DataQueryHandler>().notifyChunkForwarded(timedResultPushRequest.chunk.key)
            } when_ {
                get<DataQueryHandler>().getForwardPeerIdsFor(timedResultPushRequest.chunk.key)
            } then {
                expect(it.result is_ Empty)
            }

    private fun PeerTestEnvironment.delayForwardingOfUntimedQueries() = copy(
            configurationFactory = {
                this@delayForwardingOfUntimedQueries
                        .get<Configuration>().copy(untimedDataQueryTimeout = 1[year])
            })
}