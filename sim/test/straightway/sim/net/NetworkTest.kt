/****************************************************************************
Copyright 2016 github.com/straightway

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
 ****************************************************************************/
package straightway.sim.net

import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import straightway.sim.core.*
import straightway.testing.*
import straightway.testing.flow.*
import straightway.units.*
import java.time.LocalDateTime

class NetworkTest : TestBase<NetworkTest.Environment>() {

    inner class Environment {
        val simulator = Simulator()
        val log = TimeLog(simulator)
        val network = Network(simulator, simulator, latency = 2[second])
        var sender = ClientMock("sender", log)
        var receiver = ClientMock("receiver", log)
        var message = createMessage()
    }

    @BeforeEach fun setup() {
        sut = Environment()
    }

    @Test
    fun send_triggersTransmissionOnChannels() =
        sut.run {
            network.send(sender, receiver, message)
            expect(log.entries _is equal _to listOf(
                "00:00:00: sender_upload: Transmit $message from sender_upload to receiver_download",
                "00:00:00: receiver_download: Transmit $message from sender_upload to receiver_download"))
        }

    @Test
    fun send_schedulesReceiveCall() =
        sut.run {
            sender.uploadChannel.receiveTime = LocalDateTime.of(0, 1, 1, 0, 3)
            network.send(sender, receiver, message)
            log.entries.clear()
            simulator.run()
            expect(log.entries _is equal _to listOf("00:03:02: Receive $message from sender to receiver"))
        }
}