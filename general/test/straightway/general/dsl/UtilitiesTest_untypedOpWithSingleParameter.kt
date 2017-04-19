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
package straightway.general.dsl

import org.junit.jupiter.api.Assertions
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Test

class UtilitiesTest_untypedOpWithSingleParameter {

    @Test fun returnsLambdaWithAnyParametersAndReturnType() {
        val result = untypedOp<Int> { a -> a * 3}
        Assertions.assertTrue(result is (Any) -> Any)
    }

    @Test fun callsPassedLambda() {
        var calls = 0;
        val result = untypedOp<Int> { a -> calls++; -a }
        assertEquals(0, calls)
        val callResult = result(5)
        assertEquals(1, calls)
        assertEquals(-5, callResult)
    }
}