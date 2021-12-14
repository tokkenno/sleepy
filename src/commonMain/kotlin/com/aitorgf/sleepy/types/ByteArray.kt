package com.aitorgf.sleepy.types

import io.ktor.utils.io.core.*
import kotlin.math.min

fun ByteArray.getLongAt(index: Int, endian: ByteOrder = ByteOrder.nativeOrder()): Long {
    return this.getULongAt(index, endian).toLong()
}

fun ByteArray.getULongAt(index: Int, endian: ByteOrder = ByteOrder.nativeOrder()): ULong {
    if (this.size < index) {
        throw IndexOutOfBoundsException("The index is over the array size")
    }

    val buffer = ByteArray(8)
    for ((bufferIndex, sourceIndex) in (index..min(index + 7, this.size)).withIndex()) {
        buffer[bufferIndex] = this[sourceIndex]
    }


    return if (endian == ByteOrder.BIG_ENDIAN) {
        (buffer[7].toULong() shl 56
                or (buffer[6].toULong() and 0xffu shl 48
                ) or (buffer[5].toULong() and 0xffu shl 40
                ) or (buffer[4].toULong() and 0xffu shl 32
                ) or (buffer[3].toULong() and 0xffu shl 24
                ) or (buffer[2].toULong() and 0xffu shl 16
                ) or (buffer[1].toULong() and 0xffu shl 8
                ) or (buffer[0].toULong() and 0xffu))
    } else {
        (buffer[0].toULong() shl 56
                or (buffer[1].toULong() and 0xffu shl 48
                ) or (buffer[2].toULong() and 0xffu shl 40
                ) or (buffer[3].toULong() and 0xffu shl 32
                ) or (buffer[4].toULong() and 0xffu shl 24
                ) or (buffer[5].toULong() and 0xffu shl 16
                ) or (buffer[6].toULong() and 0xffu shl 8
                ) or (buffer[7].toULong() and 0xffu))
    }
}

fun ByteArray.setLongAt(index: Int, value: Long, endian: ByteOrder = ByteOrder.nativeOrder()) {
    this.setULongAt(index, value.toULong(), endian)
}

fun ByteArray.setULongAt(index: Int, value: ULong, endian: ByteOrder = ByteOrder.nativeOrder()) {
    if (this.size < index) {
        throw IndexOutOfBoundsException("The index is over the array size")
    }

    val buffer = if (endian == ByteOrder.BIG_ENDIAN) {
        byteArrayOf(
            value.toByte(),
            (value shr 8).toByte(),
            (value shr 16).toByte(),
            (value shr 24).toByte(),
            (value shr 32).toByte(),
            (value shr 40).toByte(),
            (value shr 48).toByte(),
            (value shr 56).toByte()
        )
    } else {
        byteArrayOf(
            (value shr 56).toByte(),
            (value shr 48).toByte(),
            (value shr 40).toByte(),
            (value shr 32).toByte(),
            (value shr 24).toByte(),
            (value shr 16).toByte(),
            (value shr 8).toByte(),
            value.toByte(),
        )
    }

    for ((bufferIndex, sourceIndex) in (index..min(index + 7, this.size)).withIndex()) {
        this[sourceIndex] = buffer[bufferIndex]
    }
}

fun ByteArray.toHex(): String {
    val hexChars = "0123456789abcdef".toCharArray()
    var result = ""

    forEach {
        val octet = it.toInt()
        val firstIndex = (octet and 0xF0).ushr(4)
        val secondIndex = octet and 0x0F
        result += hexChars[firstIndex].toString() + hexChars[secondIndex].toString()
    }

    return result
}

fun ByteArray.toBin(): String {
    return this.joinToString(" ") { it.toUByte().toString(2).padStart(8, '0') }
}
