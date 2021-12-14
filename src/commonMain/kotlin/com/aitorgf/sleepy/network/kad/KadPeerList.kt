package com.aitorgf.sleepy.network.kad

import com.aitorgf.sleepy.types.UInt128

fun List<KadPeer>.getNearest(targetID: UInt128, exceptList: List<KadPeer>): KadPeer? {
    var result = -1
    var biggerRang = -1

    for (i in this.indices) {
        val contact: KadPeer = this[i]
        if (exceptList.contains(contact)) continue
        val distance = targetID xor contact.kadDistance
        if (result == -1) {
            result = i
            biggerRang = distance.toBitSet().firstOneIndex
            continue
        }
        val rang: Int = distance.toBitSet().firstOneIndex
        if (rang < biggerRang) {
            result = i
            biggerRang = rang
        }
    }
    return if (result == -1) null else this[result]
}
