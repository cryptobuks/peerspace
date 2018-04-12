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

import straightway.peerspace.data.Chunk
import straightway.peerspace.data.Id
import straightway.peerspace.net.Configuration

data class PeerTestEnvironmentImpl(
        override val peerId: Id,
        override val knownPeersIds: List<Id> = listOf(),
        override val unknownPeerIds: List<Id> = listOf(),
        override var configuration: Configuration = Configuration(),
        override val localChunks: List<Chunk> = listOf()
) : PeerTestEnvironment {
    override val knownPeers = knownPeersIds.map { createPeerMock(it) }
    override val unknownPeers = knownPeersIds.map { createPeerMock(it) }
    override val network = createNetworkMock { knownPeers + unknownPeers }
    override val peerDirectory = createPeerDirectory { knownPeers }
    override var knownPeerQueryChooser = createChooser { knownPeersIds }
    override var knownPeerAnswerChooser = createChooser { knownPeersIds }
    override val chunkDataStore = createChunkDataStore { localChunks }
    override val sut by lazy {
        createPeerImpl(
                peerId,
                peerDirectory = peerDirectory,
                network = network,
                configuration = configuration,
                dataChunkStore = chunkDataStore,
                knownPeerQueryChooser = knownPeerQueryChooser,
                knownPeerAnswerChooser = knownPeerAnswerChooser)
    }

    override fun getPeer(id: Id) =
            knownPeers.find { it.id == id } ?: unknownPeers.find { it.id == id }!!
}
