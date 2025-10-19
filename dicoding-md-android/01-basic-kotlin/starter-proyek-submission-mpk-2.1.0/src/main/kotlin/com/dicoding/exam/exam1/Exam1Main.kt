package com.dicoding.exam.exam1

private fun main() {
    println(
        """
            ${isEvenNumber(1)}
            ${moreThanFive(1)}
            ${result(1)}
        """.trimIndent()
    )

    println(
        """
            ${isEvenNumber(10)}
            ${moreThanFive(10)}
            ${result(10)}
        """.trimIndent()
    )
}