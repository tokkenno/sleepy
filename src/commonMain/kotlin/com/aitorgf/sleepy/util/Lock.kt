package com.aitorgf.sleepy.util

internal expect class Lock() {
    fun lock()
    fun unlock()
}

internal inline fun <R> Lock.use(block: () -> R): R {
    try {
        lock()
        return block()
    } finally {
        unlock()
    }
}
