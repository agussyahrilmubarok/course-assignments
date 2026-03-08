package com.dicoding.storyapp.features.home

import androidx.lifecycle.ViewModel
import com.dicoding.storyapp.data.StoryRepository
import com.dicoding.storyapp.data.local.UserPreferences
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.runBlocking

class HomeViewModel(
    private val storyRepo: StoryRepository,
    private val userPref: UserPreferences
) : ViewModel() {

    fun getAllStories() = storyRepo.getAllStories(getToken())

    fun logout(): Boolean {
        runBlocking {
            userPref.setLogout()
        }
        return true
    }

    private fun getToken(): String {
        val token = runBlocking {
            userPref.getToken().first()
        }
        return "Bearer $token"
    }
}