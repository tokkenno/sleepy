package com.aitorgf.sleepy.types

import kotlin.test.Test
import kotlin.test.assertEquals

class UInt128Test {
    @Test
    fun constructorFromShort() {
        val testNum = 5634.toShort()

        val uint1 = UInt128(testNum)
        val uint2 = UInt128(testNum.toLong())

        assertEquals(testNum, uint1.toShort())
        assertEquals(uint1, uint2)

        val uint3 = UInt128(-testNum)
        val uint4 = UInt128(-1L, ULong.MAX_VALUE.toLong() - testNum + 1)

        assertEquals(uint3, uint4)
    }

    @Test
    fun constructorFromInt() {
        val testNum = 75634

        val uint1 = UInt128(testNum)
        val uint2 = UInt128(testNum.toLong())

        assertEquals(testNum, uint1.toInt())
        assertEquals(uint1, uint2)

        val uint3 = UInt128(-testNum)
        val uint4 = UInt128(-1L, ULong.MAX_VALUE.toLong() - testNum + 1)

        assertEquals(uint3, uint4)
    }

    @Test
    fun constructorFromLong() {
        val testNum = 756346L

        val uint1 = UInt128(testNum)

        assertEquals(testNum, uint1.toLong())

        val uint3 = UInt128(-testNum)
        val uint4 = UInt128(-1L, -testNum)

        assertEquals(uint3, uint4)
    }

    @Test
    fun toBitSet() {
        val uint1 = UInt128(3434689892479043495L, 4363476556766674563L)
        println(uint1.toBitSet())
        assertEquals(
            BitSet(
                arrayOf(
                    0, 0, 1, 0, 1, 1, 1, 1,
                    1, 0, 1, 0, 1, 0, 1, 0,
                    0, 1, 1, 1, 1, 0, 0, 0,
                    0, 1, 0, 0, 0, 0, 1, 0,
                    0, 1, 1, 0, 1, 0, 1, 0,
                    1, 1, 0, 1, 0, 1, 1, 0,
                    1, 0, 0, 0, 0, 1, 1, 1,
                    1, 0, 1, 0, 0, 1, 1, 1,
                    0, 0, 1, 1, 1, 1, 0, 0,
                    1, 0, 0, 0, 1, 1, 1, 0,
                    0, 0, 1, 0, 1, 1, 1, 0,
                    1, 1, 0, 0, 1, 1, 1, 1,
                    0, 1, 0, 0, 0, 1, 1, 1,
                    1, 1, 0, 1, 0, 0, 1, 0,
                    0, 1, 1, 1, 1, 1, 1, 0,
                    1, 0, 0, 0, 0, 0, 1, 1
                )
            ), uint1.toBitSet()
        )
    }
}
