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

import straightway.units.Time
import straightway.units.UnitNumber
import straightway.units.get
import straightway.units.minute
import straightway.units.second

/**
 * Global configuration for a peerspace peer node.
 */
@Suppress("MagicNumber")
data class Configuration(
        val maxPeersToQueryForKnownPeers: Int = 10,
        val maxKnownPeersAnswers: Int = 20,
        val untimedDataQueryTimeout: UnitNumber<Time> = 30[second],
        val timedDataQueryTimeout: UnitNumber<Time> = 5[minute])