package com.aitorgf.sleepy.types

fun Long.toBin(): String {
    val buffer = ByteArray(8)
    buffer.setLongAt(0, this)
    return buffer.toBin()
}

fun ULong.toBin(): String {
    val buffer = ByteArray(8)
    buffer.setULongAt(0, this)
    return buffer.toBin()
}
