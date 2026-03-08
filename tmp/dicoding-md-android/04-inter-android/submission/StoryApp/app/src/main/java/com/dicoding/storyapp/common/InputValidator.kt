package com.dicoding.storyapp.common

object InputValidator {
    fun isValidName(name: String?): Boolean {
        return name != null && name.trim().isNotEmpty()
    }

    fun isValidEmail(email: String?): Boolean {
        return email != null && android.util.Patterns.EMAIL_ADDRESS.matcher(email).matches()
    }
}