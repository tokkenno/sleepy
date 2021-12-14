package com.aitorgf.sleepy.network.kad.dht

import com.aitorgf.sleepy.network.kad.KadPeer
import com.aitorgf.sleepy.network.kad.getNearest
import com.aitorgf.sleepy.types.UInt128
import com.aitorgf.sleepy.util.Lock
import com.aitorgf.sleepy.util.use

class KBucket {
    companion object {
        const val MaxBucketSize: Int = 16
    }

    private val peers: ArrayDeque<KadPeer> = ArrayDeque()
    private val lock = Lock()

    /**
     * TRUE if the bucket is full
     */
    val full: Boolean
        get() = this.left() == 0

    /**
     * Get the number of peers in this bucket
     */
    fun count(): Int {
        return this.peers.count()
    }

    /**
     * Get the number of free spaces in this bucket
     */
    fun left(): Int {
        return MaxBucketSize - this.count()
    }

    /**
     * Add a peer to this bucket
     */
    fun add(peer: KadPeer) {
        this.lock.use {
            if (this.peers.contains(peer)) {
                this.peers.remove(peer)
            }
            if (this.left() > 0) {
                this.peers.add(peer)
            } else throw Full
        }
    }

    /**
     * Remove a peer from this bucket
     */
    fun remove(peer: KadPeer) {
        this.lock.use {
            if (this.peers.contains(peer)) {
                this.peers.remove(peer)
            }
        }
    }

    /**
     * Remove all peers from this bucket
     */
    fun clear() {
        this.peers.clear()
    }

    /**
     * Check if this bucket contains a concrete peer
     */
    fun contains(peer: KadPeer): Boolean {
        return this.peers.contains(peer)
    }

    /**
     * Check if this bucket contains a concrete peer by his ID
     */
    fun contains(peerId: UInt128): Boolean {
        return this.get(peerId) != null
    }

    /**
     * Get a peer by his ID
     */
    fun get(peerId: UInt128): KadPeer? {
        this.lock.use {
            return this.peers.firstOrNull { it.kadId == peerId }
        }
    }

    /**
     * Get a list of all peers on this bucket
     */
    fun getAll(): Set<KadPeer> {
        this.lock.use {
            return this.peers.toSet()
        }
    }

    /**
     * Get a list of random peers from this bucket
     * @param count The number of peers to retrieve. If bucket contains less, returns the full list.
     */
    fun getRandom(count: Int): Set<KadPeer> {
        this.lock.use {
            return if (this.peers.size <= count) {
                this.peers.toSet()
            } else {
                this.peers.toList()
                    .shuffled()
                    .take(count)
                    .toSet()
            }
        }
    }

    /**
     * Get a list of the nearest peers to a peerId from this bucket
     * @param peerId The peer ID to calculate the distance.
     * @param count The number of peers to retrieve. If bucket contains less, returns the full list.
     */
    fun getNearest(peerId: UInt128, count: Int): Set<KadPeer> {
        val originalList = this.lock.use {
            this.peers.toList()
        }

        return if (originalList.size <= count) {
            originalList.toSet()
        } else {
            val list = mutableListOf<KadPeer>()
            do {
                val contact = originalList.getNearest(peerId, list)
                if (contact != null) {
                    list.add(contact)
                }
            } while (list.size < count && contact != null)
            return list.toSet()
        }
    }

    override fun toString(): String {
        return "KBucket (${this.count() / MaxBucketSize}) {\n${this.getAll().map { "\t${it.kadId}\n" }}}"
    }
}
