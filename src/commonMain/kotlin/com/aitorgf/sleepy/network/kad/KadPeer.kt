package com.aitorgf.sleepy.network.kad

import com.aitorgf.sleepy.network.Peer
import com.aitorgf.sleepy.types.UInt128
import io.ktor.util.network.*

interface KadPeer: Peer {
    val kadId: UInt128
    val kadDistance: UInt128 // kadId xor localhostId
    override val ipv4: NetworkAddress
}
