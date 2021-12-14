package com.aitorgf.sleepy.network

import io.ktor.util.network.*

interface Peer {
    /** Public IPv4 */
    val ipv4: NetworkAddress?
}
