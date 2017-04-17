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
package straightway.test.flow

import org.junit.jupiter.api.Assertions.assertTrue
import org.junit.jupiter.api.Assertions.fail
import org.opentest4j.AssertionFailedError
import straightway.general.dsl.Expr

/**
 * Expect the given condition to be true.
 */
fun expect(condition: Expr) =
    try {
        assertTrue(condition() as Boolean) { ExpressionVisualizer(condition).string }
    } catch (e: AssertionFailedError) {
        throw e
    } catch (e: AssertionError) {
        fail("${ExpressionVisualizer(condition).string} (${e.message})")
    } catch (e: Throwable) {
        fail("${ExpressionVisualizer(condition).string}: $e")
    }

/**
 * Expect the given condition to be true.
 */
fun expect(condition: Boolean) =
        assertTrue(condition)

/**
 * Expect the given condition to be true.
 */
fun expect(condition: Boolean, lazyExplanation: () -> String) =
        assertTrue(condition, lazyExplanation)
