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
import com.nhaarman.mockito_kotlin.never
import com.nhaarman.mockito_kotlin.verify
import org.junit.jupiter.api.Test
import straightway.peerspace.data.Chunk
import straightway.peerspace.data.Id
import straightway.peerspace.data.Key
import straightway.peerspace.net.PushRequest
import straightway.peerspace.net.QueryRequest
import straightway.testing.bdd.Given

class `TimedDataQueryHandler general Test` {

    private companion object {
        val peerId = Id("peerId")
        val receiverId = Id("receiverId")
        val chunkId = Id("chunkId")
        val chunkData = "ChunkData".toByteArray()
        val chunk = Chunk(Key(chunkId, 1L), chunkData)
        val queryRequest = QueryRequest(receiverId, chunkId, 1L..2L)
    }

    private val test get() = Given {
        object : PeerTestEnvironment by PeerTestEnvironmentImpl(
                peerId,
                knownPeersIds = listOf(receiverId),
                dataQueryHandler = TimedDataQueryHandler(peerId)
        ).fixed() {
            val receiver = getPeer(receiverId)
        }
    }

    @Test
    fun `query is forwarded to chunk data store`() =
            test when_ {
                dataQueryHandler.handle(queryRequest)
            } then {
                verify(dataChunkStore).query(queryRequest)
            }

    @Test
    fun `query for not existing data does not push back`() =
            test when_ {
                dataQueryHandler.handle(queryRequest)
            } then {
                verify(receiver, never()).push(any())
            }

    @Test
    fun `query hit returns result to sender`() =
            test while_ {
                dataChunkStore.store(chunk)
            } when_ {
                dataQueryHandler.handle(queryRequest)
            } then {
                verify(receiver).push(PushRequest(peerId, chunk))
            }
}