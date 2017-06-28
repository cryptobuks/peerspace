package straightway.general.units

import straightway.general.numbers.times

class Length constructor(symbol: String, scale: UnitScale, baseMagnitude: Number)
    : QuantityBase("m", symbol, scale, baseMagnitude, { Length(symbol, it, baseMagnitude) })

val meter = Length("m", uni, 1)
val inch = Length("\"", uni, 0.0254)
val foot = Length("ft", uni, 12 * inch.baseMagnitude)
val yard = Length("yd", uni, 3 * foot.baseMagnitude)
val mile = Length("mile", uni, 1760 * yard.baseMagnitude)
val nauticalMile = Length("NM", uni, 1852)