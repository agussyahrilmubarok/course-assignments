package com.dicoding.exam.optionalexam3

private fun main() {
    println(manipulateString("Dicoding000010", 100) == "Dicoding1000")
    println(manipulateString("dicoding", 99) == "dicoding99")
    println(manipulateString("dicoding123", 10) == "dicoding1230")
}