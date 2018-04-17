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

import straightway.peerspace.net.Configuration
import straightway.peerspace.net.DataChunkStore
import straightway.peerspace.net.DataQueryHandler
import straightway.peerspace.net.ForwardStrategy
import straightway.peerspace.net.Infrastructure
import straightway.peerspace.net.Network
import straightway.peerspace.net.PeerDirectory
import straightway.random.Chooser
import straightway.utils.TimeProvider

/**
 * Infrastructure containing components required for implementing
 * the peer functionality.
 */
data class InfrastructureImpl(
        override val dataChunkStore: DataChunkStore,
        override val peerDirectory: PeerDirectory,
        override val network: Network,
        override val configuration: Configuration,
        override val knownPeerQueryChooser: Chooser,
        override val knownPeerAnswerChooser: Chooser,
        override val forwardStrategy: ForwardStrategy,
        override val timeProvider: TimeProvider,
        override val dataQueryHandler: DataQueryHandler
) : Infrastructure {
    init {
        dataQueryHandler.infrastructure = this
    }
}