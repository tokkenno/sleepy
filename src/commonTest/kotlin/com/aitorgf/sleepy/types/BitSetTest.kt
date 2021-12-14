package com.aitorgf.sleepy.types

import kotlin.test.Test
import kotlin.test.assertEquals

class BitSetTest {
    @Test
    fun creation() {
        val bSet1 = BitSet(byteArrayOf(0xA, 0x2E))
        val bSet2 = BitSet(arrayOf(0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 1, 1, 0))
        assertEquals(bSet1, bSet2)
    }
}
