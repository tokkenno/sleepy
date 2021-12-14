package com.aitorgf.sleepy.types

class BitSet(val size: Int, init: (Int) -> Boolean = { false }) {
    private val bits: Array<Boolean> = Array(size, init)

    constructor(data: ByteArray) : this(data.size * 8, {
        val dataIndex = it / 8
        data[dataIndex].toInt() shr (7 - (it % 8)) and 0x1 == 0x1
    })

    constructor(data: Array<Boolean>) : this(data.size, {
        data[it]
    })

    constructor(data: Array<Int>) : this(data.size, {
        data[it] != 0
    })

    fun get(index: Int): Boolean {
        return this.bits[index]
    }

    fun set(index: Int, value: Boolean) {
        this.bits[index] = value
    }

    /**
     * Get the position of the first one on this bit set, starting from least significant bit
     */
    val firstOneIndex: Int
        get() {
            for (i in this.size - 1 downTo 0) {
                if (this.bits[i]) return i
            }
            return 0
        }

    override fun toString(): String {
        return this.bits.mapIndexed { index, value ->
            val prefix = if (index > 0 && index % 8 == 0) {
                " "
            } else ""
            if (value) {
                prefix + "1"
            } else prefix + "0"
        }.joinToString("")
    }

    override fun equals(other: Any?): Boolean {
        return if (other == null) false
        else if (this === other) true
        else if (other is BitSet) {
            this.bits.contentEquals(other.bits)
        } else false
    }

    override fun hashCode(): Int {
        return this.bits.contentHashCode()
    }
}
