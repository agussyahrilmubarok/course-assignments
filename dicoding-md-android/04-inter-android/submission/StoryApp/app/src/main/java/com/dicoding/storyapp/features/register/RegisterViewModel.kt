package com.dicoding.storyapp.features.register

import androidx.lifecycle.ViewModel
import com.dicoding.storyapp.data.AuthRepository
import com.dicoding.storyapp.data.model.RegisterRequest

class RegisterViewModel(
    private val authRepo: AuthRepository
) : ViewModel() {

    fun register(param: RegisterRequest) = authRepo.register(param)

}