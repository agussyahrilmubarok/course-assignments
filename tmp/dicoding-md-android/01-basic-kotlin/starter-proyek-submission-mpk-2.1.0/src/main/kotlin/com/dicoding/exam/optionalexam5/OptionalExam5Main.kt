package com.dicoding.exam.optionalexam5

private fun main() {
    println(concatString("Hello", "Dicoding") == "HelloDicoding")
    println(concatString("", "Dicoding") == "Dicoding")
    println(concatString("Dicoding ", "Indonesia") == "Dicoding Indonesia")
}