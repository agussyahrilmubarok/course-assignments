package com.dicoding.storyapp.features.welcome

import androidx.lifecycle.ViewModel
import androidx.lifecycle.asLiveData
import com.dicoding.storyapp.data.CommonRepository
import com.dicoding.storyapp.data.local.UserPreferences
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.runBlocking

class WelcomeViewModel(
    private val userPref: UserPreferences,
    private val common: CommonRepository
) : ViewModel() {

    fun isLogged(): Boolean {
        val user = runBlocking {
            userPref.getUser().first()
        }

        return user.isLogged
    }

    val isOnline = common.isConnected.asLiveData()

}