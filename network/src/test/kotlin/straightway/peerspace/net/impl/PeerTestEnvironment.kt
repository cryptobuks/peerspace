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
import straightway.peerspace.data.Key
import straightway.peerspace.net.Configuration
import straightway.peerspace.net.DataChunkStore
import straightway.peerspace.net.DataPushForwarder
import straightway.peerspace.net.DataQueryHandler
import straightway.peerspace.net.ForwardStrategy
import straightway.peerspace.net.KnownPeersProvider
import straightway.peerspace.net.Peer
import straightway.peerspace.net.TransmissionResultListener
import straightway.random.Chooser
import straightway.utils.TimeProvider

/**
 * Test environment for testing the PeerImpl class.
 */
@Suppress("ComplexInterface")
interface PeerTestEnvironment {
    val peerId: Id
    val knownPeersIds: List<Id>
    val unknownPeerIds: List<Id>
    val configuration: Configuration
    val localChunks: List<Chunk>
    val knownPeers: MutableList<Peer>
    val unknownPeers: List<Peer>
    val knownPeerQueryChooser: Chooser
    val knownPeerAnswerChooser: Chooser
    val timeProvider: TimeProvider
    val peer: Peer
    val forwardStrategy: ForwardStrategy
    val dataQueryHandler: DataQueryHandler
    val dataPushForwarder: DataPushForwarder
    val knownPeersProvider: KnownPeersProvider
    val dataChunkStore: DataChunkStore
    val pushTransmissionResultListeners: MutableMap<Pair<Id, Key>, TransmissionResultListener>
    fun setPeerPushSuccess(id: Id, success: Boolean)
}

fun PeerTestEnvironment.getPeer(id: Id) =
        knownPeers.find { it.id == id } ?: unknownPeers.find { it.id == id }!!
