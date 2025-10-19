package com.dicoding.exam.exam3

private fun main() {
    println(
        """
        'Dicoding Indonesia' is String? ${checkType("Dicoding Indonesia")}
        'True' is Boolean? ${checkType(true)}
        '1' is Integer? ${checkType(1)}
        '10.01' is Double? ${checkType(10.01)}
        '[10, 9, 8 , 6]' is List? ${checkType(listOf(10, 9, 8, 6))}
        '{score=10}' is Map? ${checkType(mapOf("score" to 10))}
        """.trimIndent()
    )
}
