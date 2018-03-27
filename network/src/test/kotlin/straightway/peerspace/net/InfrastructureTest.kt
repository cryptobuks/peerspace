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

package straightway.peerspace.net

import org.junit.jupiter.api.Test
import com.nhaarman.mockito_kotlin.mock
import straightway.testing.flow.Same
import straightway.testing.flow.as_
import straightway.testing.flow.expect
import straightway.testing.flow.is_

class InfrastructureTest {

    @Test
    fun `setting and getting network`() {
        val network = mock<Network>()
        val sut = Infrastructure {
            this.network = network
        }
        expect(sut.network is_ Same as_ network)
    }

    @Test
    fun `setting and getting peerFactory`() {
        val peerFactory = mock<Factory<Peer>>()
        val sut = Infrastructure {
            this.peerStubFactory = peerFactory
        }
        expect(sut.peerStubFactory is_ Same as_ peerFactory)
    }

    @Test
    fun `setting and getting channelFactory`() {
        val channelFactory = mock<Factory<Channel>>()
        val sut = Infrastructure {
            this.channelFactory = channelFactory
        }
        expect(sut.channelFactory is_ Same as_ channelFactory)
    }
}