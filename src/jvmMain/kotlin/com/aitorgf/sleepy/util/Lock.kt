package com.aitorgf.sleepy.util

import java.util.concurrent.locks.ReentrantLock

internal actual class Lock actual constructor() {
    private val mutex = ReentrantLock()

    actual fun lock() {
        mutex.lock()
    }
    actual fun unlock() {
        mutex.unlock()
    }
}
