package com.aitorgf.sleepy.types

import kotlin.test.Test
import kotlin.test.assertContentEquals

class StringTest {
    @Test
    fun decodeHex() {
        val str = "ffab4500"
        val btr = byteArrayOf(255.toByte(), 171.toByte(), 69.toByte(), 0)

        assertContentEquals(btr, str.decodeHex())
    }
}
