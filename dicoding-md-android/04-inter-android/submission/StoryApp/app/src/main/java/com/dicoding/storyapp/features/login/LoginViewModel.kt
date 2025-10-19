package com.dicoding.storyapp.features.login

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.dicoding.storyapp.data.AuthRepository
import com.dicoding.storyapp.data.local.UserPreferences
import com.dicoding.storyapp.data.model.LoginRequest
import com.dicoding.storyapp.data.model.UserModel
import kotlinx.coroutines.launch

class LoginViewModel(
    private val autRepo: AuthRepository,
    private val userPref: UserPreferences
) : ViewModel() {

    fun login(param: LoginRequest) = autRepo.login(param)

    fun savedUser(param: UserModel) {
        viewModelScope.launch {
            userPref.saveUser(param)
        }
    }

}