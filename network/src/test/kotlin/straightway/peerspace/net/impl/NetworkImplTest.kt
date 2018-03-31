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
import com.nhaarman.mockito_kotlin.verify
import org.junit.jupiter.api.Test
import straightway.peerspace.data.Id
import straightway.peerspace.net.Factory
import straightway.peerspace.net.PushTarget
import straightway.peerspace.net.QuerySource
import straightway.testing.bdd.Given

class NetworkImplTest {

    private companion object {
        val receiverId = Id("receiver")
    }

    private val test get() =
            Given {
                object {
                    val pushTargetStubFactory = mock<Factory<PushTarget>> {
                        on { create(any()) }.thenReturn(mock())
                    }
                    val querySourceStubStubFactory = mock<Factory<QuerySource>> {
                        on { create(any()) }.thenReturn(mock())
                    }
                    val sut = NetworkImpl(pushTargetStubFactory, querySourceStubStubFactory)
                }
            }

    @Test
    fun `getPushTarget callsPeerFactory`() =
            test when_ { sut.getPushTarget(receiverId) } then {
                verify(pushTargetStubFactory).create(receiverId)
            }

    @Test
    fun `getQuerySource callsPeerFactory`() =
            test when_ { sut.getQuerySource(receiverId) } then {
                verify(querySourceStubStubFactory).create(receiverId)
            }
}