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
import com.nhaarman.mockito_kotlin.doAnswer
import com.nhaarman.mockito_kotlin.mock
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import straightway.peerspace.data.Chunk
import straightway.peerspace.data.Id
import straightway.peerspace.data.Key
import straightway.peerspace.net.Channel
import straightway.peerspace.net.Factory
import straightway.peerspace.net.Infrastructure
import straightway.peerspace.net.PushRequest
import straightway.testing.TestBase
import straightway.testing.flow.Equal
import straightway.testing.flow.expect
import straightway.testing.flow.is_
import straightway.testing.flow.to_

class PeerNetworkStubTest : TestBase<PeerNetworkStubTest.Environment>() {

    class Environment {

        val peerId = Id("id")

        val transmittedData = mutableListOf<Any>()
        val data = Chunk(Key(Id("Key")), arrayOf(1, 2, 3))

        var channelMockFactoryInvocations = 0

        private val channelFactoryMock = mock<Factory<Channel>> {
            on { create(any()) } doAnswer {
                ++channelMockFactoryInvocations
                channelMock
            }
        }

        private val channelMock = mock<Channel> {
            on { transmit(any()) } doAnswer { transmittedData.add(it.arguments[0]); null }
        }

        private val infrastructure = Infrastructure { channelFactory = channelFactoryMock }

        val sut = PeerNetworkStub(peerId, infrastructure)
    }

    @BeforeEach
    fun setup() {
        sut = Environment()
    }

    @Test
    fun `has specified id`() = sut.run {
        expect(sut.id is_ Equal to_ peerId)
    }

    @Test
    fun `push creates channel`() = sut.run {
        sut.push(PushRequest(data))
        expect(channelMockFactoryInvocations is_ Equal to_ 1)
    }

    @Test
    fun `push transmits request on channel`() = sut.run {
        val request = PushRequest(data)
        sut.push(request)
        expect(transmittedData is_ Equal to_ listOf(request))
    }
}