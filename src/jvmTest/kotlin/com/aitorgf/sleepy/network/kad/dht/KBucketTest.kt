package com.aitorgf.sleepy.network.kad.dht

import com.aitorgf.sleepy.network.kad.KadPeer
import com.aitorgf.sleepy.types.UInt128
import io.ktor.util.network.*
import org.junit.Test
import kotlin.concurrent.thread
import kotlin.test.assertEquals
import kotlin.test.assertTrue

class KBucketTest {
    data class TestKadPeer(
        override val kadId: UInt128 = UInt128.random(),
        override val kadDistance: UInt128 = UInt128.random(),
        override val ipv4: NetworkAddress = NetworkAddress("", 0)
    ) : KadPeer

    @Test
    fun multiThreadUseTest() {
        val exampleBucket = KBucket()

        val tkp1 = TestKadPeer()
        val tkp2 = TestKadPeer()
        val tkp3 = TestKadPeer()
        val tkp4 = TestKadPeer()

        var t1Count: Int = 0
        val t1 = thread(true) {
            exampleBucket.add(tkp1)
            exampleBucket.add(tkp2)
            exampleBucket.add(tkp3)
            exampleBucket.add(tkp4)
            for (i in 0 until KBucket.MaxBucketSize) {
                exampleBucket.add(TestKadPeer())
            }
            t1Count = exampleBucket.count()
        }

        var t2Count: Int = 0
        var t2Contains: Boolean = false
        val t2 = thread(true) {
            exampleBucket.add(tkp1)
            exampleBucket.add(tkp4)
            t2Count = exampleBucket.count()
            t2Contains = exampleBucket.contains(tkp1)
        }

        t2.join()
        assertTrue(t2Contains)
        assertTrue(t2Count >= 2)

        t1.join()
        assertEquals(KBucket.MaxBucketSize, t1Count)
    }
}
