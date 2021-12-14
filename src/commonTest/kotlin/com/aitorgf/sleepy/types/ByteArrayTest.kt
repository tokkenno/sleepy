package com.aitorgf.sleepy.types

import io.ktor.utils.io.core.*
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertTrue

class ByteArrayTest {
    private val littleEndianULong =
        byteArrayOf(0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 10.toByte(), 245.toByte(), 72.toByte(), 210.toByte())
    private val bigEndianULong =
        byteArrayOf(1, 0, 0, 0, 0, 0, 0, 0, 210.toByte(), 72.toByte(), 245.toByte(), 10.toByte(), 0, 0, 0, 0)

    private val littleEndianLong = byteArrayOf(
        0,
        0,
        0,
        0,
        10.toByte(),
        245.toByte(),
        72.toByte(),
        210.toByte(),
        255.toByte(),
        255.toByte(),
        255.toByte(),
        255.toByte(),
        245.toByte(),
        10.toByte(),
        183.toByte(),
        46.toByte()
    )
    private val bigEndianLong = byteArrayOf(
        210.toByte(),
        72.toByte(),
        245.toByte(),
        10.toByte(),
        0,
        0,
        0,
        0,
        46.toByte(),
        183.toByte(),
        10.toByte(),
        245.toByte(),
        255.toByte(),
        255.toByte(),
        255.toByte(),
        255.toByte()
    )

    @Test
    fun getULongAt() {
        assertEquals(1uL, this.littleEndianULong.getULongAt(0, ByteOrder.LITTLE_ENDIAN))
        assertEquals(183847122uL, this.littleEndianULong.getULongAt(8, ByteOrder.LITTLE_ENDIAN))

        assertEquals(1uL, this.bigEndianULong.getULongAt(0, ByteOrder.BIG_ENDIAN))
        assertEquals(183847122uL, this.bigEndianULong.getULongAt(8, ByteOrder.BIG_ENDIAN))
    }

    @Test
    fun setULongAt() {
        val buffer = ByteArray(16)

        buffer.setULongAt(0, 1uL, ByteOrder.LITTLE_ENDIAN)
        buffer.setULongAt(8, 183847122uL, ByteOrder.LITTLE_ENDIAN)
        assertTrue(
            this.littleEndianULong.contentEquals(buffer),
            "Content mismatch at little endian. Expected ${littleEndianULong.toHex()}. Actual ${buffer.toHex()}"
        )

        buffer.setULongAt(0, 1uL, ByteOrder.BIG_ENDIAN)
        buffer.setULongAt(8, 183847122uL, ByteOrder.BIG_ENDIAN)
        assertTrue(
            this.bigEndianULong.contentEquals(buffer),
            "Content mismatch at big endian. Expected ${bigEndianULong.toHex()}. Actual ${buffer.toHex()}"
        )
    }

    @Test
    fun getAndSetULongAt() {
        val cpy = this.littleEndianLong.copyOf()
        val value = 123456uL

        cpy.setULongAt(3, value)

        assertEquals(value, cpy.getULongAt(3))
    }

    @Test
    fun setLongAt() {
        val buffer = ByteArray(16)

        buffer.setLongAt(0, 183847122L, ByteOrder.LITTLE_ENDIAN)
        buffer.setLongAt(8, -183847122L, ByteOrder.LITTLE_ENDIAN)
        assertTrue(
            this.littleEndianLong.contentEquals(buffer),
            "Content mismatch at little endian. Expected ${littleEndianLong.toHex()}. Actual ${buffer.toHex()}"
        )

        buffer.setLongAt(0, 183847122L, ByteOrder.BIG_ENDIAN)
        buffer.setLongAt(8, -183847122L, ByteOrder.BIG_ENDIAN)
        assertTrue(
            this.bigEndianLong.contentEquals(buffer),
            "Content mismatch at big endian. Expected ${bigEndianLong.toHex()}. Actual ${buffer.toHex()}"
        )
    }

    @Test
    fun getLongAt() {
        assertEquals(183847122L, this.littleEndianLong.getLongAt(0, ByteOrder.LITTLE_ENDIAN))
        assertEquals(-183847122L, this.littleEndianLong.getLongAt(8, ByteOrder.LITTLE_ENDIAN))

        assertEquals(183847122L, this.bigEndianLong.getLongAt(0, ByteOrder.BIG_ENDIAN))
        assertEquals(-183847122L, this.bigEndianLong.getLongAt(8, ByteOrder.BIG_ENDIAN))
    }

    @Test
    fun getAndSetLongAt() {
        val cpy = this.littleEndianLong.copyOf()
        val value = -123456L

        cpy.setLongAt(5, value)

        assertEquals(value, cpy.getLongAt(5))
    }

    @Test
    fun toHex() {
        val str = "ffab4500"
        val btr = byteArrayOf(255.toByte(), 171.toByte(), 69.toByte(), 0)

        assertEquals(str, btr.toHex())
    }

    @Test
    fun toBin() {
        val str = "11111111 10101011 01000101 00000000"
        val btr = byteArrayOf(255.toByte(), 171.toByte(), 69.toByte(), 0)

        assertEquals(str, btr.toBin())
    }
}
