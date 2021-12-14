package com.aitorgf.sleepy.types

import io.ktor.utils.io.core.*
import kotlin.experimental.xor
import kotlin.math.min
import kotlin.random.Random

class UInt128 : Number, Comparable<UInt128> {
    private val bytes = ByteArray(SIZE_BYTES)

    constructor() : this(0L, 0L)

    constructor(value: Byte) : this(
        if (value >= 0) {
            0L
        } else -1L,
        if (value >= 0) {
            value.toLong() and 0x00000000000000FFL
        } else value.toLong() and 0x00000000000000FFL or 0xFFFFFFFFFFFFFF00uL.toLong()
    )

    constructor(value: Short) : this(
        if (value >= 0) {
            0L
        } else -1L,
        if (value >= 0) {
            value.toLong() and 0x000000000000FFFFL
        } else value.toLong() and 0x000000000000FFFFL or 0xFFFFFFFFFFFF0000uL.toLong()
    )

    constructor(value: Int) : this(
        if (value >= 0) {
            0L
        } else -1L,
        if (value >= 0) {
            value.toLong() and 0x00000000FFFFFFFFL
        } else value.toLong() and 0x00000000FFFFFFFFL or 0xFFFFFFFF00000000uL.toLong()
    )

    constructor(value: Long) : this(
        if (value >= 0) {
            0L
        } else -1L,
        value
    )

    constructor(upValue: Long, downValue: Long) {
        this.bytes.setLongAt(0, upValue)
        this.bytes.setLongAt(8, downValue)
    }

    constructor(value: ByteArray) {
        for (i in 0 until min(SIZE_BYTES, value.size)) {
            this.bytes[i] = value[i]
        }
    }

    fun toByteArray(): ByteArray {
        return this.bytes.copyOf()
    }

    override fun toString(): String {
        return this.bytes.toHex()
    }

    fun toBin(): String {
        return this.bytes.toBin()
    }

    /**
     * {@inheritDoc}
     */
    override operator fun compareTo(other: UInt128): Int {
        val lower63BitMask = 0x7FFFFFFFFFFFFFFFL
        val up = this.bytes.getLongAt(0)
        val low = this.bytes.getLongAt(8)
        val otherUp = other.bytes.getLongAt(0)
        val otherLow = other.bytes.getLongAt(8)

        if (up != otherUp) {
            if (up < 0 && otherUp > 0) {
                return 1
            }
            if (up > 0 && otherUp < 0) {
                return -1
            }
            return if (up and lower63BitMask
                > otherUp and lower63BitMask
            ) {
                1
            } else -1
        } else if (low != otherLow) {
            if (low < 0 && otherLow > 0) {
                return 1
            }
            if (low > 0 && otherLow < 0) {
                return -1
            }
            return if (low and lower63BitMask
                > otherLow and lower63BitMask
            ) {
                1
            } else -1
        }
        return 0
    }

    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (other == null) return false
        if (other::class != UInt128::class) return false
        other as UInt128
        return this.bytes.contentEquals(other.bytes)
    }

    override fun hashCode(): Int {
        return this.bytes.contentHashCode()
    }

    override fun toByte(): Byte {
        return if (ByteOrder.nativeOrder() == ByteOrder.BIG_ENDIAN) {
            this.bytes[SIZE_BYTES - 1]
        } else {
            this.bytes[0]
        }
    }

    override fun toChar(): Char {
        return this.toInt().toChar()
    }

    override fun toDouble(): Double {
        return this.toLong().toDouble()
    }

    override fun toFloat(): Float {
        return this.toInt().toFloat()
    }

    override fun toInt(): Int {
        var result = 0
        for ((ind, byteInd) in (SIZE_BYTES - 1 downTo SIZE_BYTES - 4).withIndex()) {
            result = result or (this.bytes[byteInd].toInt() and 0xFF shl (ind * 8))
        }
        return result
    }

    override fun toLong(): Long {
        var result = 0L
        for ((ind, byteInd) in (SIZE_BYTES - 1 downTo SIZE_BYTES - 8).withIndex()) {
            result = result or (this.bytes[byteInd].toLong() and 0xFF shl (ind * 8))
        }
        return result
    }

    override fun toShort(): Short {
        var result = 0
        for ((ind, byteInd) in (SIZE_BYTES - 1 downTo SIZE_BYTES - 2).withIndex()) {
            result = result or (this.bytes[byteInd].toInt() and 0xFF shl (ind * 8))
        }
        return result.toShort()
    }

    fun toBitSet(): BitSet {
        return BitSet(this.bytes)
    }

    infix fun xor(other: UInt128): UInt128 {
        return UInt128(this.bytes.mapIndexed { index, byte -> byte xor other.bytes[index] }.toByteArray())
    }

    companion object {
        private const val SIZE_BITS: Int = 128
        const val SIZE_BYTES: Int = SIZE_BITS / Byte.SIZE_BITS

        fun random(): UInt128 {
            return UInt128(Random.Default.nextBytes(SIZE_BYTES))
        }
    }
}
