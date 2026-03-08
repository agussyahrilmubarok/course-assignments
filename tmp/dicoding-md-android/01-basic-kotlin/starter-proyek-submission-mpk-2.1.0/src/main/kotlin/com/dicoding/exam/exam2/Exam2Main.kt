package com.dicoding.exam.exam2

private fun main() {
    val valueA = 101
    val valueB = 52
    val valueC = 99

    val resultA = calculate(valueA, valueB, valueC)
    val resultB = calculate(valueA, valueB, null)

    println(
        """
        ${result(resultA)}
        ${result(resultB)}
        """.trimIndent()
    )
}
