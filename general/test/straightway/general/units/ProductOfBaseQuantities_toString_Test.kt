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
package straightway.general.units

import org.junit.jupiter.api.Assertions
import org.junit.jupiter.api.Test

class ProductOfBaseQuantities_toString_Test : ProductOfBaseQuantities_StringRepBaseTest({ toString() }) {

    @Test fun only_linear_Mass() =
        Assertions.assertEquals("g", ProductOfBaseQuantities(NoN, NoI, NoL, NoJ, linear(gramm), NoTheta, NoT).toString())

    @Test fun withLocalScales() =
        Assertions.assertEquals(
            "MA*kg/(ms)²",
            ProductOfBaseQuantities(
                NoN,
                linear(mega(ampere)),
                NoL,
                NoJ,
                linear(kilo(gramm)),
                NoTheta,
                reciproke(square(milli(second)))).toString())

    @Test fun withGlobalScale() =
        Assertions.assertEquals(
            "kA*kg/s²",
            kilo(ProductOfBaseQuantities(
                NoN,
                linear(ampere),
                NoL,
                NoJ,
                linear(kilo(gramm)),
                NoTheta,
                reciproke(square(second)))).toString())

    @Test fun globalScaleOverrulesLocalScale() =
        Assertions.assertEquals(
            "kA*kg/s²",
            kilo(ProductOfBaseQuantities(
                NoN,
                linear(mega(ampere)),
                NoL,
                NoJ,
                linear(kilo(gramm)),
                NoTheta,
                reciproke(square(milli(second))))).toString())
}