package com.dicoding.exam.exam5

import kotlinx.coroutines.async
import kotlinx.coroutines.runBlocking

private fun main() = runBlocking {
    println("Counting...")

    val resultSum = async { sum(10, 10) }
    val resultMultiple = async { multiple(20, 20) }

    println(
        """
            Result sum: ${resultSum.await()}
            Result multiple: ${resultMultiple.await()}
        """.trimIndent()
    )
}